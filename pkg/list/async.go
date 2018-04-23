package list

import (
	"net/url"

	"github.com/TheoBrigitte/octocat/pkg/project"

	log "github.com/sirupsen/logrus"
)

// CollectProjects collecting project from the list asynchronously in a goroutine.
// Collected projects sent through urlChan channel.
// urlChan channel is closed once all projects have been collected.
// Any error is reported in errors channel.
func (l *List) CollectProjects(urlChan chan<- *url.URL, errors chan<- error) {
	go func() {
		counter := 1

	OuterLoop:
		for l.HasNext() {
			l.logger.WithFields(log.Fields{
				"url": l.searchPath,
			}).Print("list")
			projects, err := l.Fetch()

			if err != nil {
				errors <- err
				return
			}

			for _, p := range projects {
				if l.limit > -1 && counter > l.limit {
					break OuterLoop
				}
				counter++

				urlChan <- p
			}
		}

		close(urlChan)
	}()
}

// ProcessProjects does the actual project fetching asynchronously in a goroutine.
// It collects project url from urlChan channel, and send the resulting project through storageChan channel.
// For every project passed to storageChan channel one entry is sent in numberFired channel.
// numberFired channel is closed once all projects have been fetched.
// Any error is reported in errors channel.
func (l *List) ProcessProjects(urlChan <-chan *url.URL, storageChan chan<- interface{}, numberFired chan<- int, errors chan<- error) {
	go func() {
		var counter int

		for u := range urlChan {
			numberFired <- 1
			counter++

			l.logger.WithFields(log.Fields{
				"count": counter,
				"url":   u.String(),
			}).Print("project")

			go func(u *url.URL) {
				c := project.Config{
					Logger:     l.logger,
					HTTPClient: l.httpClient,
					URL:        u,
				}
				p := project.New(c)

				err := p.Fetch()
				if err != nil {
					errors <- err
					return
				}

				storageChan <- p
			}(u)
		}

		close(numberFired)
	}()
}
