package report

import (
	"testing"
)

func TestCollect_NotEmpty(t *testing.T) {
	rep := Collect()

	if rep.Host.Hostname == "" {
		t.Error("hostname should not be empty")
	}

	if rep.Memory.Total == 0 {
		t.Error("memory total should not be zero")
	}

	if len(rep.CPU.Model) == 0 {
		t.Error("cpu model should not be empty")
	}
}
