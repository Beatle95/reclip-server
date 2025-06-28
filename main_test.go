package main

import (
	"testing"
)

func TestArgsParsing(t *testing.T) {
	type testCases struct {
		Args             []string
		ExpectedSettings parsedSettings
	}

	var success_test_cases = []testCases{
		{
			Args:             []string{"-p", "21"},
			ExpectedSettings: parsedSettings{Port: 21},
		},
		{
			Args:             []string{"--port=45"},
			ExpectedSettings: parsedSettings{Port: 45},
		},
	}

	var error_test_cases = [][]string{
		{"-p"},
		{"--port="},
		{"--port=af"},
	}

	for _, elem := range success_test_cases {
		result, err := parseMainArgs(elem.Args)
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
		if result != elem.ExpectedSettings {
			t.Errorf("Expected equality: %v and %v", result, elem.ExpectedSettings)
		}
	}

	for _, elem := range error_test_cases {
		_, err := parseMainArgs(elem)
		if err == nil {
			t.Errorf("Expected an error for input: %s", elem)
		}
	}
}
