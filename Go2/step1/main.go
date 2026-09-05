package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

func main() {
	s := "привет!"
	fmt.Println(len(s))
	fmt.Println(utf8.RuneCountInString(s))
	fmt.Printf("%01x !\n", []byte(s))
}

func WriteString(s string, w io.Writer) error {
	_, err := w.Write([]byte(s))
	if err != nil {
		return err
	}

	return nil
}

func ReadString(r io.Reader) (string, error) {
	data := make([]byte, 100)
	n, err := r.Read(data)
	if err == io.EOF {
		return string(data[:n]), nil
	}
	if err != nil {
		return "", err
	}

	return string(data[:n]), nil
}

type UpperWriter struct {
	UpperString string
}

func (uw *UpperWriter) Write(p []byte) (int, error) {
	s := string(p)
	us := strings.ToUpper(s)
	uw.UpperString = us
	return len(us), nil
}

func Copy(r io.Reader, w io.Writer, n uint) error {
	tmp := make([]byte, n)
	readBytes, err := r.Read(tmp)

	if err != nil && err != io.EOF {
		return err
	}

	w.Write(tmp[:readBytes])
	return nil
}

func Contains(r io.Reader, seq []byte) (bool, error) {
	tmp := make([]byte, 100)
	n, err := r.Read(tmp)

	if err != nil && err != io.EOF {
		return false, err
	}

	return bytes.Contains(tmp[:n], seq), nil
}
