package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// ------ LOCK ------
func acquire(lock *atomic.Bool) {
	for !lock.CompareAndSwap(false, true) {
	}
}

func release(lock *atomic.Bool) {
	lock.Store(false)
}

// ----- UTILS ------
var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

// Random int in range [0, n)
func RandInt(n int) int {
	return seed.Intn(n)
}

func sumSlice(slice []int8) int {
	sum := 0
	for _, v := range slice {
		sum += int(v)
	}
	return sum
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing parameters, insert K (qty of threads) and N (qty of numbers)")
		fmt.Println("Usage: go run . <K> <N>")
		os.Exit(1)
	}

	k_str, n_str := os.Args[1], os.Args[2]

	// Convert string to int
	k, err := strconv.Atoi(k_str)
	if err != nil {
		fmt.Println("Error: Invalid parameter K")
		os.Exit(1)
	}

	n, err := strconv.Atoi(n_str)
	if err != nil {
		fmt.Println("Error: Invalid parameter N")
		os.Exit(1)
	}

	if k < 1 || n < 1 {
		fmt.Println("Error: Invalid parameter, K and N must be greater than 0")
		os.Exit(1)
	}

	// Generate random integers
	integers := make([]int8, n)
	for i := 0; i < n; i++ {
		integers[i] = int8(RandInt(201) - 100)
	}

	numbersPerThread := n / k

	var wg sync.WaitGroup
	var lock atomic.Bool
	var totalSum int
	wg.Add(k)

	startTime := time.Now()
	for i := 0; i < k; i++ {

		go func(i int) { // run in a new thread
			sum := sumSlice(integers[i*numbersPerThread : (i+1)*numbersPerThread])
			acquire(&lock)
			totalSum += sum
			release(&lock)
			wg.Done()
		}(i)

	}

	sum := sumSlice(integers[k*numbersPerThread : n])
	acquire(&lock)
	totalSum += sum
	release(&lock)

	wg.Wait() // Wait for all threads to finish

	elapsed := time.Since(startTime)
	fmt.Println("Total sum:", totalSum, "\tElapsed time:", elapsed)

	// Check result
	startTime = time.Now()

	sum2 := sumSlice(integers)

	elapsed = time.Since(startTime)
	fmt.Printf("Single thread sum: %d \t (made in %v)\n", sum2, elapsed)

	if totalSum != sum2 {
		fmt.Println("Error: Sum is not correct")
	}

}
