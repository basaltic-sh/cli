package commands

import (
	"github.com/basaltic-sh/cli/internal/output"
	"github.com/basaltic-sh/sdk-go/compute"
)

func init() { output.Present(compute.ListInstanceVolumesAttachment{}, volumeAttachmentRow) }

// volumeAttachmentRow is the table row for `compute instance list-volumes`.
//
// What the reflection renderer would show is the attachment's scalars, which
// leaves out the one thing an operator opens this table for: the `mount`
// report, an object, is what the in-guest agent last said about the volume.
// The rules its copy has to keep (CLOUD-39):
//
//   - `unknown` means NO agent has reported — every instance launched before
//     the reporting agent shipped says this forever. It is never blank and
//     never reads as healthy; it says "not reported".
//   - `mounted` can still carry a code: `fstab_write_failed` is a volume that
//     is mounted now and will not come back after a reboot. The code is a
//     column of its own so a mounted row cannot hide it.
//   - `message` is free text from inside the customer's guest. It is not in
//     the table (one line per row); json and yaml carry it, as text.
func volumeAttachmentRow(v any) []output.Field {
	var a *compute.ListInstanceVolumesAttachment
	switch t := v.(type) {
	case *compute.ListInstanceVolumesAttachment:
		a = t
	case compute.ListInstanceVolumesAttachment:
		a = &t
	}
	state, code := "", ""
	switch {
	case a.MountPath == "":
		// A bare block device the tenant mounts themselves: nothing
		// reconciles it, so nothing reports on it, and a state would be a
		// claim about something no one is watching.
		state = "not managed"
	case a.Mount == nil || a.Mount.State == "" || a.Mount.State == "unknown":
		state = "not reported"
	default:
		state, code = a.Mount.State, a.Mount.Code
	}
	return []output.Field{
		{Key: "volume_id", Value: a.VolumeID},
		{Key: "name", Value: a.Name},
		{Key: "status", Value: a.Status},
		{Key: "device", Value: a.Device},
		{Key: "mount_path", Value: a.MountPath},
		{Key: "mount_state", Value: state},
		{Key: "mount_code", Value: code},
	}
}
