package main
import (
	"fmt"
	"net/http"
	"sync"
)

type Result struct {
	Error error
	Response *http.Response
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1);
	go readErrors(&wg);
	wg.Wait()
}

// getUrl can take in any  amount of string
func getUrl(urls ...string) chan Result{
	results := make(chan Result)
	go func() { // here, i want to make the call 
		defer close(results)
		for _, url := range urls {
			var result Result
			resp, err := http.Get(url)
			result = Result{Error: err, Response: resp}
			results <- result
		}
	} ()
	return results
}

func readErrors(wg *sync.WaitGroup) {
	urls := []string {"https://www.google.com", "https://badhost"};
	defer wg.Done()
	for result := range getUrl(urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v", result.Error)
			continue
		}
	}
}

/*this is a example to make use of a goroutine to completely handle error checking 

we need to have a channel to read from an errors channel, and process it(in this case we are printing) 
*/