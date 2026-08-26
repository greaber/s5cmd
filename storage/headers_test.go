package storage

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseRequestHeaders(t *testing.T) {
	t.Parallel()

	encoded := EncodeRequestHeaders([]string{
		"X-Tigris-Consistent:true",
		"X-Test: one",
		"x-test:two",
	})
	got, err := parseRequestHeaders(encoded)
	if err != nil {
		t.Fatal(err)
	}
	want := http.Header{
		"X-Tigris-Consistent": {"true"},
		"X-Test":              {"one", "two"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("request headers differ (-want +got):\n%s", diff)
	}
}

func TestParseRequestHeadersRejectsUnsafeValues(t *testing.T) {
	t.Parallel()

	for _, value := range []string{
		"",
		"missing-colon",
		":value",
		"Bad Header:value",
		"Authorization:value",
		"Host:value",
		"X-Amz-Date:value",
		"X-Test:value\r\nInjected: true",
	} {
		value := value
		t.Run(value, func(t *testing.T) {
			t.Parallel()
			if _, err := parseRequestHeaders(EncodeRequestHeaders([]string{value})); err == nil {
				t.Fatalf("expected %q to be rejected", value)
			}
		})
	}
}
