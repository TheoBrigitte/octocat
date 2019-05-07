package project

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/TheoBrigitte/octocat/pkg/httputil"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

type Comment struct {
	Author string    `json:"author"`
	Date   time.Time `json:"date"`
	Text   string    `json:"text"`
}

type Creation struct {
	Author string    `json:"author"`
	Date   time.Time `json:"date"`
}

// Config used to initialize a new Project.
type Config struct {
	Logger     *log.Logger
	HTTPClient *httputil.Client
	URL        *url.URL
}

// Project holds various information related to a project page.
type Project struct {
	logger     *log.Logger
	httpClient *httputil.Client

	Title        string     `json:"title"`
	URL          *url.URL   `json:"url"`
	Description  string     `json:"description"`
	Localisation string     `json:"localisation"`
	Theme        string     `json:"theme"`
	Like         int        `json:"like"`
	Follower     int        `json:"follower"`
	IsPopular    bool       `json:"ispopular"`
	Year         int        `json:"year"`
	Status       string     `json:"status"`
	Cost         int        `json:"cost"`
	Author       string     `json:"author"`
	Creation     Creation   `json:"creation"`
	Attachement  []*url.URL `json:"attachement"`
	Comment      []Comment  `json:"comment"`
}

// New creates a new Project.
func New(config Config) Project {
	return Project{
		logger:     config.Logger,
		httpClient: config.HTTPClient,
		URL:        config.URL,
	}
}

// Fetch the project from its url.
func (p *Project) Fetch() error {
	res, err := p.httpClient.Get(p.URL.String())
	defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("Fetch: %v : %s", err, p.URL.String())
	}

	return p.parse(res.Body)
}

func (p *Project) parse(r io.ReadCloser) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return fmt.Errorf("parse: %v : %s", err, p.URL.String())
	}

	p.parseTitle(doc)
	p.parseDescription(doc)
	p.parseLocalisation(doc)
	p.parseTheme(doc)
	p.parsePopular(doc)
	p.parseAuthor(doc)
	p.parseComment(doc)

	err = p.parseLike(doc)
	if err != nil {
		return fmt.Errorf("parse: %v : %s", err, p.URL.String())
	}

	err = p.parseFollower(doc)
	if err != nil {
		return fmt.Errorf("parse: %v : %s", err, p.URL.String())
	}

	err = p.parseYear(doc)
	if err != nil {
		return fmt.Errorf("parse: %v : %s", err, p.URL.String())
	}

	err = p.parseStatus(doc)
	if err != nil {
		return fmt.Errorf("parse: %v : %s", err, p.URL.String())
	}

	err = p.parseCost(doc)
	if err != nil {
		return fmt.Errorf("parse: %v : %s", err, p.URL.String())
	}

	err = p.parseCreation(doc)
	if err != nil {
		return fmt.Errorf("parse: %v : %s", err, p.URL.String())
	}

	err = p.parseAttachement(doc)
	if err != nil {
		return fmt.Errorf("parse: %v : %s", err, p.URL.String())
	}

	err = p.parseComment(doc)
	if err != nil {
		return fmt.Errorf("parse: %v : %s", err, p.URL.String())
	}

	return nil
}
