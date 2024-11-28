package types

import (
	"ledis-server/utils"
	"reflect"
	"testing"
)

func TestType(t *testing.T) {
	l := NewListType().(*listType)
	got := l.Type()
	expected := utils.ListType

	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestLLen(t *testing.T) {
	tests := []struct {
		name     string
		actions  func(l *listType)
		expected int
	}{
		{
			name: "empty list",
			actions: func(l *listType) {
				// no action
			},
			expected: 0,
		},
		{
			name: "push three elements",
			actions: func(l *listType) {
				l.RPush([]*string{ptr("a"), ptr("b"), ptr("c")}...)
			},
			expected: 3,
		},
		{
			name: "push and pop elements",
			actions: func(l *listType) {
				l.RPush([]*string{ptr("a"), ptr("b"), ptr("c")}...)
				l.LPop()
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewListType().(*listType)
			tt.actions(l)
			got := l.LLen()
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestLPop(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(l *listType)
		numToPop  int
		expected  []*string
		remaining []string
	}{
		{
			name: "Pop one value from non-empty list",
			setup: func(l *listType) {
				values := []*string{ptr("a"), ptr("b"), ptr("c")}
				l.LPush(values...)
			},
			numToPop:  1,
			expected:  []*string{ptr("c")},
			remaining: []string{"b", "a"},
		},
		{
			name: "Pop multiple values",
			setup: func(l *listType) {
				values := []*string{ptr("1"), ptr("2"), ptr("3")}
				l.LPush(values...)
			},
			numToPop:  2,
			expected:  []*string{ptr("3"), ptr("2")},
			remaining: []string{"1"},
		},
		{
			name: "Pop all values",
			setup: func(l *listType) {
				values := []*string{ptr("x"), ptr("y")}
				l.LPush(values...)
			},
			numToPop:  3,
			expected:  []*string{ptr("y"), ptr("x"), nil},
			remaining: []string{},
		},
		{
			name: "Pop from empty list",
			setup: func(l *listType) {
				// No values added
			},
			numToPop:  1,
			expected:  []*string{nil},
			remaining: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewListType().(*listType)
			tt.setup(l)

			var got []*string
			for i := 0; i < tt.numToPop; i++ {
				got = append(got, l.LPop())
			}

			if !equalSlicePtrs(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}

			remaining := l.LRange(0, l.LLen()-1)
			if !reflect.DeepEqual(remaining, tt.remaining) {
				t.Errorf("remaining elements = %v; want %v", remaining, tt.remaining)
			}
		})
	}
}

func TestRPop(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(l *listType)
		numToPop  int
		expected  []*string
		remaining []string
	}{
		{
			name: "Pop one value from non-empty list",
			setup: func(l *listType) {
				values := []*string{ptr("a"), ptr("b"), ptr("c")}
				l.RPush(values...)
			},
			numToPop:  1,
			expected:  []*string{ptr("c")},
			remaining: []string{"a", "b"},
		},
		{
			name: "Pop multiple values",
			setup: func(l *listType) {
				values := []*string{ptr("1"), ptr("2"), ptr("3")}
				l.RPush(values...)
			},
			numToPop:  2,
			expected:  []*string{ptr("3"), ptr("2")},
			remaining: []string{"1"},
		},
		{
			name: "Pop all values",
			setup: func(l *listType) {
				values := []*string{ptr("x"), ptr("y")}
				l.RPush(values...)
			},
			numToPop:  3,
			expected:  []*string{ptr("y"), ptr("x"), nil},
			remaining: []string{},
		},
		{
			name: "Pop from empty list",
			setup: func(l *listType) {
				// No values added
			},
			numToPop:  1,
			expected:  []*string{nil},
			remaining: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewListType().(*listType)
			tt.setup(l)

			var got []*string
			for i := 0; i < tt.numToPop; i++ {
				got = append(got, l.RPop())
			}

			if !equalSlicePtrs(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}

			remaining := l.LRange(0, l.LLen()-1)
			if !reflect.DeepEqual(remaining, tt.remaining) {
				t.Errorf("remaining elements = %v; want %v", remaining, tt.remaining)
			}
		})
	}
}

func equalSlicePtrs(a, b []*string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != nil && b[i] != nil && *a[i] == *b[i] {
			continue
		}
		if a[i] == b[i] {
			continue
		}
		return false
	}
	return true
}

func TestGetPositiveIndex(t *testing.T) {
	tests := []struct {
		name          string
		start         int
		end           int
		length        int
		expectedStart int
		expectedEnd   int
	}{
		{
			name:          "Negative start, positive end",
			start:         -1,
			end:           3,
			length:        5,
			expectedStart: 4,
			expectedEnd:   3,
		},
		{
			name:          "Positive start, negative end",
			start:         1,
			end:           -2,
			length:        5,
			expectedStart: 1,
			expectedEnd:   3,
		},
		{
			name:          "Negative start and end",
			start:         -3,
			end:           -1,
			length:        5,
			expectedStart: 2,
			expectedEnd:   4,
		},
		{
			name:          "Out-of-bounds start and end",
			start:         -10,
			end:           10,
			length:        5,
			expectedStart: 0,
			expectedEnd:   4,
		},
		{
			name:          "Start greater than end",
			start:         4,
			end:           2,
			length:        5,
			expectedStart: 4,
			expectedEnd:   2,
		},
		{
			name:          "Start equals end",
			start:         2,
			end:           2,
			length:        5,
			expectedStart: 2,
			expectedEnd:   2,
		},
		{
			name:          "Start and end both zero",
			start:         0,
			end:           0,
			length:        5,
			expectedStart: 0,
			expectedEnd:   0,
		},
		{
			name:          "Empty list (length zero)",
			start:         -1,
			end:           1,
			length:        0,
			expectedStart: 0,
			expectedEnd:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := utils.GetPositiveStartEndIndexes(tt.start, tt.end, tt.length)
			if start != tt.expectedStart || end != tt.expectedEnd {
				t.Errorf("GetPositiveStartEndIndexes(%d, %d, %d) = (%d, %d); want (%d, %d)",
					tt.start, tt.end, tt.length, start, end, tt.expectedStart, tt.expectedEnd)
			}
		})
	}
}

func TestLRange(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(l *listType)
		start    int
		end      int
		expected []string
	}{
		{
			name: "Normal range",
			setup: func(l *listType) {
				values := []*string{ptr("a"), ptr("b"), ptr("c"), ptr("d"), ptr("e")}
				l.RPush(values...)
			},
			start:    1,
			end:      3,
			expected: []string{"b", "c", "d"},
		},
		{
			name: "Full range",
			setup: func(l *listType) {
				values := []*string{ptr("a"), ptr("b"), ptr("c"), ptr("d"), ptr("e")}
				l.RPush(values...)
			},
			start:    0,
			end:      4,
			expected: []string{"a", "b", "c", "d", "e"},
		},
		{
			name: "Out-of-bounds range",
			setup: func(l *listType) {
				values := []*string{ptr("a"), ptr("b"), ptr("c")}
				l.RPush(values...)
			},
			start:    -10,
			end:      10,
			expected: []string{"a", "b", "c"},
		},
		{
			name: "Negative indices",
			setup: func(l *listType) {
				values := []*string{ptr("a"), ptr("b"), ptr("c"), ptr("d"), ptr("e")}
				l.RPush(values...)
			},
			start:    -3,
			end:      -1,
			expected: []string{"c", "d", "e"},
		},
		{
			name: "Empty range",
			setup: func(l *listType) {
				values := []*string{ptr("a"), ptr("b"), ptr("c")}
				l.RPush(values...)
			},
			start:    2,
			end:      1,
			expected: []string{},
		},
		{
			name: "Empty list",
			setup: func(l *listType) {
				// No values added
			},
			start:    0,
			end:      4,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewListType().(*listType)
			tt.setup(l)
			got := l.LRange(tt.start, tt.end)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("LRange(%d, %d) = %v; want %v", tt.start, tt.end, got, tt.expected)
			}
		})
	}
}

func ptr(s string) *string {
	return &s
}
