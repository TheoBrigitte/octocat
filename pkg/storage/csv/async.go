package csv

import (
	"fmt"
)

// Start collecting data and writting them to the file asynchronously in a goroutine.
// Data is collected from csvChan, data must implement Formatter interface.
// numberFired channel is expected to be close when there is no more data in csvChan channel.
// true is written in done channel on success.
// Any error is reported in errors channel.
func (s CSV) Start(csvChan chan interface{}, numberFired <-chan int, done chan<- bool, errors chan<- error) {
	go func() {
		var records [][]string

		if s.header != nil {
			err := s.Write(s.header)
			if err != nil {
				errors <- err
				return
			}
		}

		for range numberFired {
			i := <-csvChan
			r, ok := i.(Formatter)
			if !ok {
				errors <- fmt.Errorf("%T is not of type Formatter", i)
				return
			}
			records = append(records, r.CSV())
		}

		s.logger.Print("writing file")
		err := s.WriteAll(records)
		if err != nil {
			errors <- err
			return
		}

		err = s.Flush()
		if err != nil {
			errors <- err
			return
		}

		close(csvChan)
		done <- true
	}()
}
