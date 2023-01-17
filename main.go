package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Token struct {
	kind  rune
	value int
	next  *Token // points to the next token
}

func (t *Token) getToken() *Token {
	return t.next
}

type TokenStream struct {
	tokens []Token // slice of tokens
	pos    int     // current position in the slice
}

func (ts *TokenStream) putBack() {
	if ts.pos > 0 {
		ts.pos--
	}
}

func (ts *TokenStream) get() *Token {
	if ts.pos < len(ts.tokens) {
		t := &ts.tokens[ts.pos]
		ts.pos++
		return t
	}
	return nil
}
func main() {
	// TODO: Expression - Term - Primary - Number
	for {
		in := input()
		switch string(in) {
		case "/exit":
			fmt.Println("Bye!")
			os.Exit(0)
		case "/help":
			fmt.Println("The program tries to be a simple calculator")
		case "":
			continue
		default:
			list := tokenize(in)
			// fill the token stream with the first token
			ts := TokenStream{list, 0}
			fmt.Println(ts.expression())
		}
	}
}

func input() []byte {
	reader := bufio.NewReader(os.Stdin)

	line, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: input: %v", err)
	}

	return bytes.TrimSpace(line)
}

func tokenize(b []byte) []Token {

	var tokens []Token
	// regular expression to match digits, operators, and parentheses
	re := regexp.MustCompile(`\d+|[+\-*/()]|-?\d+`)
	matches := re.FindAll(b, -1)

	var negative string
	for i, match := range matches {
		s := string(match)
		if s == "-" && i == 0 {
			// check if the next token is number, if not, it is just an operator
			if operator, err := strconv.Atoi(string(matches[i+1])); err != nil {
				tokens = append(tokens, Token{rune(s[0]), operator, nil})
			}
			negative = s
		} else if value, err := strconv.Atoi(s); err == nil {
			// check if the token is a `-`
			if negative == "-" {
				negative += s
				val, _ := strconv.Atoi(negative)
				tokens = append(tokens, Token{rune(s[0]), val, nil})
			} else {
				// check if the token is a number
				tokens = append(tokens, Token{rune(s[0]), value, nil})
			}

		} else {
			// if the token is not a number, it must be an operator or a parenthesis
			tokens = append(tokens, Token{rune(s[0]), 0, nil})
		}
	}
	return tokens
}

func (ts *TokenStream) expression() int {
	left := ts.term()
	//fmt.Printf("EXPRESSION(): Left->\t %v\n", left)
	for {
		t := ts.get()
		//fmt.Printf("EXPRESSION(): Token-> %v\n", t)

		if t == nil {
			return left
		}
		if t.kind != '+' && t.kind != '-' {
			ts.putBack()
			//fmt.Printf("EXPRESSION(): putBack-> %v\n", t)
			return left
		}
		right := ts.term()
		if t.kind == '+' {
			left += right
		} else {
			left -= right
		}
	}
}

func (ts *TokenStream) term() int {
	left := ts.primary()
	//fmt.Printf("TERM(): Left->\t %v\n", left)

	for {
		t := ts.get()
		//fmt.Printf("TERM(): Token-> %v\n", t)

		if t == nil {
			return left
		}
		switch t.kind {
		case '*':
			left *= ts.primary()
		case '/':
			left /= ts.primary()
		default:
			ts.putBack()
			return left
		}
	}
}

func (ts *TokenStream) primary() int {
	t := ts.get()
	//fmt.Printf("PRIMARY(): Token-> %v\n", t)
	if t == nil {
		return 0
	}
	if t.kind == '(' {
		result := ts.expression()
		t := ts.get()
		if t.kind != ')' {
			fmt.Println("Error: expected ')'")
			return 0
		}
		return result
	}
	if t.kind == '-' {
		return -ts.primary()
	}
	if t.kind == '+' {
		return ts.primary()
	}
	if t.value != 0 { // check if the token is a number
		return t.value
	}
	fmt.Println("Error: expected number or '('")
	return 0
}