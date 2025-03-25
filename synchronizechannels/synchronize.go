package main
import (
	"fmt"
	"sync"
)

func main() {
    var wg sync.WaitGroup
    intstream := make(chan int)
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go write(&wg, intstream, i)
    }
    go closeChan(&wg, intstream)
    for integer := range intstream {
        fmt.Printf("%v ", integer)
    }
}

func closeChan(wg *sync.WaitGroup, intstream chan int) {
    wg.Wait()
    close(intstream)
}

func write(wg *sync.WaitGroup, intstream chan int, num int ) {
    defer wg.Done()
    intstream <- num
}

//Testing out synchronization using channels