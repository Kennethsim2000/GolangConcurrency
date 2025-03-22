package main
import (
	"fmt"
	"sync"
)

func main() {
    var wg sync.WaitGroup
    for _, salutation := range[] string{"hello", "greetings", "good day"} {
        wg.Add(1);
        go print(salutation, &wg);
    }
    wg.Wait()
}

func print(salutation string, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Println(salutation);
}

//testing out using waitgroup to ensure all goroutines are complete