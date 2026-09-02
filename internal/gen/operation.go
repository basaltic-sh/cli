package main

import (
	"fmt"
	"strconv"
	"strings"
)

// emitOperation writes one leaf command.
func emitOperation(b *strings.Builder, svc service, r resourceGroup, op operation) error {
	fn := fmt.Sprintf("new%s%s%sCommand", exported(svc.Name), exported(r.Name), exported(op.Verb))

	use := op.Verb
	for _, p := range op.PathParams {
		use += " <" + flagName(p.Wire) + ">"
	}

	short := op.Summary
	if short == "" {
		short = title(op.Verb) + " " + r.Name
	}
	short = strings.TrimSuffix(short, ".")

	fmt.Fprintf(b, "// %s builds `basaltic %s %s %s`.\n", fn, svc.Name, r.Name, op.Verb)
	fmt.Fprintf(b, "func %s(state *cli.State) *cobra.Command {\n", fn)

	// Locals for everything cobra cannot bind to directly.
	decls, binds, applies, err := planFlags(svc, op)
	if err != nil {
		return err
	}
	for _, d := range decls {
		b.WriteString("\t" + d + "\n")
	}

	fmt.Fprintf(b, "\tcmd := &cobra.Command{\n")
	fmt.Fprintf(b, "\t\tUse:   %q,\n", use)
	fmt.Fprintf(b, "\t\tShort: %q,\n", short)
	fmt.Fprintf(b, "\t\tArgs:  cobra.ExactArgs(%d),\n", len(op.PathParams))
	if long := longHelp(svc, op); long != "" {
		fmt.Fprintf(b, "\t\tLong: %s,\n", quote(long))
	}
	b.WriteString("\t\tRunE: func(cmd *cobra.Command, args []string) error {\n")
	fmt.Fprintf(b, "\t\t\tc, err := %sClient(state)\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n", unexported(svc.Name))
	for _, a := range applies {
		b.WriteString("\t\t\t" + strings.ReplaceAll(a, "\n", "\n\t\t\t") + "\n")
	}
	emitCall(b, svc, op)
	b.WriteString("\t\t},\n\t}\n")

	b.WriteString("\tf := cmd.Flags()\n\t_ = f\n")
	for _, fl := range binds {
		b.WriteString("\t" + fl + "\n")
	}
	b.WriteString("\treturn cmd\n}\n\n")
	return nil
}

// planFlags decides how each input reaches the SDK call.
//
// Three shapes, because cobra can only bind to an addressable value of a type
// it knows:
//   - a value field on the request struct binds directly;
//   - an optional field is a pointer in the SDK, so it binds to a local and is
//     assigned only when the user actually passed the flag, which is what
//     keeps "set this to false" distinct from "leave it alone";
//   - anything structured takes JSON, because a nested object has no honest
//     flat representation.
func planFlags(svc service, op operation) (decls, binds, applies []string, err error) {
	if op.ParamsType != "" {
		decls = append(decls, fmt.Sprintf("var params %s.%s", svc.Package, op.ParamsType))
		for _, p := range op.Params {
			name, err := resolveFlagName(op, p.Wire)
			if err != nil {
				return nil, nil, nil, err
			}
			d, bind, apply := flagFor(svc.Package, "params", p, name)
			decls = append(decls, d...)
			binds = append(binds, bind...)
			applies = append(applies, apply...)
		}
		if op.Paginated {
			decls = append(decls, "var fetchAll bool")
			binds = append(binds, `f.BoolVar(&fetchAll, "all", false, "Fetch every page, not just the first.")`)
		}
	}

	switch op.BodyKind {
	case "json":
		decls = append(decls, fmt.Sprintf("var body %s.%s", svc.Package, strings.TrimPrefix(op.BodyType, "*")))
		decls = append(decls, "var bodyFile string")
		binds = append(binds, `f.StringVarP(&bodyFile, "from-file", "f", "", "Read the request body from a JSON or YAML file, or - for stdin. Flags override what it sets.")`)
		applies = append(applies, "if bodyFile != \"\" {\n\tif err := loadBody(bodyFile, &body); err != nil {\n\t\treturn err\n\t}\n}")
		for _, fd := range op.BodyFields {
			name, err := resolveFlagName(op, fd.Wire)
			if err != nil {
				return nil, nil, nil, err
			}
			d, bind, apply := flagFor(svc.Package, "body", fd, name)
			decls = append(decls, d...)
			binds = append(binds, bind...)
			applies = append(applies, apply...)
		}
	case "stream", "text":
		decls = append(decls, "var bodyFile string")
		binds = append(binds, `f.StringVarP(&bodyFile, "file", "f", "-", "File to send as the request body, or - for stdin.")`)
	}

	if op.Idempotent {
		decls = append(decls, "var idempotencyKey string")
		binds = append(binds, `f.StringVar(&idempotencyKey, "idempotency-key", "", "Makes this call replay-safe: retrying with the same key returns the original outcome instead of creating a second resource.")`)
	}
	return decls, binds, applies, nil
}

// flagFor emits the declaration, binding and assignment for one input.
//
// Four shapes, because cobra binds only to an addressable value of a type it
// already knows:
//   - a plain scalar on the request struct binds directly;
//   - an optional field is a pointer in the SDK, so it binds to a local and is
//     assigned only when the flag was actually given, which keeps "set this to
//     false" distinct from "leave it alone";
//   - a named type over string, a timestamp or a byte slice binds to a local
//     and is converted;
//   - anything structured takes JSON, because a nested object has no honest
//     flat representation.
func flagFor(pkg, target string, fd field, name string) (decls, binds, applies []string) {
	dest := target + "." + fd.GoName
	usage := flagUsage(fd)
	base := strings.TrimPrefix(fd.GoType, "*")

	required := func() {
		if fd.Required {
			binds = append(binds, fmt.Sprintf("_ = cmd.MarkFlagRequired(%q)", name))
		}
	}
	// local declares a scratch variable of the given Go type.
	local := func(goType string) string {
		v := unexported(fd.GoName) + "Flag"
		decls = append(decls, "var "+v+" "+goType)
		return v
	}

	switch fd.FlagKind {
	case "json":
		v := local("string")
		binds = append(binds, fmt.Sprintf("f.StringVar(&%s, %q, \"\", %s)", v, name, quote(usage+" (JSON)")))
		applies = append(applies, fmt.Sprintf(
			"if %s != \"\" {\n\tif err := json.Unmarshal([]byte(%s), &%s); err != nil {\n\t\treturn fmt.Errorf(\"--%s: %%w\", err)\n\t}\n}", v, v, dest, name))
		required()
		return decls, binds, applies

	case "time":
		v := local("string")
		binds = append(binds, fmt.Sprintf("f.StringVar(&%s, %q, \"\", %s)", v, name, quote(usage+" (RFC 3339)")))
		assign := dest + " = parsed"
		if fd.Pointer {
			assign = dest + " = &parsed"
		}
		applies = append(applies, fmt.Sprintf(
			"if %s != \"\" {\n\tparsed, err := parseTime(%s)\n\tif err != nil {\n\t\treturn fmt.Errorf(\"--%s: %%w\", err)\n\t}\n\t%s\n}", v, v, name, assign))
		required()
		return decls, binds, applies

	case "bytes":
		v := local("string")
		binds = append(binds, fmt.Sprintf("f.StringVar(&%s, %q, \"\", %s)", v, name, quote(usage)))
		assign := dest + " = []byte(" + v + ")"
		if fd.Pointer {
			assign = "decoded := []byte(" + v + ")\n\t" + dest + " = &decoded"
		}
		applies = append(applies, fmt.Sprintf("if %s != \"\" {\n\t%s\n}", v, assign))
		required()
		return decls, binds, applies

	case "stringSlice":
		binds = append(binds, fmt.Sprintf("f.StringSliceVar(&%s, %q, nil, %s)", dest, name, quote(usage)))
		required()
		return decls, binds, applies
	}

	// A named type whose underlying type is string — an enum from the
	// specification. It binds through a conversion rather than a local, so
	// the value lands in the struct directly.
	named := fd.FlagKind == "string" && base != "string"

	if !fd.Pointer {
		switch fd.FlagKind {
		case "string":
			ref := "&" + dest
			if named {
				ref = "(*string)(&" + dest + ")"
			}
			binds = append(binds, fmt.Sprintf("f.StringVar(%s, %q, \"\", %s)", ref, name, quote(usage)))
		case "bool":
			binds = append(binds, fmt.Sprintf("f.BoolVar(&%s, %q, false, %s)", dest, name, quote(usage)))
		case "int":
			binds = append(binds, fmt.Sprintf("f.IntVar(&%s, %q, 0, %s)", dest, name, quote(usage)))
		case "int64":
			binds = append(binds, fmt.Sprintf("f.Int64Var(&%s, %q, 0, %s)", dest, name, quote(usage)))
		case "float":
			binds = append(binds, fmt.Sprintf("f.Float64Var(&%s, %q, 0, %s)", dest, name, quote(usage)))
		}
		required()
		return decls, binds, applies
	}

	// Optional: assign through a local only when the flag was given.
	goType := map[string]string{"string": "string", "bool": "bool", "int": "int", "int64": "int64", "float": "float64"}[fd.FlagKind]
	if goType == "" {
		goType = "string"
	}
	v := local(goType)
	switch fd.FlagKind {
	case "string":
		binds = append(binds, fmt.Sprintf("f.StringVar(&%s, %q, \"\", %s)", v, name, quote(usage)))
	case "bool":
		binds = append(binds, fmt.Sprintf("f.BoolVar(&%s, %q, false, %s)", v, name, quote(usage)))
	case "int":
		binds = append(binds, fmt.Sprintf("f.IntVar(&%s, %q, 0, %s)", v, name, quote(usage)))
	case "int64":
		binds = append(binds, fmt.Sprintf("f.Int64Var(&%s, %q, 0, %s)", v, name, quote(usage)))
	case "float":
		binds = append(binds, fmt.Sprintf("f.Float64Var(&%s, %q, 0, %s)", v, name, quote(usage)))
	}
	rhs := "&" + v
	if named {
		rhs = fmt.Sprintf("(*%s.%s)(&%s)", pkg, base, v)
	}
	applies = append(applies, fmt.Sprintf("if cmd.Flags().Changed(%q) {\n\t%s = %s\n}", name, dest, rhs))
	required()
	return decls, binds, applies
}

// globalFlags are the CLI's own, available on every command. A generated flag
// with one of these names would shadow it, and cobra resolves the local one —
// so `--region` on that command would set an unrelated filter and there would
// be no way left to choose the region to call.
var globalFlags = map[string]bool{
	"profile": true, "api-key": true, "region": true, "account-id": true,
	"output": true, "no-headers": true, "insecure": true, "help": true,
}

// flagRenames resolve the collisions that exist. Each is a judgment recorded
// once, keyed by "operationId.wire_name".
//
// All three are cases where the API's parameter genuinely means something
// other than the CLI-wide flag of the same name, so neither can simply be
// dropped.
var flagRenames = map[string]string{
	// "Region to filter regional quotas by", not the region to call.
	"listQuotas.region": "for-region",
	// "Filter by emitter region", not the region to call.
	"searchLogs.region": "emitter-region",
	// "The account the resulting credentials act in", not the account this
	// request acts on.
	"assumeRoleWithWebIdentity.account_id": "target-account-id",
}

// resolveFlagName returns the flag name to emit, refusing to shadow a global.
//
// A new collision fails generation rather than producing a command with an
// unreachable global flag, which is invisible until someone needs it.
func resolveFlagName(op operation, wire string) (string, error) {
	name := flagName(wire)
	if !globalFlags[name] {
		return name, nil
	}
	if renamed, ok := flagRenames[op.ID+"."+wire]; ok {
		return renamed, nil
	}
	return "", fmt.Errorf(
		"operation %s has a parameter %q whose flag --%s would shadow the CLI-wide flag of that name,\n"+
			"making the global one unreachable on this command.\n"+
			"Add an entry to flagRenames keyed %q with the name to use instead",
		op.ID, wire, name, op.ID+"."+wire)
}

func flagUsage(fd field) string {
	u := fd.Doc
	if u == "" {
		u = title(fd.Wire)
	}
	u = strings.TrimSuffix(strings.TrimSpace(u), ".")
	if len(fd.Enum) > 0 {
		u += " (one of: " + strings.Join(fd.Enum, ", ") + ")"
	}
	return u
}

// emitCall writes the SDK invocation and the rendering of its result.
func emitCall(b *strings.Builder, svc service, op operation) {
	var callArgs []string
	callArgs = append(callArgs, "cmd.Context()")
	for i := range op.PathParams {
		callArgs = append(callArgs, fmt.Sprintf("args[%d]", i))
	}
	switch {
	case op.ParamsType != "":
		callArgs = append(callArgs, "&params")
	case op.BodyKind == "json":
		if strings.HasPrefix(op.BodyType, "*") {
			callArgs = append(callArgs, "&body")
		} else {
			callArgs = append(callArgs, "body")
		}
	case op.BodyKind == "stream":
		b.WriteString("\t\t\treader, closeBody, err := openBody(bodyFile)\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\tdefer closeBody()\n")
		callArgs = append(callArgs, "reader")
	case op.BodyKind == "text":
		b.WriteString("\t\t\ttext, err := readBody(bodyFile)\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n")
		callArgs = append(callArgs, "text")
	}

	if op.Idempotent {
		b.WriteString("\t\t\tvar reqOpts []basaltic.RequestOption\n")
		b.WriteString("\t\t\tif idempotencyKey != \"\" {\n\t\t\t\treqOpts = append(reqOpts, basaltic.WithIdempotencyKey(idempotencyKey))\n\t\t\t}\n")
		callArgs = append(callArgs, "reqOpts...")
	}

	call := fmt.Sprintf("c.%s(%s)", op.GoName, strings.Join(callArgs, ", "))

	switch op.ResultKind {
	case "none":
		fmt.Fprintf(b, "\t\t\tif err := %s; err != nil {\n\t\t\t\treturn err\n\t\t\t}\n", call)
		fmt.Fprintf(b, "\t\t\tstate.Printer().Done(%q)\n\t\t\treturn nil\n", doneMessage(op))
	case "page":
		if op.Paginated {
			allArgs := append([]string{}, callArgs...)
			// The iterator takes the same arguments; only the method differs.
			fmt.Fprintf(b, "\t\t\tif fetchAll {\n\t\t\t\treturn state.Printer().Iter(c.%sAll(%s))\n\t\t\t}\n",
				op.GoName, strings.Join(allArgs, ", "))
		}
		fmt.Fprintf(b, "\t\t\tpage, err := %s\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n", call)
		b.WriteString("\t\t\treturn state.Printer().Page(page)\n")
	case "stream":
		fmt.Fprintf(b, "\t\t\tstream, err := %s\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n", call)
		b.WriteString("\t\t\treturn state.Printer().Stream(stream)\n")
	default: // value
		fmt.Fprintf(b, "\t\t\tout, err := %s\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n", call)
		b.WriteString("\t\t\treturn state.Printer().Value(out)\n")
	}
}

// doneMessage is what a command with no output says on success. Silence would
// be ambiguous — a delete that printed nothing looks the same as one that did
// not run.
func doneMessage(op operation) string {
	switch op.Verb {
	case "delete":
		return "Deleted."
	case "revoke":
		return "Revoked."
	case "detach":
		return "Detached."
	case "attach":
		return "Attached."
	}
	return title(op.Verb) + " requested."
}

func longHelp(svc service, op operation) string {
	var parts []string
	if op.Summary != "" {
		parts = append(parts, strings.TrimSuffix(op.Summary, ".")+".")
	}
	if op.Paginated {
		parts = append(parts, "Returns one page. Pass --all to walk every page.")
	}
	if op.Idempotent {
		parts = append(parts, "Pass --idempotency-key to make this call replay-safe, which also\nmakes it safe for the CLI to retry.")
	}
	if len(parts) < 2 {
		return ""
	}
	return strings.Join(parts, "\n\n")
}

func quote(s string) string { return strconv.Quote(s) }

func quoteList(items []string) string {
	parts := make([]string, len(items))
	for i, s := range items {
		parts[i] = strconv.Quote(s)
	}
	return strings.Join(parts, ", ")
}
