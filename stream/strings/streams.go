package strings

import (
	"io"
	"strings"

	"github.com/sbchaos/consume/stream"
)

type StringStream struct {
	content string
	len     int
	offset  int
}

func NewStringStream(content string) stream.SimpleStream[rune] {
	return &StringStream{
		content: content,
		len:     len(content),
		offset:  0,
	}
}

func (s *StringStream) Peek() (rune, error) {
	if s.offset < s.len {
		return rune(s.content[s.offset]), nil
	}

	return 0, io.EOF
}

func (s *StringStream) Offset() int {
	return s.offset
}

func (s *StringStream) Seek(n int) {
	s.offset = n
}

func (s *StringStream) Take() (rune, error) {
	if s.offset < s.len {
		val := s.content[s.offset]
		s.offset++
		return rune(val), nil
	}

	return 0, io.EOF
}

func (s *StringStream) TakeN(num int) ([]rune, error) {
	if num == 0 {
		return []rune(""), nil
	}
	if s.offset == s.len {
		return nil, io.EOF
	}
	offset := s.offset
	if s.offset+num > s.len {
		s.offset = s.len
		return []rune(s.content[offset:]), nil
	}

	s.offset = s.offset + num
	return []rune(s.content[offset : offset+num]), nil
}

// TakeWhile will Extract chunk of the stream taking tokens while the supplied
// predicate returns 'True'. Return the chunk and the rest of the stream.
func (s *StringStream) TakeWhile(p stream.Predicate[rune], escape rune) []rune {
	end := -1
	offset := s.offset
	for i, ch := range s.content[s.offset:] {
		if p(ch) {
			continue
		}
		if escape > 0 && i > 0 && rune(s.content[s.offset+i-1]) == escape {
			continue
		}

		end = i
		break
	}
	if end == 0 {
		return nil
	}
	if end == -1 {
		s.offset = s.len
	} else {
		s.offset = offset + end
	}
	return []rune(s.content[offset:s.offset])
}

func (s *StringStream) TakeUntil(str []rune) []rune {
	if len(str) == 0 {
		return nil
	}

	idx := strings.Index(s.content[s.offset:], string(str))
	if idx < 1 {
		return nil
	}

	val := s.content[s.offset : s.offset+idx]
	s.offset += idx
	return []rune(val)
}
