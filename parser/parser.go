package parser

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

	ret, err := ParenParser(ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func ParenParser(ret []Token) ([]Token, error) {
	tokens := []Token{}
	count := 0
	tmp := ""
	for i := 0; i < len(ret); i++ {
		b1 := func(s string) bool {
			r1 := false
			r2 := false
			for {
				r1 = strings.HasPrefix(s, "(")
				if !r1 {
					break
				}
				r2 = r2 || r1
				count++
				s = strings.TrimPrefix(s, "(")
			}
			return r2
		}(ret[i].Val)

		b2 := func(s string) bool {
			r1 := false
			r2 := false
			for {
				r1 = strings.HasSuffix(s, ")")
				if !r1 {
					break
				}
				r2 = r2 || r1
				count--
				s = strings.TrimSuffix(s, ")")
			}
			return r2
		}(ret[i].Val)

		if b1 || b2 || count > 0 {
			tmp = concat(tmp, ret[i].Val)
		}

		if count == 0 {
			if tmp != "" {
				tmp = strings.TrimPrefix(tmp, "(")
				tmp = strings.TrimSuffix(tmp, ")")
				tmp = strings.TrimSpace(tmp)
				tokens = append(tokens, Token{
					Val:  tmp,
					Kind: SubstitutionToken,
				})
				tmp = ""
				continue
			}
			tokens = append(tokens, ret[i])
			continue
		}
	}

	if count > 0 {
		return nil, fmt.Errorf("parentheses not terminated")
	}

	return tokens, nil
}

func concat(s1, s2 string) string {
	if s2 == "" {
		return s1
	}
	if s1 != "" {
		s1 += " "
	}
	return s1 + s2
}
