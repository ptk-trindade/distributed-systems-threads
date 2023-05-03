package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var MAX_CONSUMED int = int(1e5)
var MAX_NUMBER_PRODUCED int = int(1e7)

func consumer(sv *sharedVector) {
	val, keepGoing := sv.pop()
	for keepGoing {
		if isPrime(val) {
			fmt.Println(val, "is prime")
		} else {
			fmt.Println(val, "is not prime")
		}
		val, keepGoing = sv.pop()
	}
}

func producer(sv *sharedVector) {
	keepGoing := true
	for keepGoing {
		keepGoing = sv.insert(RandInt(MAX_NUMBER_PRODUCED))
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing parameters, insert N (length of shared vector), Np (qty of producers) and Nc (qty of consumers)")
		fmt.Println("Usage: go run . <N> <Np> <Nc>")
		os.Exit(1)
	}

	n_str, np_str, nc_str := os.Args[1], os.Args[2], os.Args[3]

	n, err := strconv.Atoi(n_str)
	if err != nil {
		fmt.Println("Error: Invalid N")
		os.Exit(1)
	}

	np, err := strconv.Atoi(np_str)
	if err != nil {
		fmt.Println("Error: Invalid Np")
		os.Exit(1)
	}

	nc, err := strconv.Atoi(nc_str)
	if err != nil {
		fmt.Println("Error: Invalid Nc")
		os.Exit(1)
	}

	sv := newSharedVector(n)

	startTime := time.Now()

	// Start producers
	for i := 0; i < np; i++ {
		go producer(sv)
	}

	// Start consumers
	for i := 1; i < nc; i++ {
		go consumer(sv)
	}

	consumer(sv) // blocking consumer

	elapsedTime := time.Since(startTime)

	fmt.Println("Elapsed time: ", elapsedTime)

	writeFile(fmt.Sprintf("%v", sv.historic), "data.txt")
}
