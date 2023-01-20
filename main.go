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
)

const (
	LHS = iota
	RHS
)

// Token describes a single character from the input
type Token struct {
	Kind  rune
	Value int
	Name  string
}

func main() {
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
		default:
			// the input is an expression
			list, err := tokenize(in)
			if err != nil {
				fmt.Println(err)
			} else {
				// fill the token stream
				ts := TokenStream{list, 0}
				statement, sErr := ts.statement(&a)
				if sErr != nil {
					fmt.Println(sErr)
				}
				fmt.Println(statement)
			}
		}
	}
}

// validate the left-hand side and right-hand side
func checkOperands(operand *Token, a *Assignment, side int) error {
	switch side {
	case LHS:
		if !onlyLetters([]byte(operand.Name)) {
			//fmt.Println("checkLhs(): Invalid identifier LHS")
			return errors.New("Invalid identifier")
		}
	case RHS:
		if unicode.IsLetter(operand.Kind) {
			if !onlyLetters([]byte(operand.Name)) {
				//fmt.Println("checkRhs(): Invalid assignment RHS")
				return errors.New("Invalid assignment")
			}
			val, lookupErr := a.lookup(operand.Name)
			if lookupErr != nil {
				return lookupErr
			}
			operand.Value = val
			return nil
		}
		if !onlyDigits([]byte(operand.Name)) {
			//fmt.Println("checkRhs(): Invalid assignment RHS")
			return errors.New("Invalid assignment")
		}
	}
	return nil
}

// add a `Variable` to the `Assignment` map
func defineVar(name string, val int, a *Assignment) error {
	err := a.add(Variable{name, strconv.Itoa(val)})
	if err != nil {
		return err
	}
	return nil
}

// determines if the input is a command
func isCmd(b []byte) bool {
	if bytes.HasPrefix(b, []byte("/")) || string(b) == "" {
		return true
	}
	return false
}

// executes command
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

// checks if the bytes only contain letters
func onlyLetters(b ...[]byte) bool {
	for _, r := range bytes.Runes(b[0]) {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// checks if the bytes only contain digits
func onlyDigits(b ...[]byte) bool {
	for _, r := range bytes.Runes(b[0]) {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
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

// checks for digits and letters at the same time
func isMixed(b ...[]byte) bool {
	return onlyDigits(b...) || onlyLetters(b...)
}

// creates a slice of tokens out of the input bytes
func tokenize(b []byte) ([]Token, error) {
	var tokens []Token

	if bytes.Count(b, []byte("=")) > 1 {
		//fmt.Println("makeVar(): Invalid assignment (`=`) > 1")
		return nil, errors.New("Invalid assignment")
	}

	rs := bytes.Runes(b)
	if !unicode.IsDigit(rs[len(rs)-1]) && !isMixed(b) && !bytes.Contains(b, []byte("=")) {
		//fmt.Println("tokenize(): Invalid expression (last rune)")
		return nil, errors.New("Invalid expression")
	}

	// regular expression to match digits, operators, and parentheses
	re := regexp.MustCompile(`\w+|\d+|[+\-*/()%= ]|-?\d+`)
	matches := re.FindAll(b, -1)

	var negative string
	for i, match := range matches {
		s := string(match)
		//fmt.Printf("tokenize(): ->\t%q <-> %s\n", rune(s[0]), s)
		if s == " " {
			continue
		}

		// check if matched string consists of letters
		if onlyLetters(match) {
			//fmt.Printf("tokenize(): onlyLetters ->\t%q\n", match)
			tokens = append(tokens, Token{rune(s[0]), 0, s})

		} else if s == "-" && i == 0 {
			// check if the next rune is number, if not, it is just an operator
			if operator, err := strconv.Atoi(string(matches[i+1])); err != nil {
				tokens = append(tokens, Token{rune(s[0]), operator, s})
			}
			negative = s
		} else if s == "+" && i == 0 {
			// check if next rune is number, if it is
			if _, err := strconv.Atoi(string(matches[i+1])); err == nil {
				continue
			}
		} else if value, err := strconv.Atoi(s); err == nil {
			// check if the rune is a `-`
			if negative == "-" {
				negative += s
				val, _ := strconv.Atoi(negative)
				tokens = append(tokens, Token{rune(s[0]), val, ""})
			} else {
				//check if the rune is a number
				tokens = append(tokens, Token{rune(s[0]), value, ""})
			}
		} else {
			// if the rune is not a number, it must be an operator or a parenthesis
			tokens = append(tokens, Token{rune(s[0]), 0, s})
		}
	}

	// if there is no operator between numbers, it is invalid
	if len(tokens) > 1 &&
		(unicode.IsDigit(tokens[0].Kind) &&
			unicode.IsDigit(tokens[1].Kind)) {
		//fmt.Println("tokenize(): Invalid expression (' ')")
		return nil, errors.New("Invalid expression")
	}

	return tokens, nil
}
