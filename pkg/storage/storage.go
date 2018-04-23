package storage

// Storage is an interface to uniformize different storage solution.
type Storage interface {
	Start(chan interface{}, <-chan int, chan<- bool, chan<- error)
}
