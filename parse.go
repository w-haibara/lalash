package main

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
			if s.Peek() == ' ' {
				continue
			}
			for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
				tmp += s.TokenText()
				if s.Peek() == ' ' {
					break
				}
			}
			if tmp == " " {
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
