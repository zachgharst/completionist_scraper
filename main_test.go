package main

import "testing"

func TestNormalizeValue(t *testing.T) {
	t.Run("normalizing percentages", func(t *testing.T) {
		val := "12.34 %"

		expected := "0.1234"
		actual, err := normalizeValue(val)
		if err != nil {
			t.Errorf("%s", err)
		}

		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("normalizing large integers", func(t *testing.T) {
		val := "12,345"

		expected := "12345"
		actual, err := normalizeValue(val)
		if err != nil {
			t.Errorf("%s", err)
		}

		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("normalizing timespans", func(t *testing.T) {
		val := "12h 34m"

		expected := "12.57"
		actual, err := normalizeValue(val)
		if err != nil {
			t.Errorf("%s", err)
		}

		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("normalizing timespans with hour only", func(t *testing.T) {
		val := "12h"

		expected := "12.00"
		actual, err := normalizeValue(val)
		if err != nil {
			t.Errorf("%s", err)
		}

		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})

	t.Run("normalizing timespans with minutes only", func(t *testing.T) {
		val := "34m"

		expected := "0.57"
		actual, err := normalizeValue(val)
		if err != nil {
			t.Errorf("%s", err)
		}

		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})
}
