package config

import (
	"testing"
	"time"
)

func TestParseRPSSchedule(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []RPSStep
		wantErr  bool
	}{
		{
			name:     "empty string disables the schedule",
			input:    "",
			expected: nil,
		},
		{
			name:     "whitespace only disables the schedule",
			input:    "   ",
			expected: nil,
		},
		{
			name:  "multi step with explicit hold",
			input: "1000:60,5000:120,10000:0",
			expected: []RPSStep{
				{RPS: 1000, Duration: 60 * time.Second},
				{RPS: 5000, Duration: 120 * time.Second},
				{RPS: 10000, Duration: 0},
			},
		},
		{
			name:  "last step may omit duration",
			input: "1000:60,5000",
			expected: []RPSStep{
				{RPS: 1000, Duration: 60 * time.Second},
				{RPS: 5000, Duration: 0},
			},
		},
		{
			name:  "single step with duration",
			input: "1250:300",
			expected: []RPSStep{
				{RPS: 1250, Duration: 300 * time.Second},
			},
		},
		{
			name:  "spaces around items are tolerated",
			input: " 100:10 , 200:20 ",
			expected: []RPSStep{
				{RPS: 100, Duration: 10 * time.Second},
				{RPS: 200, Duration: 20 * time.Second},
			},
		},
		{name: "zero rps", input: "0:60", wantErr: true},
		{name: "negative rps", input: "-10:60", wantErr: true},
		{name: "non-numeric rps", input: "abc:60", wantErr: true},
		{name: "negative duration", input: "100:-1", wantErr: true},
		{name: "non-numeric duration", input: "100:xyz", wantErr: true},
		{name: "zero duration on non-last step", input: "100:0,200:60", wantErr: true},
		{name: "omitted duration on non-last step", input: "100,200:60", wantErr: true},
		{name: "too many fields", input: "100:60:30", wantErr: true},
		{name: "empty item", input: "100:60,,200:60", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			steps, err := ParseRPSSchedule(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseRPSSchedule(%q) expected error, got %v", tt.input, steps)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseRPSSchedule(%q) unexpected error: %v", tt.input, err)
			}
			if len(steps) != len(tt.expected) {
				t.Fatalf("ParseRPSSchedule(%q) expected %d steps, got %d", tt.input, len(tt.expected), len(steps))
			}
			for i := range steps {
				if steps[i] != tt.expected[i] {
					t.Errorf("ParseRPSSchedule(%q) step %d: expected %+v, got %+v", tt.input, i, tt.expected[i], steps[i])
				}
			}
		})
	}
}
