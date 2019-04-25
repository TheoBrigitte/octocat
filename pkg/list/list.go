package list

import (
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/TheoBrigitte/octocat/pkg/httputil"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

const (
	projectSelector = "#prop-list > div > div.prop-panel > div.prop-card > a"
	nextSelector    = "#paginator > span.paginator-current-page + a"

	itemPerPageKey = "items_per_page"
)

// Config used to initialize a new List.
type Config struct {
	Logger *log.Logger

	BaseURL      *url.URL
	HTTPClient   *httputil.Client
	ItemsPerPage int
	Limit        int
	SearchPath   *url.URL
}

// List is a collection of projects.
type List struct {
	logger *log.Logger

	baseURL      *url.URL
	httpClient   *httputil.Client
	itemsPerPage int
	limit        int
	searchPath   *url.URL
}

// New creates a new List.
func New(c Config) List {
	return List{
		logger: c.Logger,

		baseURL:      c.BaseURL,
		httpClient:   c.HTTPClient,
		itemsPerPage: c.ItemsPerPage,
		limit:        c.Limit,
		searchPath:   c.BaseURL.ResolveReference(c.SearchPath),
	}
}

// Fetch retrieve all project on the current page and look for the next page if any.
func (l *List) Fetch() ([]*url.URL, error) {
	q := l.searchPath.Query()
	q.Set(itemPerPageKey, strconv.Itoa(l.itemsPerPage))
	l.searchPath.RawQuery = q.Encode()

	res, err := l.httpClient.Get(l.searchPath.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return l.process(res.Body)
}

func (l *List) process(r io.ReadCloser) ([]*url.URL, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	projects, err := l.getProjects(doc)
	if err != nil {
		return nil, err
	}

	next, err := l.getNext(doc)
	if err != nil {
		return nil, err
	}
	l.searchPath = next

	return projects, nil
}

func (l List) getProjects(doc *goquery.Document) (projects []*url.URL, mainErr error) {
	doc.Find(projectSelector).EachWithBreak(func(i int, s *goquery.Selection) bool {
		val, exist := s.Attr("href")
		if !exist {
			mainErr = fmt.Errorf("getProjects: no link: %v", s)
			return false
		}

		u, err := url.Parse(val)
		if err != nil {
			mainErr = err
			return false
		}

		projects = append(projects, l.baseURL.ResolveReference(u))

		return true
	})

	return projects, nil
}

func (l List) getNext(doc *goquery.Document) (*url.URL, error) {
	nodes := doc.Find(nextSelector)

	if nodes.Length() < 1 {
		return nil, nil
	}

	val, exist := nodes.Attr("href")
	if !exist {
		return nil, fmt.Errorf("getNext: no link: %v", nodes)
	}

	u, err := url.Parse(val)
	if err != nil {
		return nil, err
	}

	return l.baseURL.ResolveReference(u), nil
}

// HasNext determine is there is a next page to fetch.
func (l List) HasNext() bool {
	return l.searchPath != nil
}
