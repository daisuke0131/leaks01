package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

func main() {
	timeout := 3 * time.Second

	c := make(chan struct{}, 1)
	r1 := make(chan float64, 1)
	r2 := make(chan int, 1)

	e, result1, result2 := func(timeout time.Duration) (e error, res1 *float64, res2 *int) {
		var wg sync.WaitGroup
		timerT := time.After(timeout)
		log.Print("Start")

		wg.Add(1)
		go func(r chan<- float64) {
			defer wg.Done()
			r <- F1()
			log.Print("Finish F1")
		}(r1)

		wg.Add(1)
		go func(r chan<- int) {
			defer wg.Done()
			r <- F2()
			log.Print("Finish F2")
		}(r2)

		go func() {
			defer close(c)
			wg.Wait()
		}()

		//log.Println("Number of goroutine: ", runtime.NumGoroutine())

		run := true

		for run {
			select {
			case result1 := <-r1:
				{
					res1 = &result1
					log.Print("Received from F1")
				}
			case result2 := <-r2:
				{
					res2 = &result2
					log.Print("Received from F2")
				}
			case <-c:
				run = false
			case <-timerT:
				{
					e = errors.New("Got Timeout")
					run = false
				}
			}
		}
		return
	}(timeout)

	if e != nil {
		log.Println(e)
	}
	log.Printf("Got %v and %d\n", result1, *result2)

	//Debug
	//log.Println("Number of goroutine: ", runtime.NumGoroutine())
	//time.Sleep(6 * time.Second)
	//log.Println("Number of goroutine: ", runtime.NumGoroutine())
}

func F1() float64 {
	time.Sleep(5 * time.Second)
	return 80.3
}
func F2() int {
	time.Sleep(2 * time.Second)
	return 45
}
