package main

import "fmt"

type B struct {
	Some int
}
type U struct {
	b B
}

func main() {
	u := U{b: B{Some: 1}}

	fmt.Println(u)
}
