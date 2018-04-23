// Package list provide a way to collect projects.
//
// 	l := list.New(list.Config{...})
// 	for l.HasNext() {
// 		err := l.Fetch()
// 		...
//
// 		// do something with l.Projects
// 	}
package list
