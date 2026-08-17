package compat

import (
	"fmt"
	"testing"
)

func TestResolveLowVersion(t *testing.T) {
	// 直接测试（不注入 SystemSetting，会走默认值）
	tests := []struct {
		name           string
		apiVersion     string
		appJsonVersion string
		expected       string
	}{
		// apiVersion 范围查询
		{"apiVersion 精确命中 1.0.18", "1.0.18", "1.0", "2.0.1"},
		{"apiVersion 1.0.20 >= 1.0.18 向下匹配", "1.0.20", "1.0", "2.0.1"},
		{"apiVersion 1.0.19 >= 1.0.18 向下匹配", "1.0.19", "1.0", "2.0.1"},
		{"apiVersion 1.0.16 >= 1.0.15 命中", "1.0.16", "1.0", "1.8.0"},
		{"apiVersion 精确命中 1.0.15", "1.0.15", "1.0", "1.8.0"},
		{"apiVersion 1.0.14 无匹配但 appJson 命中", "1.0.14", "1.0", "1.8.0"},

		// appJsonVersion 范围查询
		{"appJsonVersion 1.1 精确命中", "1.0.14", "1.1", "2.0.0"},
		{"appJsonVersion 精确命中 1.0", "1.0.14", "1.0", "1.8.0"},
		{"appJsonVersion 1.0.5 < 1.1 但 >= 1.0", "1.0.14", "1.0.5", "1.8.0"},
		{"appJsonVersion 1.2 > 1.1 向下匹配", "1.0.14", "1.2", "2.0.0"},
		{"appJsonVersion 0.9 无匹配", "1.0.14", "0.9", ""},

		// 两个版本取 max
		{"两个都匹配取 max: api→2.0.1, appJson→1.8.0", "1.0.18", "1.0", "2.0.1"},
		{"两个都匹配取 max: api→1.8.0, appJson→2.0.0", "1.0.15", "1.1", "2.0.0"},

		// 都查不到
		{"都查不到", "0.0.1", "0.0", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolveLowVersion(tt.apiVersion, tt.appJsonVersion)
			if result != tt.expected {
				t.Errorf("ResolveLowVersion(%q, %q) = %q, want %q",
					tt.apiVersion, tt.appJsonVersion, result, tt.expected)
			}
			fmt.Printf("✓ %s: api=%s, appJson=%s → lowVersion=%s\n",
				tt.name, tt.apiVersion, tt.appJsonVersion, result)
		})
	}
}

func TestLookupRange(t *testing.T) {
	entries := []VersionEntry{
		{Version: "1.0.18", LowVersion: "2.0.1"},
		{Version: "1.0.15", LowVersion: "1.8.0"},
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"1.0.20", "2.0.1"},
		{"1.0.18", "2.0.1"},
		{"1.0.17", "1.8.0"},
		{"1.0.16", "1.8.0"},
		{"1.0.15", "1.8.0"},
		{"1.0.14", ""},
		{"1.0.0", ""},
		{"0.9.0", ""},
	}

	for _, tt := range tests {
		t.Run("input_"+tt.input, func(t *testing.T) {
			result := lookupRange(entries, tt.input)
			if result != tt.expected {
				t.Errorf("lookupRange(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"2.0.0", "1.9.9", 1},
		{"1.0.10", "1.0.9", 1},
		{"1.0.15", "1.0.18", -1},
		{"1.0.18", "1.0.15", 1},
		{"1.0", "1.0.0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("compareVersions(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestIsVersionCompatible(t *testing.T) {
	tests := []struct {
		current  string
		minimum  string
		expected bool
	}{
		{"2.0.0", "1.8.0", true},
		{"1.8.0", "1.8.0", true},
		{"1.7.9", "1.8.0", false},
		{"2.0.1", "2.0.1", true},
		{"1.0.15", "1.0.18", false},
		{"1.0.20", "1.0.18", true},
	}

	for _, tt := range tests {
		t.Run(tt.current+"_gte_"+tt.minimum, func(t *testing.T) {
			result := IsVersionCompatible(tt.current, tt.minimum)
			if result != tt.expected {
				t.Errorf("IsVersionCompatible(%q, %q) = %v, want %v", tt.current, tt.minimum, result, tt.expected)
			}
		})
	}
}

func TestMarshalUnmarshalEntries(t *testing.T) {
	entries := []VersionEntry{
		{Version: "1.0.18", LowVersion: "2.0.1"},
		{Version: "1.0.15", LowVersion: "1.8.0"},
	}

	jsonStr, err := MarshalEntriesToJSON(entries)
	if err != nil {
		t.Fatalf("MarshalEntriesToJSON error: %v", err)
	}

	decoded, err := UnmarshalEntriesFromJSON(jsonStr)
	if err != nil {
		t.Fatalf("UnmarshalEntriesFromJSON error: %v", err)
	}

	if len(decoded) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(decoded))
	}
	if decoded[0].Version != "1.0.18" || decoded[0].LowVersion != "2.0.1" {
		t.Errorf("entry 0 mismatch: got %+v", decoded[0])
	}
	if decoded[1].Version != "1.0.15" || decoded[1].LowVersion != "1.8.0" {
		t.Errorf("entry 1 mismatch: got %+v", decoded[1])
	}

	fmt.Printf("✓ JSON round-trip: %s → %d entries\n", jsonStr, len(decoded))
}
