package internal

import (
	"testing"
)

func TestArgsParsing(t *testing.T) {
	type testCases struct {
		Args             []string
		ExpectedSettings AppSettings
	}

	var success_test_cases = []testCases{
		{
			Args:             []string{"-p", "21"},
			ExpectedSettings: AppSettings{Port: 21},
		},
		{
			Args:             []string{"--port=45"},
			ExpectedSettings: AppSettings{Port: 45},
		},
		{
			Args:             []string{"--app-data-dir=/user/path"},
			ExpectedSettings: AppSettings{Port: DefaultServerPort, AppDataDir: "/user/path"},
		},
		{
			Args:             []string{"--port=8888", "--app-data-dir=/user/path"},
			ExpectedSettings: AppSettings{Port: 8888, AppDataDir: "/user/path"},
		},
	}

	var error_test_cases = [][]string{
		{"-p"},
		{"--port="},
		{"--port=af"},
	}

	for _, elem := range success_test_cases {
		result, err := ParseMainArgs(elem.Args)
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
		if result != elem.ExpectedSettings {
			t.Errorf("Expected equality: %v and %v", result, elem.ExpectedSettings)
		}
	}

	for _, elem := range error_test_cases {
		_, err := ParseMainArgs(elem)
		if err == nil {
			t.Errorf("Expected an error for input: %s", elem)
		}
	}
}
