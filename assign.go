package main

import (
	"errors"
	"strconv"
)

type Assignment map[string]int

// Variable represents a variable that can hold the value
// of a `statement`
type Variable struct {
	Name  string
	Value string
}

// adds or updates a single `Variable` to the `Assignment` map
func (a *Assignment) add(v Variable) error {
	// assign value of another Variable
	if _, ok := (*a)[v.Value]; ok {
		(*a)[v.Name] = (*a)[v.Value]
		return nil
	}

	num, err := strconv.Atoi(v.Value)
	if err != nil {
		return errors.New("Invalid assignment")
	}

	// create / update var
	if _, ok := (*a)[v.Name]; !ok {
		(*a)[v.Name] = num
	} else {
		(*a)[v.Name] = num
	}

	return nil
}

// looks up a single key in the `Assignment` and returns its value
func (a *Assignment) lookup(key string) (int, error) {
	if num, ok := (*a)[key]; ok {
		return num, nil
	}
	return 0, errors.New("Unknown variable")
}
