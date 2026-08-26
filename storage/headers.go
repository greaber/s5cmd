package storage

import (
	"fmt"
	"net/http"
	"strings"
)

const requestHeaderSeparator = "\x00"

var reservedRequestHeaders = map[string]struct{}{
	"Authorization":        {},
	"Connection":           {},
	"Content-Length":       {},
	"Host":                 {},
	"Transfer-Encoding":    {},
	"X-Amz-Content-Sha256": {},
	"X-Amz-Date":           {},
	"X-Amz-Security-Token": {},
}

// EncodeRequestHeaders stores repeated CLI header values in a comparable
// representation so Options remains usable as a session-cache key.
func EncodeRequestHeaders(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return requestHeaderSeparator + strings.Join(values, requestHeaderSeparator)
}

// ValidateRequestHeaders validates encoded request headers without exposing
// their parsed, mutable representation.
func ValidateRequestHeaders(encoded string) error {
	_, err := parseRequestHeaders(encoded)
	return err
}

func parseRequestHeaders(encoded string) (http.Header, error) {
	headers := make(http.Header)
	if encoded == "" {
		return headers, nil
	}

	if !strings.HasPrefix(encoded, requestHeaderSeparator) {
		return nil, fmt.Errorf("invalid encoded request headers")
	}
	for _, value := range strings.Split(strings.TrimPrefix(encoded, requestHeaderSeparator), requestHeaderSeparator) {
		name, headerValue, ok := strings.Cut(value, ":")
		name = strings.TrimSpace(name)
		if !ok || name == "" {
			return nil, fmt.Errorf("request header %q must use NAME:VALUE syntax", value)
		}
		if !validHeaderName(name) {
			return nil, fmt.Errorf("request header %q has an invalid name", value)
		}

		name = http.CanonicalHeaderKey(name)
		if _, reserved := reservedRequestHeaders[name]; reserved {
			return nil, fmt.Errorf("request header %q is managed by the HTTP transport or AWS signer", name)
		}

		headerValue = strings.TrimSpace(headerValue)
		if !validHeaderValue(headerValue) {
			return nil, fmt.Errorf("request header %q contains an invalid value", name)
		}
		headers.Add(name, headerValue)
	}

	return headers, nil
}

func validHeaderValue(value string) bool {
	for i := 0; i < len(value); i++ {
		c := value[i]
		if (c < ' ' && c != '\t') || c == 0x7f {
			return false
		}
	}
	return true
}

func validHeaderName(name string) bool {
	for i := 0; i < len(name); i++ {
		c := name[i]
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') {
			continue
		}
		switch c {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		default:
			return false
		}
	}
	return true
}
