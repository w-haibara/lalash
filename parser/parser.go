package parser

import (
	"errors"
	"fmt"
	"strings"
	"text/scanner"
)

const (
	CommandToken      = "command"
	StringToken       = "string"
	RawStringToken    = "raw-string"
	SubstitutionToken = "substitution"
	SplitToken        = "split"
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
				Kind: CommandToken,
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
				Kind: CommandToken,
				Val:  tmp,
			})
			continue
		}
		ret = append(ret, Token{
			Kind: CommandToken,
			Val:  s.TokenText(),
		})
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}

	for i, v := range ret {
		if strings.HasPrefix(v.Val, "#") {
			ret = ret[:i]
		}
	}

	for i := 0; i < len(ret); i++ {
		if strings.HasSuffix(ret[i].Val, ";") {
			ret[i].Val = strings.TrimSuffix(ret[i].Val, ";")

			if ret[i].Val == "" {
				ret[i] = Token{Kind: SplitToken}
				break
			}

			ret = append(ret, Token{})
			copy(ret[i+2:], ret[i+1:])
			ret[i+1] = Token{Kind: SplitToken}
			i -= 1
			continue
		}
	}

	var err error

	ret, err = ParenParser(ret, "{", "}", RawStringToken)
	if err != nil {
		return nil, err
	}

	ret, err = ParenParser(ret, "(", ")", SubstitutionToken)
	if err != nil {
		return nil, err
	}

	for i, v := range ret {
		if v.Kind == RawStringToken {
			continue
		}
		if strings.HasPrefix(v.Val, "\"") && strings.HasSuffix(v.Val, "\"") {
			ret[i].Kind = StringToken
			ret[i].Val = strings.TrimPrefix(v.Val, "\"")
			ret[i].Val = strings.TrimSuffix(ret[i].Val, "\"")
		}
	}

	return ret, nil
}

func ParenParser(tok []Token, start, end, kind string) ([]Token, error) {
	res := []Token{}
	count := 0
	tmp := ""
	for i := 0; i < len(tok); i++ {
		if tok[i].Kind == StringToken || tok[i].Kind == RawStringToken {
			res = append(res, tok[i])
			continue
		}

		b1 := func(s string) bool {
			r1 := false
			r2 := false
			for {
				r1 = strings.HasPrefix(s, start)
				if !r1 {
					break
				}
				r2 = r2 || r1
				count++
				s = strings.TrimPrefix(s, start)
			}
			return r2
		}(tok[i].Val)

		b2 := func(s string) bool {
			r1 := false
			r2 := false
			for {
				r1 = strings.HasSuffix(s, end)
				if !r1 {
					break
				}
				r2 = r2 || r1
				count--
				s = strings.TrimSuffix(s, end)
			}
			return r2
		}(tok[i].Val)

		if b1 || b2 || count > 0 {
			tmp = concat(tmp, tok[i].Val)
		}

		if count == 0 {
			if tmp != "" {
				tmp = strings.TrimPrefix(tmp, start)
				tmp = strings.TrimSuffix(tmp, end)
				tmp = strings.TrimSpace(tmp)
				res = append(res, Token{
					Val:  tmp,
					Kind: kind,
				})
				tmp = ""
				continue
			}
			res = append(res, tok[i])
			continue
		}
	}

	if count > 0 {
		return nil, fmt.Errorf("parentheses not terminated")
	}

	return res, nil
}

func concat(s1, s2 string) string {
	if s2 == "" {
		return s1
	}
	if s1 == "" {
		return s2
	}
	s1 = strings.TrimSpace(s1)
	s2 = strings.TrimSpace(s2)
	return s1 + " " + s2
}
