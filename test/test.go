package test

import "fmt"

func TestA() {
	for _, val := range values {
		go func() {
			fmt.Println(val)
		}()
	}
}
