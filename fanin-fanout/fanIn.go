package main
import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Event struct {
	id int64
	procTime time.Duration
}

type Worker struct {
	inputChan <- chan Event
	outputChan chan <- Event
}

type EventFunc func(Event) Event

func main() {
	numChannel := runtime.NumCPU();
	processedChannel := make(chan Event, 5)  
	finalChannel := make(chan Event, 5)  
	inputStream := generateEvents()
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background()) 
	defer cancel()

	for i := 0; i < numChannel; i++ {
		wg.Add(1)
		worker := Worker {
			inputChan: inputStream,
			outputChan: processedChannel,
		}
		worker.start(ctx, &wg, func(e Event) Event {
			time.Sleep(e.procTime)
			return e
		})
	}
	serializer(processedChannel, finalChannel)
	readerDone := make(chan interface{})
	go read(finalChannel, readerDone)
	wg.Wait()
	close(processedChannel)
	<- readerDone
}

func read(finalChannel <- chan Event, readerDone chan interface{}) {
	go func() {
		for output := range finalChannel {
			fmt.Printf("Event id: %d \n", output.id)
		}
		close(readerDone)
	}()
}

func generateEvents() <- chan Event{
	dataChan := make(chan Event);
	go func() {
		for i := 0; i <= 30; i++ {
			event := Event{
				id: int64(i),
				procTime: time.Duration(i%3+1) * time.Second,
			}
			dataChan <- event
		}
		close(dataChan)
	}()
	return dataChan;
}

func serializer(processedChannel <- chan Event, finalChannel chan <- Event) {
	eventMap := make(map[int64]Event)
	var currentEventId int64 = 1
	go func() {
		for event := range processedChannel {
			if(event.id == currentEventId) {
				finalChannel <- event
				currentEventId += 1
				for {
					if event, present := eventMap[currentEventId]; present {
						finalChannel <- event
						currentEventId += 1
					} else {
						break;
					}
				}
			} else {
				eventMap[event.id] = event
			}
		}
		close(finalChannel)
	}()	
}

func (w *Worker) start (ctx context.Context, wg *sync.WaitGroup, fn EventFunc) {
	go func() {
		defer wg.Done()
		for {
			select {
			case e, more := <- w.inputChan:
				if !more {
					return
				}
				select {
				case w.outputChan <- fn(e):
				case <- ctx.Done():
					return
	
				}
			
			case <- ctx.Done():
				return
			}
		}
	}()
}

/*
1.have a function that randomly generates events and pushes into a channel, and returns this channel. 

2.We will also multiplex this work to workers, so we can create a worker type with an input channel and
output channel. We can use a unidirectional type when specifying our worker.

3.We can create a method that belongs to worker (w* worker) before the start method, and this takes in a 
context, waitgroup to wait for worker, and a eventhandler to handle fn.

4.So now, we can multiplex the work to workers, in our main method, we can create a processed Channel of 
Event, as well as the generatedChannel returned from generateEvents function. Then, we can loop through
and create 10 workers, and call their start method, providing a new function.

5. We can have another read method that takes in the processed channel, reads from it.

6. Now, we can implement the serialization of the Event output. We can have a serializer goroutine that 
reads from the processed channel, and outputs to a final channel. In this serializer goroutine, we can
create a map, as well as an currentEventId. For each event read from the processed channel, we check if
the eventId is the currentEventId, if it is, we can send it to the final channel, and have a for loop 
to keep pushing those events in the map as long as it matches the currentEventId. Else, we just store
the event we received in the map.
 */