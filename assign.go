package main

import (
	"errors"
	"strconv"
)

type Assignment map[string]int

type Variable struct {
	Name  string
	Value string
}

func (a *Assignment) add(v Variable) error {
	// assign value of another Variable
	//fmt.Printf("add(): name=%q value=%q\n", v.Name, v.Value)
	if _, ok := (*a)[v.Value]; ok {
		(*a)[v.Name] = (*a)[v.Value]
		return nil
	}

	num, err := strconv.Atoi(v.Value)
	if err != nil {
		return errors.New("add(): Invalid assignment")
	}

	// create / update var
	if _, ok := (*a)[v.Name]; !ok {
		//fmt.Printf("add(): update %q %v\n", v.Name, num)
		(*a)[v.Name] = num
	} else {
		//fmt.Printf("add(): update %q %v\n", v.Name, num)
		(*a)[v.Name] = num
	}

	return nil
}

func (a *Assignment) lookup(key string) (int, error) {
	if num, ok := (*a)[key]; ok {
		return num, nil
	}
	return 0, errors.New("Unknown variable")
}
