package sqlscanner

import (
	"errors"
	"io"
	"strings"
)

var (
	errorEndOfFile = errors.New("end of file")
)

// SQLScanner reads reader and extracts sql queries from it
type SQLScanner struct {
	r     io.Reader
	Error error
	buf   []byte
	bufI  int
}

// NewSQLScanner creates new sql scanner with given reader.
// Usage sample:
//
//	 	sqlReader := NewSQLScanner(reader)
//
//		var query string
//		for sqlReader.Next(&query) {
//			// do something with query
// 		}
//
// 		if sqlReader.Error != nil {
//			// handle error
//		}
func NewSQLScanner(r io.Reader) SQLScanner {
	return SQLScanner{
		r: r,
	}
}

// Next reads reader and finds next query. Better to use it in a loop.
// `sqlResult` will be filled with the next query.
// Returns `true` if query found and written to the `sqlResult`.
// Returns `false` when reader ended or when error occures.
func (s *SQLScanner) Next(sqlResult *string) bool {
	defer func() {
		if s.Error == errorEndOfFile || s.Error == io.EOF {
			s.Error = nil
		}
	}()

	s.Error = s.skipTillBeginning()
	if s.Error != nil {
		return false
	}

	result := ""

	for {
		p, err := s.indexOf(";")
		if err != nil {
			s.Error = err
			return false
		}

		if p < 0 {
			return false
		}

		d, err := s.pop(p + 1)
		if err != nil {
			s.Error = err
			return false
		}

		result += string(d)

		count := strings.Count(result, "'")
		if count == 0 || count%2 == 0 {
			break
		}

	}

	*sqlResult = result

	return true
}

func (s *SQLScanner) skipTillBeginning() error {
	for {
		d, err := s.peek(2)
		if err != nil {
			return err
		}

		if len(d) < 2 {
			return errorEndOfFile
		}

		switch string(d[0]) {
		case "\n", "\t", " ":
			s.skip(0, 1)
			continue
		}

		if string(d) == "--" {
			s.skip(0, 2)
			s.skipUntilFound("\n")
			continue
		}

		return nil
	}
}

func (s *SQLScanner) peek(l int) ([]byte, error) {
	if len(s.buf) < l {
		r, err := s.read(l)
		if err != nil && err != errorEndOfFile {
			return nil, err
		}
		s.buf = append(s.buf, r...)
	}

	maxLen := len(s.buf) - s.bufI
	if maxLen > l {
		maxLen = l
	}

	return s.buf[s.bufI:maxLen], nil
}

func (s *SQLScanner) pop(l int) ([]byte, error) {
	if len(s.buf) < l {
		r, err := s.read(l)
		if err != nil && err != errorEndOfFile {
			return nil, err
		}
		s.buf = append(s.buf, r...)
	}

	maxLen := len(s.buf) - s.bufI
	if maxLen > l {
		maxLen = l
	}

	result := s.buf[s.bufI:maxLen]

	s.buf = s.buf[s.bufI+maxLen:]
	s.bufI = 0
	return result, nil
}

func (s *SQLScanner) skip(start, l int) {
	s.buf = append(s.buf[0:start], s.buf[l:]...)
}

func (s *SQLScanner) indexOf(str string) (int, error) {
	i := 0
	readSize := 0
	for {
		readSize = i + len(str)*2
		d, err := s.peek(readSize)
		if err != nil {
			return -1, err
		}

		if len(d) == i {
			return -1, nil
		}

		p := strings.Index(string(d), str)
		if p >= 0 {
			return p, nil
		}
		i = readSize
	}
}

func (s *SQLScanner) skipUntilFound(str string) error {
	for {
		d, err := s.peek(len(str) * 2)
		if err != nil {
			return err
		}
		p := strings.Index(string(d), str)
		if p < 0 {
			s.skip(0, len(d))
			continue
		} else if p == 0 {
			s.skip(0, 1)
		} else {
			s.skip(0, p)
		}
		return nil
	}
}

func (s *SQLScanner) read(l int) ([]byte, error) {
	b := make([]byte, l)
	n, err := s.r.Read(b)
	if err != nil {
		return nil, err
	}

	if n < 1 {
		return nil, errorEndOfFile
	}

	return b, nil
}
