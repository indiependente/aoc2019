package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	m, err := NewIntcodeMachineWithStdISet(f, ",")
	if err != nil {
		log.Fatal(err)
	}
	err = m.Execute()
	if err != nil && err != EOX {
		log.Fatal(err)
	}
	fmt.Print(m)
}
