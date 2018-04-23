package json

// Start collecting data and writting them to the file asynchronously in a goroutine.
// Data is collected from jsonChan.
// numberFired channel is expected to be close when there is no more data in jsonChan channel.
// true is written in done channel on success.
// Any error is reported in errors channel.
func (s JSON) Start(jsonChan chan interface{}, numberFired <-chan int, done chan<- bool, errors chan<- error) {
	go func() {
		var records []interface{}

		for range numberFired {
			records = append(records, <-jsonChan)
		}

		s.logger.Print("writing file")
		err := s.Encode(records)
		if err != nil {
			errors <- err
			return
		}

		close(jsonChan)
		done <- true
	}()
}
