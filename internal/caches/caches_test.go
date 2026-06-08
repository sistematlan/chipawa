package caches

import "testing"

func TestSplitJBVersion(t *testing.T) {
	cases := []struct {
		in              string
		product, suffix string
	}{
		{"PhpStorm2026.1", "PhpStorm", "2026.1"},
		{"DataGrip2025.3", "DataGrip", "2025.3"},
		{"IntelliJIdea2025.1", "IntelliJIdea", "2025.1"},
		{"Toolbox", "", ""},
		{"Daemon", "", ""},
		{"", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			gotP, gotS := splitJBVersion(tc.in)
			if gotP != tc.product || gotS != tc.suffix {
				t.Fatalf("splitJBVersion(%q) = (%q,%q); want (%q,%q)",
					tc.in, gotP, gotS, tc.product, tc.suffix)
			}
		})
	}
}

func TestParseDockerSize(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"0B", 0},
		{"512B", 512},
		{"1.5KB", 1536},                  // 1.5 * 1024
		{"2MB", 2 * 1024 * 1024},
		{"1.2GB", 1288490188},            // int64(1.2 * 1024^3) — truncated
		{"3TB", 3 * 1024 * 1024 * 1024 * 1024},
		{"  1.0GB  ", 1024 * 1024 * 1024},
		{"garbage", 0},
		{"", 0},
	}
	for _, tc := range cases {
		got := parseDockerSize(tc.in)
		if got != tc.want {
			t.Errorf("parseDockerSize(%q) = %d; want %d", tc.in, got, tc.want)
		}
	}
}
