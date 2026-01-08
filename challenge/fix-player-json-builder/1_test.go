package main

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	type testCase struct {
		name     string
		level    int
		online   bool
		expected string
	}

	runCases := []testCase{
		{"Alice", 5, true, "{\"name\":\"Alice\",\"level\":5,\"online\":true}"},
		{"Bob", 1, false, "{\"name\":\"Bob\",\"level\":1,\"online\":false}"},
		{"Clara", 10, true, "{\"name\":\"Clara\",\"level\":10,\"online\":true}"},
	}

	submitCases := append(runCases, []testCase{
		{"Zero", 0, false, "{\"name\":\"Zero\",\"level\":0,\"online\":false}"},
		{"HighLevel", 99, true, "{\"name\":\"HighLevel\",\"level\":99,\"online\":true}"},
	}...)

	testCases := runCases
	if withSubmit {
		testCases = submitCases
	}

	skipped := len(submitCases) - len(testCases)

	passCount := 0
	failCount := 0

	for _, test := range testCases {
		result := buildPlayerJSON(test.name, test.level, test.online)
		if result != test.expected {
			failCount++
			t.Errorf(`---------------------------------
Input:
  name:   %q
  level:  %d
  online: %t

Expected: %s
Actual:   %s
Fail
`, test.name, test.level, test.online, test.expected, result)
		} else {
			passCount++
			fmt.Printf(`---------------------------------
Input:
  name:   %q
  level:  %d
  online: %t

Expected: %s
Actual:   %s
Pass
`, test.name, test.level, test.online, test.expected, result)
		}
	}

	fmt.Println("---------------------------------")
	if skipped > 0 {
		fmt.Printf("%d passed, %d failed, %d skipped\n", passCount, failCount, skipped)
	} else {
		fmt.Printf("%d passed, %d failed\n", passCount, failCount)
	}
}

// withSubmit is set at compile time depending
// on which button is used to run the tests
var withSubmit = true
