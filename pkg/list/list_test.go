package list

import (
	"net/url"
	"os"
	"testing"
)

func TestProcess(t *testing.T) {
	file, err := os.Open("../../fixtures/list.html")
	if err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse("http://foo.bar")
	if err != nil {
		t.Fatal(err)
	}

	list := List{
		baseURL: u,
	}

	urls, err := list.process(file)
	if err != nil {
		t.Fatal(err)
	}

	expectedURLs := []string{
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=2",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=4",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=5",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=6",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=7",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=9",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=10",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=11",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=12",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=13",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=14",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=16",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=18",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=20",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=21",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=23",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=24",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=25",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=26",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=27",
		"https://budgetparticipatif.paris.fr/bp/jsp/site/Portal.jsp?page=idee&campagne=C&idee=29",
	}

	for _, e := range expectedURLs {
		found := func(u string, urls []*url.URL) bool {
			for _, t := range urls {
				if u == t.String() {
					return true
				}
			}
			return false
		}(e, urls)
		if !found {
			t.Fatalf("url %#q not found", e)
		}
	}
}
