package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/quick"
)

func TestEcho(t *testing.T) {
	conv, err := makeStringConverter([]string{"cat"})
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(conv)
	defer ts.Close()
	f := func(x string) bool {
		resp, err := http.Post(ts.URL, "", strings.NewReader(x))
		if err != nil {
			t.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		return bytes.Equal([]byte(x), body)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestError(t *testing.T) {
	errStr := "Who knew?"
	var conv stringConverter = func(x string) (string, error) {
		return "", errors.New(errStr)
	}
	ts := httptest.NewServer(conv)
	defer ts.Close()
	resp, err := http.Post(ts.URL, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 500 {
		t.Error("Didn't have 500 status code")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte(errStr), body) {
		t.Error("Didn't have expected error body")
	}
}
