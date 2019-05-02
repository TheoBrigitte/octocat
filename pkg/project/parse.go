package project

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	titleSelector         = "#prop-desc > div.prop-desc-titre > h1 > p"
	descriptionSelector   = "#prop-desc > div.prop-desc-txt"
	localisationSelector  = "#prop-desc > div.prop-desc-localisation > div"
	themeSelector         = "#prop-header > span"
	likeSelector          = "#prop-data-actions > a:nth-child(1) > div"
	followerSelector      = "#prop-data-actions > a:nth-child(2) > div"
	popularSelector       = "#prop-data-pastilles > img"
	yearSelector          = "#prop-data > div:nth-child(4) > div.prop-data-value"
	statusSelector        = "#prop-data > div:nth-child(5) > div.prop-data-value"
	costSelector          = "#prop-data > div.prop-data-kv.last > div.prop-data-value"
	authorSelector        = "#prop-data-avatar > a.color-theme-transport"
	creationSelector      = "#prop-data-avatar"
	attachmentSelector    = "#prop-desc > div.prop-desc-pj" 
	commentSelector       = "#comments > div.row > div > div.commentBlock.row" // suis allée jusqu'ici :D
	commentAuthorSelector = "div > h4 > span"
	commentDateSelector   = "div > h4"
	commentTextSelector   = "div.separator + div.paddingoverride + div"

	noCost              = "Le coût n'a pas été évalué"
	costPattern         = `[\d,]+`
	commentCleanPattern = `\s*\n\s*\n`
	creationDateFormat  = "02 January 2006"
	creationPattern     = `Déposé\sle\s(.+)\spar\s*(.*)`
	statusPattern       = `Statut\s*:(.+)`
	yearPattern         = `[0-9]{4}`
	commentDatePattern  = `le\s([\d/]+)`
	commentDateFormat   = `02/01/2006`
)

func (p *Project) parseTitle(doc *goquery.Document) {
	p.Title = strings.TrimSpace(doc.Find(titleSelector).Text())
}

func (p *Project) parseDescription(doc *goquery.Document) {
	p.Description = strings.TrimSpace(doc.Find(descriptionSelector).Text())
}

func (p *Project) parseLocalisation(doc *goquery.Document) {
	p.Localisation = strings.TrimSpace(doc.Find(localisationSelector).Text())
}

func (p *Project) parseTheme(doc *goquery.Document) {
	p.Theme = strings.TrimSpace(doc.Find(themeSelector).Text())
}

func (p *Project) parseLike(doc *goquery.Document) error {
	text := doc.Find(likeSelector).Text()
	val, err := strconv.Atoi(text)
	if err != nil {
		return fmt.Errorf("parseLike: %v : %s", err, text)
	}

	p.Like = val
	return nil
}

func (p *Project) parseFollower(doc *goquery.Document) error {
	text := doc.Find(followerSelector).Text()
	val, err := strconv.Atoi(text)
	if err != nil {
		return err
	}

	p.Follower = val
	return nil
}

func (p *Project) parsePopular(doc *goquery.Document) {
	nodes := doc.Find(popularSelector)
	p.IsPopular = nodes.Length() > 0
}

func (p *Project) parseYear(doc *goquery.Document) error {
	text := strings.TrimSpace(doc.Find(yearSelector).Text())

	match := regexp.MustCompile(yearPattern).Find([]byte(text))

	val, err := strconv.Atoi(string(match))
	if err != nil {
		return fmt.Errorf("parseYear: %v : %s", err, string(match))
	}

	p.Year = val
	return nil
}

func (p *Project) parseStatus(doc *goquery.Document) error {
	text := doc.Find(statusSelector).Text()

	matches := regexp.MustCompile(statusPattern).FindStringSubmatch(text)
	if len(matches) < 2 {
		return fmt.Errorf("parseStatus: no matches in %q", text)
	}

	p.Status = strings.TrimSpace(matches[1])
	return nil
}

func (p *Project) parseCost(doc *goquery.Document) error {
	text := strings.TrimSpace(doc.Find(costSelector).Text())

	if text == noCost {
		p.Cost = -1
		return nil
	}

	matches := regexp.MustCompile(costPattern).FindStringSubmatch(text)
	if len(matches) < 1 {
		return fmt.Errorf("parseCost: no matches in %q", text)
	}

	str := strings.Replace(matches[0], ",", "", -1)
	val, err := strconv.Atoi(str)
	if err != nil {
		return fmt.Errorf("parseCost: %v : %s", err, str)
	}

	p.Cost = val
	return nil
}

func (p *Project) parseAuthor(doc *goquery.Document) {
	p.Author = strings.TrimSpace(doc.Find(authorSelector).Text())
}

func (p *Project) parseCreation(doc *goquery.Document) error {
	text := strings.TrimSpace(doc.Find(creationSelector).Text())

	matches := regexp.MustCompile(creationPattern).FindStringSubmatch(text)
	if len(matches) < 2 {
		return fmt.Errorf("parseCreation: no matches in %q", text)
	}

	date, err := time.Parse(creationDateFormat, matches[1])
	if err != nil {
		return fmt.Errorf("parseCreation: %v : %s", err, matches[1])
	}
	c := creation{Date: date}

	if len(matches) >= 3 {
		c.Author = matches[2]
	}

	p.Creation = c

	return nil
}

func (p *Project) parseAttachement(doc *goquery.Document) (mainErr error) {

	doc.Find(attachmentSelector).EachWithBreak(func(i int, s *goquery.Selection) bool {
		val, exist := s.Attr("href")
		if !exist {
			mainErr = fmt.Errorf("parseAttachement: no link in %q", s.Text())
			return false
		}

		u, err := url.Parse(val)
		if err != nil {
			mainErr = fmt.Errorf("parseAttachement: %v : %s", err, val)
			return false
		}

		p.Attachement = append(p.Attachement, p.URL.ResolveReference(u))
		return true
	})

	return mainErr
}

func (p *Project) parseComment(doc *goquery.Document) (mainErr error) {
	doc.Find(commentSelector).EachWithBreak(func(i int, s *goquery.Selection) bool {
		author := strings.TrimSpace(s.Find(commentAuthorSelector).Text())

		dateText := strings.TrimSpace(s.Find(commentDateSelector).Text())
		matches := regexp.MustCompile(commentDatePattern).FindStringSubmatch(dateText)
		if len(matches) < 2 {
			mainErr = fmt.Errorf("parseComment: no matches in %q", dateText)
			return false
		}
		date, err := time.Parse(commentDateFormat, matches[1])
		if err != nil {
			mainErr = fmt.Errorf("parseComment: %v : %s", err, matches[1])
			return false
		}

		text := strings.TrimSpace(doc.Find(commentTextSelector).Text())
		text = regexp.MustCompile(commentCleanPattern).ReplaceAllString(text, "\n")

		p.Comment = append(p.Comment, comment{author, date, text})
		return true
	})

	return mainErr
}
