// Package csv provide CSV storage which is a superset of csv.Writer.
//
// 	s := csv.New(csv.Config{...})
// 	err := s.Write([]string{"some data"})
// 	...
// 	err = Flush()
// 	...
package csv
