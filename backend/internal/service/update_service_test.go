package service

import "testing"

func TestCompareVersions_IgnoresSuffixes(t *testing.T) {
	t.Parallel()

	if got := compareVersions("0.1.119-resin.1", "0.1.119"); got != 0 {
		t.Fatalf("compareVersions() = %d, want 0", got)
	}
}

func TestCompareVersions_BasicSemver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		current string
		latest  string
		want    int
	}{
		{name: "older patch", current: "0.1.118", latest: "0.1.119", want: -1},
		{name: "same version", current: "0.1.119", latest: "0.1.119", want: 0},
		{name: "newer patch", current: "0.1.120", latest: "0.1.119", want: 1},
		{name: "ignores build metadata", current: "0.1.119+local", latest: "0.1.119", want: 0},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := compareVersions(tc.current, tc.latest); got != tc.want {
				t.Fatalf("compareVersions(%q, %q) = %d, want %d", tc.current, tc.latest, got, tc.want)
			}
		})
	}
}
