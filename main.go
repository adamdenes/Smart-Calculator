package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Assignment map[string]int

type Variable struct {
	name  string
	value string
}

type Token struct {
	kind  rune
	value int
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

func (a *Assignment) add(v Variable) {
	num, err := strconv.Atoi(v.value)

	if err != nil {
		prev := a.lookup(v)
		fmt.Printf("found it! val=%v\n", prev)
		fmt.Println("NUM", prev)
		(*a)[v.name] = prev
	} else if _, ok := (*a)[v.name]; !ok {
		(*a)[v.name] = num
	} else {
		(*a)[v.name] = num
	}
}

func (a *Assignment) lookup(v Variable) int {
	// check if the value is the same as one of the keys
	for key := range *a {
		if v.value == key {
			return (*a)[key]
		}
	}

	return 0
}

func (v *Variable) printVar() {
	fmt.Println(v.value)
}

func main() {
	var list []Token
	a := make(Assignment)

	for {
		in := input()

		// the input is a command
		if bytes.HasPrefix(in, []byte("/")) || string(in) == "" {
			switch string(in) {
			case "/exit":
				fmt.Println("Bye!")
				os.Exit(0)
			case "/help":
				fmt.Println("The program tries to be a simple calculator")
			case "":
				continue
			default:
				fmt.Println("Unknown command")
			}
		} else if bytes.Contains(in, []byte("=")) {
			// the input is an assignment
			v := makeVar(in)
			a.add(v)
			fmt.Printf("a: %v\n", a)
		} else {
			// the input is an expression
			err := sanitize(in)
			if err != nil {
				fmt.Println(err)
				continue
			}
			list = tokenize(in)
			// fill the token stream
			ts := TokenStream{list, 0}

			fmt.Println(ts.expression())
		}
	}
}

func makeVar(b []byte) Variable {
	myVar := Variable{}

	re := regexp.MustCompile(`\w+|-?\d+`)
	matches := re.FindAll(b, -1)
	fmt.Printf("%q\n", matches)

	myVar.name = string(matches[0])
	myVar.value = string(matches[1])

	return myVar
}

func assignment() {

}

// read the input from stdin, line by line
func input() []byte {
	reader := bufio.NewReader(os.Stdin)

	line, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: input: %v", err)
	}

	return bytes.TrimSpace(line)
}

// check for invalid expressions
func sanitize(b []byte) error {
	err := "Invalid expression"

	// err if operator is last or not digit
	lastRune, _ := utf8.DecodeLastRune(b)
	if !unicode.IsDigit(lastRune) {
		return errors.New(err)
	}

	return nil
}

// creates a slice of tokens out of the input bytes
func tokenize(b []byte) []Token {

	var tokens []Token
	// regular expression to match digits, operators, and parentheses
	re := regexp.MustCompile(`\d+|[+\-*/()%]|-?\d+`)
	matches := re.FindAll(b, -1)

	var negative string
	for i, match := range matches {
		s := string(match)
		if s == "-" && i == 0 {
			// check if the next token is number, if not, it is just an operator
			if operator, err := strconv.Atoi(string(matches[i+1])); err != nil {
				tokens = append(tokens, Token{rune(s[0]), operator})
			}
			negative = s
		} else if value, err := strconv.Atoi(s); err == nil {
			// check if the token is a `-`
			if negative == "-" {
				negative += s
				val, _ := strconv.Atoi(negative)
				tokens = append(tokens, Token{rune(s[0]), val})
			} else {
				// check if the token is a number
				tokens = append(tokens, Token{rune(s[0]), value})
			}
		} else {
			// if the token is not a number, it must be an operator or a parenthesis
			tokens = append(tokens, Token{rune(s[0]), 0})
		}
	}
	return tokens
}

/*
Expression Grammar

Expression:
    Term
    Expression '+' Term
    Expression '-' Term

Term:
    Primary
    Term '*' Primary
    Term '/' Primary
    Term '%' Primary

Primary:
    Number
    '('Expression')'

Number:
    integer literal
*/

// expression recursively evaluates the terms
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

// `term` will provide the primary for the expression
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
		case '%':
			left %= ts.primary()
		default:
			ts.putBack()
			return left
		}
	}
}

// `primary` returns the most basic component, a number
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
	//fmt.Println("Error: expected number or '('")
	return 0
}
