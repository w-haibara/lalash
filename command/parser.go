package command

import (
	"errors"
	"fmt"
	"strings"
	"text/scanner"
)

func Parse(expr string) ([]string, error) {
	ret := []string{}
	errs := []string{}
	var s scanner.Scanner
	s.Init(strings.NewReader(expr))
	s.Whitespace ^= 1 << ' '
	s.Error = func(s *scanner.Scanner, msg string) {
		errs = append(errs, fmt.Sprintf("%s %s", s.Pos(), msg))
	}

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if s.TokenText() == " " {
			tmp := ""
			for {
				tok := s.Peek()
				if tok == ' ' {
					s.Scan()
					continue
				}
				break
			}
			for {
				tok := s.Peek()
				if tok == ' ' || tok == scanner.EOF {
					break
				}
				s.Scan()
				tmp += s.TokenText()
			}
			if tmp == "" {
				continue
			}
			ret = append(ret, tmp)
			continue
		}
		ret = append(ret, s.TokenText())
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	for i, v := range ret {
		if strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"") {
			ret[i] = strings.TrimPrefix(v, "\"")
			ret[i] = strings.TrimSuffix(ret[i], "\"")
		}
	}
	return ret, nil
}
