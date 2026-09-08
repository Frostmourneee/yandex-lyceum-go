package main

import (
	"slices"
	"testing"
)

func main() {

}

func TestSortIntegers(t *testing.T) {
	cases := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			"first",
			[]int{3, 2, 1},
			[]int{1, 2, 3},
		},
		{
			"second",
			[]int{1, 2, 3},
			[]int{1, 2, 3},
		},
		{
			"third",
			[]int{3, 2, 1},
			[]int{1, 2, 3},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SortIntegers(tc.input)
			if !slices.Equal(got, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestContains(t *testing.T) {
	cases := []struct {
		name string
		inputArr []int
		inputTarget int
		expected bool
	} {
		{
			"first",
			[]int{1, 2, 3},
			2,
			true,
		},
		{
			"second",
			[]int{1, 2, 3},
			4,
			false,
		},
		{
			"third",
			[]int{1, 2, 3},
			1,
			true,
		},
		{
			"fourth",
			[]int{},
			0,
			false,
		},
		{
			"fifth",
			[]int{1},
			1,
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Contains(tc.inputArr, tc.inputTarget)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestReverseString(t *testing.T) {
	cases := []struct {
		name string
		input string
		expected string
	} {
		{
			"first",
			"hello",
			"olleh",
		},
		{
			"second",
			"",
			"",
		},
		{
			"third",
			"1234567890",
			"0987654321",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ReverseString(tc.input)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestAreAnagrams(t *testing.T) {
	cases := []struct {
		name string
		input1 string
		input2 string
		expected bool
	} {
		{
			"first",
			"listen",
			"silent",
			true,
		},
		{
			"second",
			"hello",
			"world",
			false,
		},
		{
			"third",
			"",
			"",
			true,
		},
		{
			"fourth",
			"abc",
			"cba",
			true,
		},
		{
			"fifth",
			"abc",
			"def",
			false,
		},
		{
			"sixth",
			"abc",
			"abc",
			true,
		},
		{
			"seventh",
			"a",
			"a",
			true,
		},
		{
			"eighth",
			"AbCDEgF",
			"FEDgCBA",
			true,
		},
		{
			"ninth",
			"a",
			"bc",
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := AreAnagrams(tc.input1, tc.input2)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}