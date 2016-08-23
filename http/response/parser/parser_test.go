package parser

import (
	"bytes"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	str := `HTTP/1.1 200 OK
Someheader: valua
Content-Length: 14
other-header: test
andthisis:testheader

1234567890
asd
`
	testParse(str, "1234567890\nasd", 14, t)
}

func TestReadEmpty(t *testing.T) {
	str := `HTTP/1.1 200 OK
Date: Sun, 14 Aug 2016 11:44:13 GMT
Content-Length: 0
Content-Type: text/plain, charset=utf8

`
	testParse(str, "", 0, t)
}

func testParse(str, exBody string, exContentLength int, t *testing.T) {
	p := Parser{}
	input := []byte(str)
	for i := 0; i < len(input); i += 10 {
		limit := i
		if i >= len(input) {
			limit = len(input) - 1
		}
		finished, err := p.Parse(input, limit)
		if err != nil {
			t.Error(err)
		}
		if finished {
			if limit < len(input)-1 {
				t.Error("Parsed finished, but not whole input was read")
			}
		}
	}
	if p.ContentLength != exContentLength {
		t.Error(fmt.Sprintf("Expected cl %d, got %d", exContentLength, p.ContentLength))
	}
	body := input[p.BodyStartPosition : p.BodyStartPosition+p.ContentLength]
	if bytes.Compare(body, []byte(exBody)) != 0 {
		t.Error(fmt.Sprintf("Expected body %s, but got '%s'", exBody, body))
	}

}
