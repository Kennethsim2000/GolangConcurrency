package main
import (
	"fmt"
)

func main() {
	inputsGenerated := generate()
	multiplyRes := multiply(inputsGenerated, 2)
	addRes := add(multiplyRes, 2)
	for num := range addRes {
		fmt.Println(num)
	}
}

func generate() <- chan int {
	generate := make(chan int)
	go func() {
		numsArray := [5]int{1, 2, 3, 4, 5}
		for _ , num := range numsArray {
			generate <- num
		}
		close(generate)
	}()
	return generate
}

func multiply(input <- chan int, constant int) <- chan int{
	output := make(chan int)
	go func() {
		for num := range input {
			newNum := num * constant;
			output <- newNum
		}
		close(output)
	}()
	return output
}

func add(input <- chan int, constant int) <- chan int {
	output := make(chan int)
	go func() {
		for num := range input {
			newNum := num + constant;
			output <- newNum
		}
		close(output)
	}()
	return output
}

/*we want to have a function generate, which returns a channel of input. 
We also want our second stage of our pipeline, which is our multiple stage, this takes in an input channel, creates an output channel, and return*/