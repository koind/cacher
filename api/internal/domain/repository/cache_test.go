package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCache(t *testing.T) {
	testCases := map[string]struct {
		key      string
		value    string
		expected *Cache
	}{
		"Пустые данные": {
			key:      "",
			value:    "",
			expected: &Cache{},
		},
		"Передан только ключ": {
			key:      "test-kye",
			value:    "",
			expected: &Cache{Kye: "test-kye"},
		},
		"Передано только значение": {
			key:      "",
			value:    "test-value",
			expected: &Cache{Value: "test-value"},
		},
		"Переданы все данные": {
			key:      "test-kye",
			value:    "test-value",
			expected: &Cache{Kye: "test-kye", Value: "test-value"},
		},
	}

	for title, testCase := range testCases {
		t.Run(title, func(t *testing.T) {
			assert.EqualValues(t, testCase.expected, NewCache(testCase.key, testCase.value))
		})
	}
}
