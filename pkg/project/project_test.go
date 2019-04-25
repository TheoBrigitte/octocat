package project

import (
	"net/url"
	"os"
	"testing"
)

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
