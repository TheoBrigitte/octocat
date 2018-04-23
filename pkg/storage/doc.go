// Package storage provide Storage client which is a superset of csv.Writer.
//
// 	s := storage.New(storage.Config{...})
// 	err := s.Write([]string{"some data"})
// 	...
// 	err = Flush()
// 	...
package storage
