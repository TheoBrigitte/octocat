package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/TheoBrigitte/octocat/pkg/httputil"
	"github.com/TheoBrigitte/octocat/pkg/list"
	"github.com/TheoBrigitte/octocat/pkg/project"
	"github.com/TheoBrigitte/octocat/pkg/storage"
	csvStorage "github.com/TheoBrigitte/octocat/pkg/storage/csv"
	jsonStorage "github.com/TheoBrigitte/octocat/pkg/storage/json"

	log "github.com/sirupsen/logrus"
)

var (
	baseURL      = flag.String("baseURL", "https://budgetparticipatif.paris.fr/bp/", "absolute url to the website")
	limit        = flag.Int("limit", -1, "limit of projects to fetch")
	logFile      = flag.String("log", "stdout", "log file")
	itemsPerPage = flag.Int("pagination", 100, "number of item per page")
	outputFile   = flag.String("o", "", "output file")
	outputType   = flag.String("type", "json", "output file type")
	searchPath   = flag.String("searchPath", "jsp/site/Portal.jsp?page=search-solr&conf=list_idees", "relative path to the search page")
)

func pause() {
	fmt.Println("bye")
	fmt.Scanln()
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\nExtract projects data from budgetparticipatif.paris.fr and store them into a file.\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	log.RegisterExitHandler(pause)
	defer pause()

	var (
		storageChan = make(chan interface{})
		urlChan     = make(chan *url.URL, *itemsPerPage)
		numberFired = make(chan int, *itemsPerPage)
		done        = make(chan bool)
		errors      = make(chan error)
	)

	{
		var logOutput io.Writer
		if *logFile != "stdout" {
			file, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("log: os.OpenFile: %v", err)
			}
			defer file.Close()
			logOutput = file
			log.Printf("log file: %s", file.Name())
		} else {
			logOutput = os.Stdout
		}
		log.SetOutput(logOutput)
	}

	var l list.List
	{
		u, err := url.Parse(*baseURL)
		if err != nil {
			log.Fatal(err)
		}

		s, err := url.Parse(*searchPath)
		if err != nil {
			log.Fatal(err)
		}

		// Allow for as many http connection as needed.
		h := httputil.NewWithTransport(&http.Transport{
			MaxIdleConns:        0,
			IdleConnTimeout:     30 * time.Second,
			MaxIdleConnsPerHost: *itemsPerPage,
		})

		c := list.Config{
			Logger:       log.StandardLogger(),
			BaseURL:      u,
			HTTPClient:   h,
			ItemsPerPage: *itemsPerPage,
			Limit:        *limit,
			SearchPath:   s,
		}
		l = list.New(c)
	}

	var s storage.Storage
	{
		if *outputFile == "" {
			path, err := os.Executable()
			if err != nil {
				log.Fatal(err)
			}
			dir := filepath.Dir(path)
			*outputFile = dir + "/octocat_output_" + time.Now().Format("2006-01-02T15-04-05") + "." + *outputType
		}

		_, err := os.Stat(*outputFile)
		if err == nil {
			log.Fatalf("file %v already exist", *outputFile)
		}

		file, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		switch *outputType {
		case "json":
			c := jsonStorage.Config{
				Logger: log.StandardLogger(),
				Writer: file,
			}
			s = jsonStorage.New(c)
		case "csv":
			c := csvStorage.Config{
				Logger: log.StandardLogger(),
				Header: project.CSVHeader(),
				Writer: file,
			}
			s = csvStorage.New(c)
		default:
			log.Fatalf("outputType %q not supported", *outputType)
		}
	}

	log.Println("start octocat")
	log.WithFields(log.Fields{
		"baseURL":      *baseURL,
		"limit":        *limit,
		"logFile":      *logFile,
		"itemsPerPage": *itemsPerPage,
		"outputFile":   *outputFile,
		"outputType":   *outputType,
		"searchPath":   *searchPath,
	}).Print("config")

	s.Start(storageChan, numberFired, done, errors)
	l.ProcessProjects(urlChan, storageChan, numberFired, errors)
	l.CollectProjects(urlChan, errors)

	select {
	case <-done:
		log.Print("done")
	case err := <-errors:
		log.Fatal(err)
	}
}
