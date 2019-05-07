package project

import (
	"encoding/json"
)

// MarshalJSON encode project in JSON.
// Project URL is encoded as a string and attachements URL are encoded as strings too.
func (p Project) MarshalJSON() ([]byte, error) {
	var attachement []string
	for _, a := range p.Attachement {
		attachement = append(attachement, a.String())
	}

	basicProject := struct {
		Title        string    `json:"title"`
		URL          string    `json:"url"`
		Description  string    `json:"description"`
		Localisation string    `json:"localisation"`
		Theme        string    `json:"theme"`
		Like         int       `json:"like"`
		Follower     int       `json:"follower"`
		IsPopular    bool      `json:"ispopular"`
		Year         int       `json:"year"`
		Status       string    `json:"status"`
		Cost         int       `json:"cost"`
		Author       string    `json:"author"`
		Creation     Creation  `json:"creation"`
		Attachement  []string  `json:"attachement"`
		Comment      []Comment `json:"comment"`
	}{
		Title:        p.Title,
		URL:          p.URL.String(),
		Description:  p.Description,
		Localisation: p.Localisation,
		Theme:        p.Theme,
		Like:         p.Like,
		Follower:     p.Follower,
		IsPopular:    p.IsPopular,
		Year:         p.Year,
		Status:       p.Status,
		Cost:         p.Cost,
		Author:       p.Author,
		Creation:     p.Creation,
		Attachement:  attachement,
		Comment:      p.Comment,
	}

	return json.Marshal(basicProject)
}
