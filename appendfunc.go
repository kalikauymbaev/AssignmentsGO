package main

import (
	"fmt"
)

func SortSlice(slice []int) {
	n := len(slice)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if slice[j] > slice[j+1] {
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
}

func IncrementOdd(slice []int) {
	for i := 0; i < len(slice); i++ {
		if i%2 == 0 {
			continue
		}
		slice[i] += 1
	}
}

func PrintSlice(slice []int) {
	fmt.Println(slice)
}

func ReverseSlice(slice []int) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func appendFunc(dst func([]int), src ...func([]int)) func([]int) {
	return func(s []int) {
		dst(s)
		for _, fn := range src {
			fn(s)
		}
	}
}

func main() {
	slice := []int{5, 3, 4, 1, 2}

	compositeFunc := appendFunc(SortSlice, IncrementOdd, ReverseSlice, PrintSlice)

	compositeFunc(slice)
}
