package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	log, _ := NewLogFile("log.txt")
	defer log.Close()
	var wg sync.WaitGroup

	for i := 0; i < 1000000; i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log.WriteLog(Logger{
				Level:     rand.Intn(5),
				Message:   fmt.Sprintf("这是第%d条日志", i),
				TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
			})
		}(i)
	}
	wg.Wait()
}
