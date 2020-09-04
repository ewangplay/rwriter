package rwriter

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestRotateWriterWithoutConfig(t *testing.T) {
	var err error
	var w io.Writer

	w, err = NewRotateWriter(nil)
	if err != nil {
		t.Fatal(err)
	}

	verifyRotateWriter(t, w, "test.log")
}

func TestRotateWriter(t *testing.T) {
	var err error
	var w io.Writer

	cfg := &Config{
		Module:      "kitty",
		Path:        ".",
		MaxSize:     10,
		RotateDaily: true,
	}
	w, err = NewRotateWriter(cfg)
	if err != nil {
		t.Fatalf("Create rotate writer failed: %v\n", err)
	}

	verifyRotateWriter(t, w, "kitty.log")
}

func verifyRotateWriter(t *testing.T, w io.Writer, filename string) {
	log.SetOutput(w)
	msg := "Hello, rotate writer!"
	log.Println(msg)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), msg) {
		t.Fatalf("The log file %s should contains '%s'", filename, msg)
	}
	os.Remove(filename)
}
