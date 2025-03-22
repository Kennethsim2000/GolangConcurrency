package main
import (
	"fmt"
	"math/rand"
)

func main() {
    done := make (chan interface {})
	randStream := make(chan int)
	go write(done, randStream)
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d \n", i, <-randStream)
	}
	close(done)

}

//takes in a randstream channel to write, as well as a done channel
func write( done <- chan interface{}, randStream chan <- int ) {
	for {
		select {
		case randStream <- rand.Intn(100):
		case <- done:
			return
		}
	}

}


// this example is to use a done channel to stop writer goroutines from writing