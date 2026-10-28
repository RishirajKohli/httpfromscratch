package request

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request-line")
var ERROR_MALFORMED_HTTP_MESSAGE = fmt.Errorf("malformed http message")
var ERROR_INVALID_METHOD_NAME = fmt.Errorf("invalid method name")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")

var validMethodName = regexp.MustCompile(`^[A-Z]+$`)
var validHTTPVersions = map[string]struct{}{
	"1.1": {},
}

type Request struct {
	RequestLine RequestLine
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawMessage, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Unable to read from reader: %w", err)
	}
	rl, _, err := parseRequestLine(string(rawMessage))
	if err != nil {
		return nil, fmt.Errorf("Unable to parse request line: %w", err)
	}
	return &Request{RequestLine: *rl}, nil
}

func parseRequestLine(message string) (*RequestLine, string, error) {
	startLine, remainder, found := strings.Cut(message, "\r\n")
	if !found {
		return nil, message, ERROR_MALFORMED_HTTP_MESSAGE
	}

	startlineParts := strings.Split(startLine, " ")
	if len(startlineParts) != 3 {
		return nil, message, ERROR_MALFORMED_REQUEST_LINE
	}

	//validate METHOD
	if !validMethodName.MatchString(startlineParts[0]) {
		return nil, message, ERROR_INVALID_METHOD_NAME
	}

	httpParts := strings.Split(startlineParts[2], "/")

	httpVersion := httpParts[1]
	//validate HTTP version
	if _, ok := validHTTPVersions[httpVersion]; httpParts[0] != "HTTP" || !ok {
		return nil, message, ERROR_UNSUPPORTED_HTTP_VERSION
	}

	rl := RequestLine{
		Method:        startlineParts[0],
		RequestTarget: startlineParts[1],
		HttpVersion:   httpVersion,
	}
	return &rl, remainder, nil

}
