package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var fuelTotal float64
	s := bufio.NewScanner(f)
	for s.Scan() {
		mass, err := strconv.Atoi(s.Text())
		if err != nil {
			log.Fatal(err)
		}
		fuelTotal += calculateFuel(float64(mass))
	}

	fmt.Printf("Total fuel required: %.2f\n", fuelTotal)
}

func calculateFuel(x float64) float64 {
	f := math.Floor(x/3) - 2
	if f < 0 {
		return 0
	}
	return f + calculateFuel(f)
}
