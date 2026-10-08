package commands

import (
	"testing"

	"github.com/basaltic-sh/cli/internal/output"
	"github.com/basaltic-sh/sdk-go/compute"
)

func cell(fields []output.Field, key string) string {
	for _, f := range fields {
		if f.Key == key {
			return f.Value.(string)
		}
	}
	return "<missing>"
}

func TestVolumeRowSaysWhenNothingHasReported(t *testing.T) {
	// An instance older than the reporting agent: mount_path set, no report.
	fields := volumeAttachmentRow(&compute.ListInstanceVolumesAttachment{VolumeID: "v-1", MountPath: "/data"})
	if got := cell(fields, "mount_state"); got != "not reported" {
		t.Errorf("no report must read as not reported, never blank or healthy; got %q", got)
	}
	fields = volumeAttachmentRow(&compute.ListInstanceVolumesAttachment{VolumeID: "v-1", MountPath: "/data", Mount: &compute.VolumeMount{State: "unknown"}})
	if got := cell(fields, "mount_state"); got != "not reported" {
		t.Errorf("unknown must read as not reported; got %q", got)
	}
}

func TestVolumeRowKeepsTheCodeOnAMountedVolume(t *testing.T) {
	fields := volumeAttachmentRow(compute.ListInstanceVolumesAttachment{
		VolumeID: "v-1", MountPath: "/data",
		Mount: &compute.VolumeMount{State: "mounted", Code: "fstab_write_failed", Message: "<b>x</b>"},
	})
	if cell(fields, "mount_state") != "mounted" || cell(fields, "mount_code") != "fstab_write_failed" {
		t.Errorf("a mounted volume with a code must show both; got %v", fields)
	}
	for _, f := range fields {
		if f.Key == "message" || f.Value == "<b>x</b>" {
			t.Errorf("the guest's free text does not belong in the table row: %v", f)
		}
	}
}

func TestVolumeRowDistinguishesPendingFromFailedAndBareDevices(t *testing.T) {
	pending := cell(volumeAttachmentRow(&compute.ListInstanceVolumesAttachment{MountPath: "/a", Mount: &compute.VolumeMount{State: "pending"}}), "mount_state")
	failed := cell(volumeAttachmentRow(&compute.ListInstanceVolumesAttachment{MountPath: "/a", Mount: &compute.VolumeMount{State: "failed", Code: "mkfs_failed"}}), "mount_state")
	if pending == failed {
		t.Errorf("pending and failed must not look alike: %q", pending)
	}
	if got := cell(volumeAttachmentRow(&compute.ListInstanceVolumesAttachment{VolumeID: "v-2"}), "mount_state"); got != "not managed" {
		t.Errorf("a bare block device is not managed; got %q", got)
	}
}
