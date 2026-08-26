package command

import (
	"path/filepath"
	"testing"

	"github.com/peak/s5cmd/v2/storage/url"
)

func TestJoinLocalDestinationConfinesObjectPath(t *testing.T) {
	t.Parallel()

	root, err := url.New(filepath.Join(t.TempDir(), "downloads"))
	if err != nil {
		t.Fatal(err)
	}

	for _, relative := range []string{"../escape", "nested/../../escape"} {
		if _, err := joinLocalDestination(root, relative); err == nil {
			t.Fatalf("expected %q to be rejected", relative)
		}
	}

	got, err := joinLocalDestination(root, "nested/../object.bin")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root.Absolute(), "object.bin")
	if got.Absolute() != want {
		t.Fatalf("expected %q, got %q", want, got.Absolute())
	}
}

func TestGenerateDestinationURLRejectsEscapingRemoteKey(t *testing.T) {
	t.Parallel()

	src, err := url.New("s3://bucket/prefix/*")
	if err != nil {
		t.Fatal(err)
	}
	if !src.Match("prefix/../escape") {
		t.Fatal("test key did not match source wildcard")
	}
	src = src.Clone()
	src.Path = "prefix/../escape"

	dst, err := url.New(filepath.Join(t.TempDir(), "downloads"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := generateDestinationURL(src, dst, true); err == nil {
		t.Fatal("expected escaping remote key to be rejected")
	}
}
