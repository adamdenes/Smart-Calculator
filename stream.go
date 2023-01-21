package main

import (
	"fmt"
)

// TokenStream creates a buffer to hold the tokens
// after tokenization and keeps track of them
type TokenStream struct {
	Tokens []Token // slice of tokens
	Pos    int     // current position in the slice
}

// puts a `Token` back to the buffer
func (ts *TokenStream) putBack() {
	if ts.Pos > 0 {
		ts.Pos--
	}
}

// retrieves a `Token` from the buffer
func (ts *TokenStream) get() *Token {
	if ts.Pos < len(ts.Tokens) {
		t := &ts.Tokens[ts.Pos]
		ts.Pos++
		return t
	}
	return nil
}

// declares a variable
func (ts *TokenStream) declaration(a *Assignment) (int, error) {
	lhs := ts.get() // var name

	lhsErr := checkOperands(lhs, a, LHS)
	if lhsErr != nil {
		//fmt.Println(lhsErr)
		return 0, fmt.Errorf("%w", lhsErr)
	}

	// there is only an LHS
	if len(ts.Tokens) == 1 {
		lookup, lookupErr := a.lookup(lhs.Name)
		if lookupErr != nil {
			return 0, fmt.Errorf("%w", lookupErr)
		}
		//fmt.Println(lookup)
		return lookup, nil
	}

	t2 := ts.get()
	if t2.Kind != '=' {
		//fmt.Println("declaration(): Invalid expression (=)")
		return 0, fmt.Errorf("Invalid expression")
	}

	if len(ts.Tokens) < 3 {
		//fmt.Println("declaration(): Invalid assignment RHS")
		return 0, fmt.Errorf("Invalid assignment")
	}

	rhs := ts.get()
	if rhs.Kind == '-' {
		next := ts.get()
		rhs = next
		rhs.Value *= -1
	}
	rhsErr := checkOperands(rhs, a, RHS)
	if rhsErr != nil {
		//fmt.Println(rhsErr)
		return 0, fmt.Errorf("%w", rhsErr)
	}
	ts.putBack()

	res := ts.expression()
	err := defineVar(lhs.Name, res, a)
	if err != nil {
		//fmt.Println(err)
		return 0, fmt.Errorf("%w", err)
	}
	return res, nil
}

// converts/looks up the variables
func (ts *TokenStream) changeVarsInTokenStream(a *Assignment) {
	for i := range ts.Tokens {
		// NOTE: range loop wouldn't change the values of the slice...
		if num, ok := checkVar(ts.Tokens[i].Name, a); ok {
			//fmt.Printf("checkTokens(): ->%v : %v\n", ts.Tokens[i], num)
			ts.Tokens[i].Value = num
			ts.Tokens[i].Name = ""
		}
		//fmt.Println("changeVarsInTokenStream(): ->", ts.Tokens[i])
	}
}

// statement will evaluate the input, either by processing the expression
// or by doing a declaration
func (ts *TokenStream) statement(a *Assignment) (int, error) {
	ts.changeVarsInTokenStream(a)
	//fmt.Printf("statement(): ->%v\n", ts)
	t := ts.get()
	//fmt.Printf("statement(): ->%v Pos=%v\n", t, ts.Pos)

	switch {
	case t.Name == "":
		ts.putBack()
		return ts.expression(), nil
	default:
		ts.putBack()
		result, err := ts.declaration(a)
		if err != nil {
			return 0, err
		}
		return result, nil
	}
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
		if t.Kind != '+' && t.Kind != '-' {
			ts.putBack()
			//fmt.Printf("EXPRESSION(): putBack-> %v\n", t)
			return left
		}
		right := ts.term()
		if t.Kind == '+' {
			left += right
		} else {
			left -= right
		}
	}
}

// term will provide the primary for the expression
func (ts *TokenStream) term() int {
	left := ts.primary()
	//fmt.Printf("TERM(): Left->\t %v\n", left)

	for {
		t := ts.get()
		//fmt.Printf("TERM(): Token-> %v\n", t)

		if t == nil {
			return left
		}
		switch t.Kind {
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

// primary returns the most basic component, a number
func (ts *TokenStream) primary() int {
	t := ts.get()
	//fmt.Printf("PRIMARY(): Token-> %v\n", t)
	if t == nil {
		return 0
	}
	if t.Kind == '(' {
		result := ts.expression()
		t := ts.get()
		if t.Kind != ')' {
			fmt.Println("Error: expected ')'")
			return 0
		}
		return result
	}
	if t.Kind == '-' {
		return -ts.primary()
	}
	if t.Kind == '+' {
		return ts.primary()
	}
	if t.Value != 0 { // check if the token is a number
		return t.Value
	}
	//fmt.Println("Error: expected number or '('")
	return 0
}
