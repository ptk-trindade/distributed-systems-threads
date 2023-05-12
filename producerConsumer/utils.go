package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

// Random int in range [0, n)
func RandInt(n int) int {
	return seed.Intn(n)
}

func isPrime(n int) bool {
	sqrtn := int(math.Sqrt(float64(n)))
	for i := 2; i < sqrtn; i++ {
		if n%i == 0 {
			return false
		}
	}
	return n > 1
}

func writeFile(txt, fileName string) {
	f, err := os.Create(fileName)

	if err != nil {
		fmt.Println("Error: Couldn't create file: ", err)
	}

	defer f.Close()

	_, err2 := f.WriteString(txt)

	if err2 != nil {
		fmt.Println("Error: Couldn't write string: ", err)
	}

}
