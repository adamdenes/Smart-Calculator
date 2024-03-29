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
		return 0, fmt.Errorf("Invalid expression")
	}

	if len(ts.Tokens) < 3 {
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
		return 0, fmt.Errorf("%w", rhsErr)
	}
	ts.putBack()

	res := ts.expression()
	err := defineVar(lhs.Name, res, a)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return res, nil
}

// converts/looks up the variables
func (ts *TokenStream) changeVarsInTokenStream(a *Assignment) {
	for i := range ts.Tokens {
		if num, ok := checkVar(ts.Tokens[i].Name, a); ok {
			ts.Tokens[i].Value = num
			ts.Tokens[i].Name = ""
		}
	}
}

// statement will evaluate the input, either by processing the expression
// or by doing a declaration
func (ts *TokenStream) statement(a *Assignment) (int, error) {
	t := ts.get()

	// if RHS is a variable, assign the looked up values
	if _, ok := checkVar(t.Name, a); ok {
		t2 := ts.get()
		if t2 != nil && t2.Name != "=" {
			ts.changeVarsInTokenStream(a)
		}
		ts.putBack()
	}

	switch {
	case t.Name == "":
		ts.changeVarsInTokenStream(a)
		ts.putBack()
		return ts.expression(), nil
	default:
		ts.putBack()
		return ts.declaration(a)
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
	for {
		t := ts.get()
		if t == nil {
			return left
		}
		if t.Kind != '+' && t.Kind != '-' {
			ts.putBack()
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
	for {
		t := ts.get()
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
	if t == nil {
		return 0
	}
	if t.Kind == '(' {
		result := ts.expression()
		t = ts.get()
		if t == nil /* t.Kind != ')'*/ {
			fmt.Println("Invalid expression")
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
	return 0
}
