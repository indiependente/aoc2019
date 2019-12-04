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
	v, err := m.Print(0)
	fmt.Println(v)
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			m.Reset()
			err = m.SetOpCode(1, noun)
			if err != nil {
				log.Fatal(err)
			}
			err = m.SetOpCode(2, verb)
			if err != nil {
				log.Fatal(err)
			}
			err = m.Execute()
			if err != nil && err != EOX {
				log.Fatal(err)
			}
			v, err := m.Print(0)
			if err != nil {
				log.Fatal(err)
			}
			if v == "19690720" {
				fmt.Printf("noun = %d, verb = %d\n", noun, verb)
				fmt.Printf("100 * noun + verb = %d\n", 100*noun+verb)
				break
			}
		}
	}
}
