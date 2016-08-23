package parser

import (
	"bytes"
	"errors"
	"fmt"
	"math"
)

var CONTENT_LENGTH_BYTES []byte = []byte("content-length")

const (
	STATUS_LINE int = iota
	HEADER_NAME_STARTED
	HEADER_NAME_FINISHED
	HEADER_VALUE_STARTED
	HEADER_VALUE_FINISHED
	BODY_STARTED
	BODY_FINISHED
)

const (
	HEADER_INVALID int = iota
	HEADER_CONTENT_LENGTH
)

type Parser struct {
	state                     int
	currentHeader             int
	currentTokenStartPosition int
	ContentLength             int
	BodyStartPosition         int
	position                  int
}

func (p *Parser) Parse(buffer []byte, limit int) (bool, error) {
	for ; p.position < limit; p.position++ {
		switch p.state {
		case STATUS_LINE:
			if buffer[p.position] == '\n' {
				p.state = HEADER_NAME_STARTED
				p.currentTokenStartPosition = p.position
			}
		case HEADER_NAME_STARTED:
			if buffer[p.position] == ':' {
				if bytes.Compare(buffer[p.currentTokenStartPosition:p.position], CONTENT_LENGTH_BYTES) == 0 {
					p.currentHeader = HEADER_CONTENT_LENGTH
				}
				p.state = HEADER_VALUE_STARTED
				p.currentTokenStartPosition = p.position + 1
			} else if buffer[p.position] == '\n' {
				if p.ContentLength == 0 {
					p.state = BODY_FINISHED
				} else {
					p.state = BODY_STARTED
					p.BodyStartPosition = p.position + 1
				}
			} else if buffer[p.position] >= 'A' && buffer[p.position] <= 'Z' {
				buffer[p.position] += 'a' - 'A' //change to lowercase
			}
		case HEADER_VALUE_STARTED:
			if p.currentTokenStartPosition == p.position && buffer[p.position] == ' ' {
				p.currentTokenStartPosition++
			} else if buffer[p.position] == '\n' {
				if p.currentHeader == HEADER_CONTENT_LENGTH {
					end := p.position - 1
					if buffer[p.position-1] == '\r' {
						end -= 1
					}
					for i := end; i >= p.currentTokenStartPosition; i-- {
						if buffer[i] < '0' || buffer[i] > '9' {
							return false, errors.New(fmt.Sprintf("Invalid numeric content length value %q", buffer[i]))
						}
						p.ContentLength += int((buffer[i] - '0')) * int(math.Pow10(end-i))
					}
					p.currentHeader = HEADER_INVALID
				}
				p.state = HEADER_NAME_STARTED
				p.currentTokenStartPosition = p.position + 1
			}
		case BODY_STARTED:
			if p.position-p.currentTokenStartPosition >= p.ContentLength {
				p.state = BODY_FINISHED
			}
		}
	}
	return p.state == BODY_FINISHED, nil
}
