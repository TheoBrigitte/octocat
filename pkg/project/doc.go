// Package project provide a way to retrive project informations.
//
// 	p := project.New(project.Config{...})
// 	err := p.Parse()
// 	...
//
// 	// store in csv file
// 	csv <- project.CSVHeader()
// 	csv <- p.CSV()
package project
