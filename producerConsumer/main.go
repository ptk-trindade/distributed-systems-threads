package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var MAX_CONSUMED int = int(1e5)
var MAX_NUMBER_PRODUCED int = int(1e7)

func consumer(sv sharedVector, wg *sync.WaitGroup) {
	val, keepGoing := sv.pop()
	for keepGoing {
		if isPrime(val) {
			fmt.Println(val, "is prime")
		} else {
			fmt.Println(val, "is not prime")
		}
		val, keepGoing = sv.pop()
	}
	wg.Done()
}

func producer(sv sharedVector) {
	keepGoing := true
	for keepGoing {
		keepGoing = sv.insert(RandInt(MAX_NUMBER_PRODUCED))
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing parameters, insert N (length of shared vector), Np (qty of producers), Nc (qty of consumers)")
		fmt.Println("Usage: go run . <N> <Np> <Nc>")
		os.Exit(1)
	}

	args := os.Args[1:]
	n_str, np_str, nc_str := args[0], args[1], args[2]

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

	if n < 1 || np < 1 || nc < 1 {
		fmt.Println("Error: Params must be greater than 0")
		os.Exit(1)
	}

	var sv sharedVector
	if len(args) > 3 && args[3] == "v2" {
		sv = newSharedVectorV2(n)
	} else {
		sv = newSharedVectorV1(n)
	}

	startTime := time.Now()

	// Start producers
	for i := 0; i < np; i++ {
		go producer(sv)
	}

	// Start consumers
	var wg sync.WaitGroup
	wg.Add(nc)
	for i := 0; i < nc; i++ {
		go consumer(sv, &wg)
	}

	wg.Wait()

	elapsedTime := time.Since(startTime)

	fmt.Println("Elapsed time: ", elapsedTime)

	historic := sv.getHistoric()
	writeFile(fmt.Sprintf("%v", historic), "data.txt")
}
