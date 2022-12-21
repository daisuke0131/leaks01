package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

func F1() float64 {
	time.Sleep(5 * time.Second)
	return 80.3
}
func F2() int {
	time.Sleep(2 * time.Second)
	return 45
}

func main() {
	timeout := 3 * time.Second
	var wg sync.WaitGroup
	c := make(chan struct{})
	r1 := make(chan float64)
	r2 := make(chan int)
	wg.Add(2)
	e, result1, result2 := func(wg *sync.WaitGroup, timeout time.Duration) (e error, res1 float64, res2 int) {
		timerT := time.After(timeout)
		log.Print("Start")
		go func() {
			defer close(c)
			wg.Wait()
		}()
		go func(wg *sync.WaitGroup, r chan<- float64) {
			defer wg.Done()
			r <- F1()
			log.Print("Finish F1")
		}(wg, r1)
		go func(wg *sync.WaitGroup, r chan<- int) {
			defer wg.Done()
			r <- F2()
			log.Print("Finish F2")
		}(wg, r2)
		for {
			select {
			case res1 = <-r1:
				{
					log.Print("Received from F1")
				}
			case res2 = <-r2:
				{
					log.Print("Received from F2")
				}
			case <-c:
				return
			case <-timerT:
				{
					e = errors.New("Got Timeout")
					return
				}
			}
		}
	}(&wg, timeout)
	if e != nil {
		log.Println(e)
	}
	log.Printf("Got %f and %d\n", result1, result2)
}
