package main

import (
	"fmt"
	"strings"
	"text/scanner"
)

func Parse(expr string) ([]string, error) {
	ret := []string{}
	var s scanner.Scanner
	s.Init(strings.NewReader(expr))
	s.Whitespace ^= 1 << ' '
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if s.TokenText() == " " {
			tmp := ""
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
	for i, v := range ret {
		fmt.Println(i, v)
		if strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"") {
			ret[i] = strings.TrimPrefix(v, "\"")
			ret[i] = strings.TrimSuffix(ret[i], "\"")
		}
	}
	return ret, nil
}
