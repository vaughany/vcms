package vcms

import (
	"runtime"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	version := Version("Test App Name")

	if !strings.Contains(version, "Test App Name") {
		t.Error("Version() was incorrect, appName not present.")
	}

	if !strings.Contains(version, AppVersion) {
		t.Error("Version() was incorrect, appVersion not present.")
	}

	if !strings.Contains(version, AppDate) {
		t.Error("Version() was incorrect, appDate not present.")
	}

	if !strings.Contains(version, runtime.Version()) {
		t.Error("Version() was incorrect, runtime.Version not present.")
	}
}

func BenchmarkVersion(b *testing.B) {
	// Run the Version function b.N times.
	for n := 0; n < b.N; n++ {
		Version("Benchmarking App Name")
	}
}
