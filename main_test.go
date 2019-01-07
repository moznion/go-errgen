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
		"[ERR-1] this is FOO error":                    test.FooErr(),
		"[ERR-2] this is BAR error [123, hello world]": test.BarErr(123, "hello world"),
	}

	for expected, got := range dataset {
		if got.Error() != expected {
			t.Errorf(`got unexpected result: expected="%s", got="%s"`, expected, got)
		}
	}

	expected := []string{
		`[ERR-1] this is FOO error`,
		`[ERR-2] this is BAR error [%d, %s]`,
	}
	for i, got := range test.BasicErrMsgList() {
		if exp := expected[i]; got != exp {
			t.Errorf(`got unexpected result: expected="%s", got="%s"`, exp, got)
		}
	}
}

func TestGeneratedErrorMessageWithPrefix(t *testing.T) {
	filePath := "test/prefix_err_msg_errmsg_gen.go"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("file doesn't exist: %s", filePath)
	}

	dataset := map[string]error{
		"[PREF-1] this is BUZ error":                    test.BuzErr(),
		"[PREF-2] this is QUX error [123, hello world]": test.QuxErr(123, "hello world"),
	}

	for expected, got := range dataset {
		if got.Error() != expected {
			t.Errorf(`got unexpected result: expected="%s", got="%s"`, expected, got)
		}
	}

	expected := []string{
		`[PREF-1] this is BUZ error`,
		`[PREF-2] this is QUX error [%d, %s]`,
	}
	for i, got := range test.PrefixErrMsgList() {
		if exp := expected[i]; got != exp {
			t.Errorf(`got unexpected result: expected="%s", got="%s"`, exp, got)
		}
	}
}

func TestGeneratedErrorMessageWithArbitraryOutputFilePath(t *testing.T) {
	filePath := "test/foobar_errmsg_gen.go"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("file doesn't exist: %s", filePath)
	}

	dataset := map[string]error{
		"[ERR-1] this is FOOBAR error": test.FooBarErr(),
	}

	for expected, got := range dataset {
		if got.Error() != expected {
			t.Errorf(`got unexpected result: expected="%s", got="%s"`, expected, got)
		}
	}

	expected := []string{
		"[ERR-1] this is FOOBAR error",
	}
	for i, got := range test.PathSpecifiedErrMsgList() {
		if exp := expected[i]; got != exp {
			t.Errorf(`got unexpected result: expected="%s", got="%s"`, exp, got)
		}
	}
}

func TestGeneratedErrorMessageWithObsoleted(t *testing.T) {
	filePath := "test/obsoleted_err_msg_errmsg_gen.go"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("file doesn't exist: %s", filePath)
	}

	dataset := map[string]error{
		"[ERR-1] this is error 1": test.OneErr(),
		"[ERR-3] this is error 3": test.ThreeErr(),
	}

	for expected, got := range dataset {
		if got.Error() != expected {
			t.Errorf(`got unexpected result: expected="%s", got="%s"`, expected, got)
		}
	}

	expected := []string{
		`[ERR-1] this is error 1`,
		`[ERR-3] this is error 3`,
	}
	for i, got := range test.ObsoletedErrMsgList() {
		if exp := expected[i]; got != exp {
			t.Errorf(`got unexpected result: expected="%s", got="%s"`, exp, got)
		}
	}
}
