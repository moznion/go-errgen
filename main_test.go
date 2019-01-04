package errgen

import (
	"os"
	"testing"

	"github.com/moznion/go-errgen/test"
)

func TestGeneratedBasicErrorMessage(t *testing.T) {
	filePath := "test/basic_err_msg_errmsg_gen.go"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("file doesn't exist: %s", filePath)
	}

	dataset := map[string]error{
		"[ERR-0] this is FOO error":                    test.FooErr(),
		"[ERR-1] this is BAR error [123, hello world]": test.BarErr(123, "hello world"),
	}

	for expected, got := range dataset {
		if got.Error() != expected {
			t.Errorf(`got unexpected result: expected="%s", got="%s"`, expected, got)
		}
	}
}

func TestGeneratedErrorMessageWithPrefix(t *testing.T) {
	filePath := "test/prefix_err_msg_errmsg_gen.go"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("file doesn't exist: %s", filePath)
	}

	dataset := map[string]error{
		"[PREF-0] this is BUZ error":                    test.BuzErr(),
		"[PREF-1] this is QUX error [123, hello world]": test.QuxErr(123, "hello world"),
	}

	for expected, got := range dataset {
		if got.Error() != expected {
			t.Errorf(`got unexpected result: expected="%s", got="%s"`, expected, got)
		}
	}
}

func TestGeneratedErrorMessageWithArbitraryOutputFilePath(t *testing.T) {
	filePath := "test/foobar_errmsg_gen.go"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("file doesn't exist: %s", filePath)
	}

	dataset := map[string]error{
		"[ERR-0] this is FOOBAR error": test.FooBarErr(),
	}

	for expected, got := range dataset {
		if got.Error() != expected {
			t.Errorf(`got unexpected result: expected="%s", got="%s"`, expected, got)
		}
	}
}
