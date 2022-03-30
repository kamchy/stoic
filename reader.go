package stoic

import (
	"bufio"
	"io"
	"log"
	"strings"
)

// Quote represents single quote
type Quote struct {
	Text   string
	Author string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func nonEmptyLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	res := make([]string, 0)
	for scanner.Scan() {
		t := scanner.Text()
		if s := strings.Trim(t, " \t\n"); s != "" {
			res = append(res, s)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return res, err
	}

	return res, nil
}

// ReadQuotes reads lines using io.Reader and returns
// slice of Quotes and error
func ReadQuotes(from io.Reader) ([]Quote, error) {
	bufreader := bufio.NewReader(from)
	lines, err := nonEmptyLines(bufreader)
	qs := make([]Quote, 0)
	if err != nil {
		return qs, err
	}
	for i := 0; i < len(lines)/2; i += 2 {
		q := strings.Split(lines[i], "\"")[1]
		a := strings.TrimLeft(lines[i+1], "-–")
		a = strings.TrimRight(a, ".")
		qs = append(qs, Quote{q, a})
	}
	return qs, nil
}
