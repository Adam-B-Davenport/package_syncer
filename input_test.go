package main

import "testing"

func runRangeTest(t *testing.T, exp []int, actual []int) {
	if len(exp) != len(actual) {
		t.Errorf("Expected length of %d got %d", len(exp), len(actual))
	}
	for i, expVal := range exp {
		if expVal != actual[i] {
			t.Errorf("Expected %d got %d", exp, actual[i])
		}
	}
}

func TestComplex(t *testing.T) {

	exp := []int{1, 2, 3, 4, 5}
	actual := make([]int, 0)
	AppendRange(&actual, 1, 5)
	runRangeTest(t, exp, actual)
}

func TestRange1(t *testing.T) {
	exp := []int{1, 2, 3, 4, 5}
	actual := make([]int, 0)
	err := ParseRange(&actual, "1-5")
	if err != nil {
		t.Errorf("Encountered error while parsing range.\n%s", err)
	}
	runRangeTest(t, exp, actual)
}

func TestRange2(t *testing.T) {
	exp := []int{9, 10, 11}
	actual := make([]int, 0)
	err := ParseRange(&actual, "9- 11")
	if err != nil {
		t.Errorf("Encountered error while parsing range.\n%s", err)
	}
	runRangeTest(t, exp, actual)
}

func TestRange3(t *testing.T) {
	exp := []int{7}
	actual := make([]int, 0)
	err := ParseRange(&actual, "7 - 7")
	if err != nil {
		t.Errorf("Encountered error while parsing range.\n%s", err)
	}
	runRangeTest(t, exp, actual)
}
