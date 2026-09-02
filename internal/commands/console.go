package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/basaltic-sh/cli/internal/auth"
	"github.com/basaltic-sh/cli/internal/cli"
	"github.com/basaltic-sh/sdk-go/compute"
)

func init() { cli.RegisterAt([]string{"compute", "instance"}, newConsoleCommand) }

// escapeKey ends the session: Ctrl-] , 0x1d.
//
// Raw mode is the point of a serial console — Ctrl-C has to reach the guest so
// that a command running inside the instance can be interrupted — which means
// Ctrl-C cannot also be how you leave. Ctrl-] is what telnet and virsh console
// use, so the muscle memory transfers.
const escapeKey = 0x1d

// newConsoleCommand is written by hand rather than generated.
//
// The endpoint is a WebSocket upgrade carrying a raw tty: bidirectional binary
// frames, a terminal in raw mode, and an escape key. A generated
// request/response command would compile and then not work.
func newConsoleCommand(state *cli.State) *cobra.Command {
	var backlogBytes int

	cmd := &cobra.Command{
		Use:   "serial-console <instance-id>",
		Short: "Open an interactive serial console on an instance",
		Long: "Open an interactive serial console on an instance.\n\n" +
			"The terminal is put in raw mode, so every key reaches the guest — including\n" +
			"Ctrl-C, which is usually what you want it for. Press Ctrl-] to disconnect.\n\n" +
			"This is a live tty, not a log. Use `compute instance console-output` to read\n" +
			"what the instance has already printed.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSerialConsole(cmd.Context(), state, args[0], backlogBytes)
		},
	}
	cmd.Flags().IntVar(&backlogBytes, "backlog-bytes", 0,
		"Replay this many bytes of recent output before the live stream begins")

	return cmd
}

func runSerialConsole(ctx context.Context, state *cli.State, instanceID string, backlogBytes int) error {
	if strings.TrimSpace(instanceID) == "" {
		return errors.New("instance_id must not be empty")
	}
	cfg, err := state.SDK()
	if err != nil {
		return err
	}

	endpoint, err := cfg.EndpointResolver.ResolveEndpoint(compute.ServiceID, cfg.Region)
	if err != nil {
		return err
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("the compute endpoint is not a URL: %w", err)
	}
	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	case "http":
		u.Scheme = "ws"
	}
	u.Path = "/v1/instances/" + url.PathEscape(instanceID) + "/console/serial"
	if backlogBytes > 0 {
		u.RawQuery = "backlog_bytes=" + strconv.Itoa(backlogBytes)
	}

	// Headers, not a ticket. A ticket exists because a browser's WebSocket
	// constructor takes a URL and nothing else, so it cannot send an
	// Authorization header; it costs a round trip and puts a credential in the
	// edge's access log. A Go client can send headers, so it should.
	token, err := cfg.TokenSource.Token(ctx)
	if err != nil {
		return err
	}
	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)
	if cfg.AccountID != "" {
		// The query parameter exists only for the browser; the header wins.
		header.Set("X-Account-Id", cfg.AccountID)
	}
	header.Set("User-Agent", "basaltic-cli/"+auth.Version)

	dialer := &websocket.Dialer{
		HandshakeTimeout: 30 * time.Second,
		Proxy:            http.ProxyFromEnvironment,
	}
	if tr, ok := cfg.HTTPClient.Transport.(*http.Transport); ok && tr.TLSClientConfig != nil {
		// Mirror the profile's TLS decision rather than inventing another.
		dialer.TLSClientConfig = tr.TLSClientConfig
	}

	conn, resp, err := dialer.DialContext(ctx, u.String(), header)
	if err != nil {
		return dialError(err, resp)
	}
	defer conn.Close()

	return pump(ctx, conn, state)
}

// pump copies the terminal to the connection and back until the escape key,
// an error, or the far end closing.
func pump(ctx context.Context, conn *websocket.Conn, state *cli.State) error {
	stdinFD := int(os.Stdin.Fd())
	interactive := term.IsTerminal(stdinFD)

	if interactive {
		restore, err := term.MakeRaw(stdinFD)
		if err != nil {
			return fmt.Errorf("putting the terminal in raw mode: %w", err)
		}
		// Restoring matters more than anything else here: leaving a terminal
		// raw makes the shell unusable afterwards.
		defer term.Restore(stdinFD, restore)
		fmt.Fprint(state.Printer().Err, "Connected. Press Ctrl-] to disconnect.\r\n")
	}

	errc := make(chan error, 2)

	// Guest to terminal.
	go func() {
		for {
			msgType, data, err := conn.ReadMessage()
			if err != nil {
				errc <- readError(err)
				return
			}
			if msgType != websocket.BinaryMessage && msgType != websocket.TextMessage {
				continue
			}
			if _, err := os.Stdout.Write(data); err != nil {
				errc <- err
				return
			}
		}
	}()

	// Terminal to guest.
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				if i := indexByte(buf[:n], escapeKey); i >= 0 {
					// Send whatever preceded the escape, then leave.
					if i > 0 {
						_ = conn.WriteMessage(websocket.BinaryMessage, buf[:i])
					}
					errc <- nil
					return
				}
				if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					errc <- err
					return
				}
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					errc <- nil
					return
				}
				errc <- err
				return
			}
		}
	}()

	var result error
	select {
	case result = <-errc:
	case <-ctx.Done():
		result = ctx.Err()
	}

	_ = conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(time.Second))

	if interactive {
		fmt.Fprint(state.Printer().Err, "\r\nDisconnected.\r\n")
	}
	return result
}

// readError turns a closed connection into either a clean exit or a readable
// failure.
//
// The upgrade itself succeeds even when the request is refused — the platform
// authenticates the WebSocket and only then decides it cannot serve the
// console, closing with a reason. Surfacing that raw gives
// "websocket: close 1008 (policy violation): Instance not found", which buries
// the only part that matters.
func readError(err error) error {
	if websocket.IsCloseError(err,
		websocket.CloseNormalClosure,
		websocket.CloseGoingAway,
		websocket.CloseNoStatusReceived) {
		return nil
	}
	if errors.Is(err, io.EOF) {
		return nil
	}
	var closeErr *websocket.CloseError
	if errors.As(err, &closeErr) {
		reason := strings.TrimSpace(closeErr.Text)
		if reason == "" {
			return fmt.Errorf("the serial console closed unexpectedly (code %d)", closeErr.Code)
		}
		switch closeErr.Code {
		case websocket.ClosePolicyViolation:
			// The platform's own words: "Instance not found", "instance is
			// not running", and so on.
			return errors.New(reason)
		case websocket.CloseInternalServerErr, websocket.CloseServiceRestart, websocket.CloseTryAgainLater:
			return fmt.Errorf("%s (the platform closed the console; this is worth retrying)", reason)
		}
		return errors.New(reason)
	}
	return err
}

// dialError explains an upgrade that was refused. The handshake response
// carries the reason, and without reading it every failure reads as "bad
// handshake".
func dialError(err error, resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("connecting to the serial console: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	detail := strings.TrimSpace(string(body))

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("the serial console refused the credential (http 401). Check `basaltic auth status`")
	case http.StatusForbidden:
		return fmt.Errorf("not permitted to open a serial console on this instance (http 403)")
	case http.StatusNotFound:
		return fmt.Errorf("no such instance, or it is not visible to this account (http 404)")
	case http.StatusConflict:
		return fmt.Errorf("the instance is not in a state that can serve a console (http 409): %s", detail)
	}
	if detail != "" {
		return fmt.Errorf("connecting to the serial console: http %d: %s", resp.StatusCode, detail)
	}
	return fmt.Errorf("connecting to the serial console: http %d: %w", resp.StatusCode, err)
}

func indexByte(b []byte, c byte) int {
	for i, x := range b {
		if x == c {
			return i
		}
	}
	return -1
}
