package services

import "testing"

func TestPrepareServerFiles_RejectsNewlineInjection(t *testing.T) {
	dir := t.TempDir()
	err := PrepareServerFiles(dir, false, true, map[string]string{
		"motd": "hello\nrcon.password=pwned",
	})
	if err == nil {
		t.Fatal("expected an error for a property value containing a newline")
	}
}

func TestPrepareServerFiles_AllowsValidProperties(t *testing.T) {
	dir := t.TempDir()
	if err := PrepareServerFiles(dir, false, true, map[string]string{"motd": "Hello World"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
