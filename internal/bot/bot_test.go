package bot

import (
	"testing"
)

func TestCountSelectedOptions(t *testing.T) {
	tests := []struct {
		name     string
		selected map[int][]string
		expected map[int]int
	}{
		{
			name: "Пустой словарь",
			selected: map[int][]string{
				0: {},
			},
			expected: map[int]int{},
		},
		{
			name: "Один вариант",
			selected: map[int][]string{
				0: {"1"},
			},
			expected: map[int]int{1: 1},
		},
		{
			name: "Несколько вариантов",
			selected: map[int][]string{
				0: {"1", "2"},
				1: {"3", "4"},
			},
			expected: map[int]int{1: 1, 2: 1, 3: 1, 4: 1},
		},
		{
			name: "Повторяющиеся варианты",
			selected: map[int][]string{
				0: {"1", "1"},
				1: {"2", "2"},
			},
			expected: map[int]int{1: 2, 2: 2},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := countSelectedOptions(test.selected)
			if !mapsEqual(actual, test.expected) {
				t.Errorf("Ожидаемый результат: %+v, фактический результат: %+v", test.expected, actual)
			}
		})
	}
}

func TestGetResultDescription(t *testing.T) {
	tests := []struct {
		name     string
		counts   map[int]int
		expected string
	}{
		{
			name:     "Пустой словарь",
			counts:   map[int]int{},
			expected: "Нет выбранных вариантов.",
		},
		{
			name:     "Один вариант",
			counts:   map[int]int{1: 1},
			expected: "Описание для варианта 1",
		},
		{
			name:     "Несколько вариантов",
			counts:   map[int]int{1: 1, 2: 1},
			expected: "Описание для варианта 1",
		},
		{
			name:     "Повторяющиеся варианты",
			counts:   map[int]int{1: 2, 2: 1},
			expected: "Описание для варианта 1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := getResultDescription(test.counts)
			if actual != test.expected {
				t.Errorf("Ожидаемый результат: %s, фактический результат: %s", test.expected, actual)
			}
		})
	}
}

func mapsEqual(a, b map[int]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
