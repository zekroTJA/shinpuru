package multierror

import (
	"errors"
	"strings"
	"testing"
)

var (
	testFormatFunc = func(errors []error) string {
		lines := make([]string, len(errors))
		for i, err := range errors {
			lines[i] = err.Error()
		}
		return strings.Join(lines, "\n")
	}

	testErrors = []error{
		errors.New("error 1"),
		errors.New("error 2"),
		errors.New("error 3"),
		errors.New("error 4"),
	}

	poisonedTestErrors = []error{
		errors.New("error 1"),
		errors.New("error 2"),
		nil,
		errors.New("error 3"),
		nil,
		errors.New("error 4"),
	}
)

func TestError(t *testing.T) {
	err := New()
	err.Append(testErrors...)

	result := err.Error()
	expected := defaultFormatFunc(testErrors)
	if result != expected {
		t.Errorf(
			"result was:\n>\n%s\n<\nbut should have been:\n>\n%s\n<\n",
			result, expected)
	}

	err = New(testFormatFunc)
	err.Append(testErrors...)

	result = err.Error()
	expected = testFormatFunc(testErrors)
	if result != expected {
		t.Errorf(
			"result was:\n>\n%s\n<\nbut should have been:\n>\n%s\n<\n",
			result, expected)
	}

	err = New()
	err.Append(poisonedTestErrors...)

	result = err.Error()
	expected = defaultFormatFunc(testErrors)
	if result != expected {
		t.Errorf(
			"result was:\n>\n%s\n<\nbut should have been:\n>\n%s\n<\n",
			result, expected)
	}

	err = New()
	if err.Error() != "" {
		t.Error("empty error should return an empty string")
	}
}

func TestLen(t *testing.T) {
	err := New()
	err.Append(testErrors...)

	result := err.Len()
	expected := len(testErrors)
	if result != expected {
		t.Errorf("returned len was %d but should have been %d",
			result, expected)
	}
}
