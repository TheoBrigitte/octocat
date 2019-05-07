package project

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/TheoBrigitte/octocat/pkg/project"
)

var p project.Project

func init() {
	u, err := url.Parse("")
	if err != nil {
		panic(err)
	}

	p = project.Project{
		Title:        "",
		Description:  "",
		Localisation: "",
		Theme:        "",
		Like:         0,
		Follower:     0,
		IsPopular:    false,
		Year:         0,
		Status:       "",
		Cost:         0,
		Author:       "",
		URL:          u,
		Creation: project.Creation{
			Author: "",
			Date:   time.Date(2019, 12, 31, 12, 0, 0, 0, time.UTC),
		},
		Comment: []project.Comment{
			{
				Author: "",
				Date:   time.Date(2019, 12, 31, 12, 0, 0, 0, time.UTC),
				Text:   "",
			},
		},
	}
}

func TestParse(t *testing.T) {
	file, err := os.Open("../../fixtures/project.html")
	if err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse("http://foo.bar")
	if err != nil {
		t.Fatal(err)
	}

	project := Project{
		URL: u,
	}

	err = project.parse(file)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Test project fields.
}
