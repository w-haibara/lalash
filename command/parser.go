package command

import (
	"errors"
	"fmt"
	"strings"
	"text/scanner"
)

const (
	StringToken       = "string"
	SubstitutionToken = "substitution"
)

type Token struct {
	Kind string
	Val  string
}

func Parse(expr string) ([]Token, error) {
	ret := []Token{}
	errs := []string{}
	var s scanner.Scanner
	s.Init(strings.NewReader(expr))
	s.Whitespace ^= 1 << ' '
	s.Error = func(s *scanner.Scanner, msg string) {
		errs = append(errs, fmt.Sprintf("%s %s", s.Pos(), msg))
	}

	first := true
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if first {
			first = false
			tmp := s.TokenText()
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
			ret = append(ret, Token{
				Kind: StringToken,
				Val:  tmp,
			})
			continue
		}

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
			ret = append(ret, Token{
				Kind: StringToken,
				Val:  tmp,
			})
			continue
		}
		ret = append(ret, Token{
			Kind: StringToken,
			Val:  s.TokenText(),
		})
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	for i, v := range ret {
		if strings.HasPrefix(v.Val, "\"") && strings.HasSuffix(v.Val, "\"") {
			ret[i].Val = strings.TrimPrefix(v.Val, "\"")
			ret[i].Val = strings.TrimSuffix(ret[i].Val, "\"")
		}
	}

	for i, v := range ret {
		if strings.HasPrefix(v.Val, "`") && strings.HasSuffix(v.Val, "`") {
			ret[i].Val = strings.TrimPrefix(v.Val, "`")
			ret[i].Val = strings.TrimSuffix(ret[i].Val, "`")
			ret[i].Kind = SubstitutionToken
		}
	}

	return ret, nil
}
