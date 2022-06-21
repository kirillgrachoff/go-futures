package main

import "github.com/kirillgrachoff/go-futures/future"
import "fmt"

func main() {
	value := future.Async(func() (int, error) {
		return 10, nil
	}).Map(func(s int) (int, error) {
		return s + 5, nil
	}).GetUnsafe()
	fmt.Println(value)
}
