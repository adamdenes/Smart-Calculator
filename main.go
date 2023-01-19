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

func (a *Assignment) add(v Variable) error {
	// assign value of another Variable
	if _, ok := (*a)[v.value]; ok {
		(*a)[v.name] = (*a)[v.value]
		return nil
	}
	num, err := strconv.Atoi(v.value)
	if err != nil {
		return errors.New("add(): Invalid assignment")
	}

	// create / update var
	if _, ok := (*a)[v.name]; !ok {
		(*a)[v.name] = num
	} else {
		(*a)[v.name] = num
	}

	return nil
}

func (a *Assignment) lookup(key string) error {
	if len(*a) == 0 {
		//	fmt.Println("lookup(): len(a) =", len(*a))
		return nil
	}

	if content, ok := (*a)[key]; ok {
		fmt.Println(content)
		return nil
	}
	return errors.New("Unknown variable")
}

func main() {
	var list []Token
	a := make(Assignment)

	for {
		in := input()

		switch {
		case string(in) == "":
			// no input, go next
			continue
		case isCmd(in):
			// the input is a command
			doCmd(in)
		case isAssignment(in):
			// the input is an assignment
			doAssign(in, a)
		default:
			// the input is an expression
			list = tokenize(in)
			if list != nil {
				// fill the token stream
				ts := TokenStream{list, 0}
				fmt.Println(ts.expression())
			}
		}
	}
}

func isCmd(b []byte) bool {
	if bytes.HasPrefix(b, []byte("/")) || string(b) == "" {
		return true
	}
	return false
}

func doCmd(b []byte) {
	switch string(b) {
	case "/exit":
		fmt.Println("Bye!")
		os.Exit(0)
	case "/help":
		fmt.Println("The program tries to be a simple calculator")
	default:
		fmt.Println("Unknown command")
	}
}

func isAssignment(b []byte) bool {
	if bytes.Contains(b, []byte("=")) &&
		unicode.IsLetter(bytes.Runes(b)[0]) {
		return true
	}
	return false
}

func doAssign(b []byte, m Assignment) {
	v, assignErr := makeVar(b)

	if assignErr != nil {
		fmt.Println(assignErr)
		return
	}

	err := m.lookup(v.name)
	if err != nil {
		fmt.Println(err)
		return
	}

	addErr := m.add(v)
	if addErr != nil {
		fmt.Println(addErr)
		return
	}
}

func makeVar(b []byte) (Variable, error) {
	if bytes.Count(b, []byte("=")) > 1 {
		return Variable{}, errors.New("makeVar(): Invalid assignment")
	}

	re := regexp.MustCompile(`\w+|-?\d+`)
	matches := re.FindAll(b, -1)

	var lhs string
	if len(matches) != 0 {
		lhs = string(matches[0])
	}

	for _, r := range lhs {
		if !unicode.IsLetter(r) {
			return Variable{}, errors.New("makeVar(): Invalid identifier")
		}
	}

	var rhs string
	if len(matches) == 1 {
		rhs = string(matches[0])
	} else {
		rhs = string(matches[1])
	}

	return Variable{
		name:  lhs,
		value: rhs,
	}, nil
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

// creates a slice of tokens out of the input bytes
func tokenize(b []byte) []Token {
	lastRune, _ := utf8.DecodeLastRune(b)
	if !unicode.IsDigit(lastRune) {
		fmt.Println("Invalid expression")
		return nil
	}

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
