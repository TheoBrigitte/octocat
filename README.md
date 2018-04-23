# Octocat

Octocat is a content scrapper for [budgetparticipatif.paris.fr](https://budgetparticipatif.paris.fr/).
It fetch every projects informations found on the [search page](https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=search-solr&conf=list_idees) and store them in json or csv.

It uses:
* [goquery](https://github.com/PuerkitoBio/goquery) for html parsing.
* [logrus](https://github.com/sirupsen/logrus) for logging.

[https://godoc.org/github.com/TheoBrigitte/octocat/pkg](https://godoc.org/github.com/TheoBrigitte/octocat/pkg)

### Scraped data

For every project found it scrap the following data:
* Title
* URL
* Description
* Localisation
* Theme
* Like
* Follower
* IsPopular
* Year
* Status
* Cost
* Author
* Creation date
* Creation author
* Attachements
* Comments

For a visual description of those elements on a project page see [example.png](https://github.com/TheoBrigitte/octocat/project-example.png)

### Download

`go get -v github.com/TheoBrigitte/octocat/cmd/scrape`

### Run

`scrape -o result.json`

### Dependencies

Dependencies are managed using [dep](https://github.com/golang/dep)
