package test

import "fmt"

func TestA() {
	for _, val := range values {
		go func() {
			fmt.Println(val)
		}()
	}
}

func TestB() {
	defer func() {
		fmt.Println("recovered:", recover())
	}()

	panic("not good")
}

func TestC() {
}
