package main

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain_handler_StatusCode(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler(w, r)
	rw := w.Result()
	defer rw.Body.Close()

	expected := http.StatusOK
	actual := rw.StatusCode
	if actual != expected {
		t.Errorf(`unexpected status code: expected: "%d" actual: "%d"`, expected, actual)
	}
}

func TestMain_handler_ResponseBody(t *testing.T) {
	cases := map[string]struct {
		seed      int64
		nameParam string
		expected  string
	}{
		"KYOU":       {seed: 0, nameParam: "", expected: "{\"name\":\"Gopher\",\"fortune\":\"凶\"}\n"},
		"DAIKYOU":    {seed: 1, nameParam: "", expected: "{\"name\":\"Gopher\",\"fortune\":\"大凶\"}\n"},
		"SUEKICHI":   {seed: 2, nameParam: "", expected: "{\"name\":\"Gopher\",\"fortune\":\"末吉\"}\n"},
		"KICHI":      {seed: 3, nameParam: "", expected: "{\"name\":\"Gopher\",\"fortune\":\"吉\"}\n"},
		"CHUKICHI":   {seed: 4, nameParam: "", expected: "{\"name\":\"Gopher\",\"fortune\":\"中吉\"}\n"},
		"SHOKICHI":   {seed: 5, nameParam: "", expected: "{\"name\":\"Gopher\",\"fortune\":\"小吉\"}\n"},
		"DAICHIKI":   {seed: 9, nameParam: "", expected: "{\"name\":\"Gopher\",\"fortune\":\"大吉\"}\n"},
		"name param": {seed: 9, nameParam: "hioki-daichi", expected: "{\"name\":\"hioki-daichi\",\"fortune\":\"大吉\"}\n"},
	}

	for n, c := range cases {
		c := c
		t.Run(n, func(t *testing.T) {
			rand.Seed(c.seed)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			if c.nameParam != "" {
				q := r.URL.Query()
				q.Add("name", c.nameParam)
				r.URL.RawQuery = q.Encode()
			}

			handler(w, r)
			rw := w.Result()
			defer rw.Body.Close()

			b, err := ioutil.ReadAll(rw.Body)
			if err != nil {
				t.Fatalf("err %s", err)
			}

			expected := c.expected
			actual := string(b)
			if actual != expected {
				t.Errorf(`unexpected response body: expected: "%s" actual: "%s"`, expected, actual)
			}
		})
	}
}

func TestMain_handler_DuringTheNewYear(t *testing.T) {
	isDuringTheNewYearFunc = func() bool {
		return true
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler(w, r)
	rw := w.Result()
	defer rw.Body.Close()

	b, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		t.Fatalf("err %s", err)
	}

	expected := "{\"name\":\"Gopher\",\"fortune\":\"大吉\"}\n"
	actual := string(b)
	if actual != expected {
		t.Errorf(`unexpected response body: expected: "%s" actual: "%s"`, expected, actual)
	}
}

func TestMain_handler_ToJSONError(t *testing.T) {
	toJSONFunc = func(v interface{}) (string, error) {
		return "", errors.New("error in TestMain_handler_ToJSONError")
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler(w, r)
	rw := w.Result()
	defer rw.Body.Close()

	b, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		t.Fatalf("err %s", err)
	}

	expected := "Internal Server Error\n"
	actual := string(b)
	if actual != expected {
		t.Errorf(`unexpected response body: expected: "%s" actual: "%s"`, expected, actual)
	}
}
