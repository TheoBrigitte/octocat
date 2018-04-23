package project

import (
	"fmt"
	"strconv"
	"strings"
)

// CSVHeader returns csv header aligned with Project.CSV() columns.
func CSVHeader() []string {
	return []string{
		"title",
		"url",
		"description",
		"localisation",
		"theme",
		"like",
		"follower",
		"is popular",
		"year",
		"status",
		"cost",
		"author",
		"creation author",
		"creation date",
		"attachement",
		"comment",
	}
}

// CSV return csv columns from the project.
func (p Project) CSV() []string {
	var (
		attachement []string
		comment     []string
	)

	for _, a := range p.Attachement {
		attachement = append(attachement, a.String())
	}

	for _, c := range p.Comment {
		comment = append(comment, fmt.Sprintf("auteur:%s, date:%s\n%s", c.Author, c.Date.String(), c.Text))
	}

	return []string{
		p.Title,
		p.URL.String(),
		p.Description,
		p.Localisation,
		p.Theme,
		strconv.Itoa(p.Like),
		strconv.Itoa(p.Follower),
		strconv.FormatBool(p.IsPopular),
		strconv.Itoa(p.Year),
		p.Status,
		strconv.Itoa(p.Cost),
		p.Author,
		p.Creation.Author,
		p.Creation.Date.String(),
		strings.Join(attachement, "\n"),
		strings.Join(comment, "\n"),
	}
}
