package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	s := MakeCurlCommand(
		"POST",
		"https://example.com/api/users",
		"Content-Type: application/json\nAuthorization: Bearer abc123\n",
		`{"name":"John Doe","email":"johndoe@example.com","password":"123456"}`,
	)

	fmt.Println(s)
}

func BuildHTTPRequest(method, url, host, headers, body string) string {
	jsonSep := "\r\n"

	builder := strings.Builder{}
	line1 := fmt.Sprintf("%s %s HTTP/1.1%s", method, url, jsonSep)
	builder.WriteString(line1)

	line2 := fmt.Sprintf("Host: %s%s", host, jsonSep)
	builder.WriteString(line2)

	splitedHeaders := strings.SplitSeq(headers, "jsonSep")
	for line := range splitedHeaders {
		builder.WriteString(line)
		builder.WriteString(jsonSep)
	}

	builder.WriteString(body)

	return builder.String()
}

func BuildHTTPResponse(statusLine, headers, body string) string {
	jsonSep := "\r\n"

	builder := strings.Builder{}
	l1 := fmt.Sprintf("%s%s", statusLine, jsonSep)
	builder.WriteString(l1)

	for line := range strings.SplitSeq(headers, jsonSep) {
		l := fmt.Sprintf("%s%s", line, jsonSep)
		builder.WriteString(l)
	}

	builder.WriteString(body)

	return builder.String()
}

func ParseHTTPStatus(s string) (code int, reason string) {
	spaceParsed := strings.SplitN(s, " ", 3)
	val, _ := strconv.Atoi(spaceParsed[1])
	return val, spaceParsed[2]
}

func MakeCurlCommand(method, url, headers, body string) string {
	builder := &strings.Builder{}
	builder.WriteString("curl ")
	if method != "GET" {
		fmt.Fprintf(builder, "-X %s ", method)
	}

	for line := range strings.SplitSeq(headers, "\n") {
		if line == "" {
			break
		}
		fmt.Fprintf(builder, "-H '%s' ", line)
	}

	if body != "" {
		builder.WriteString("--data ")
		fmt.Fprintf(builder, "'%s' ", body)
	}

	builder.WriteString(url)

	return builder.String()
}
