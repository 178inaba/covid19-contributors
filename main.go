package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/PuerkitoBio/goquery"
)

type contributor struct {
	nickname          string
	name              string
	displayName       string
	isGitHubUser      bool
	contributionCount int
}

func main() {
	// Key is name.
	contributorMap := map[string]*contributor{}
	commitsRawurl := "https://github.com/tokyo-metropolitan-gov/covid19/commits/development"
	for {
		d, err := goquery.NewDocument(commitsRawurl)
		if err != nil {
			log.Fatal(err)
		}

		d.Find(".commit-author").Each(func(i int, s *goquery.Selection) {
			c := contributorMap[s.Text()]
			if c == nil {
				c = &contributor{
					nickname: s.Text(),
				}
			}

			if !c.isGitHubUser && goquery.NodeName(s) == "a" {
				c.isGitHubUser = true

				d, err := goquery.NewDocument("https://github.com/" + c.nickname)
				if err != nil {
					log.Fatal(err)
				}

				c.name = d.Find(".p-name").Text()
			}

			c.contributionCount++

			contributorMap[s.Text()] = c
		})

		rawurl, isExists := d.Find(".paginate-container .BtnGroup-item").Eq(1).Attr("href")
		if !isExists {
			break
		}

		commitsRawurl = rawurl
	}

	cs := make([]contributor, len(contributorMap))
	var i int
	for _, c := range contributorMap {
		c.displayName = c.name
		if c.displayName == "" {
			c.displayName = c.nickname
		}

		cs[i] = *c
		i++
	}

	sort.Slice(cs, func(i, j int) bool {
		if cs[i].contributionCount != cs[j].contributionCount {
			return cs[i].contributionCount > cs[j].contributionCount
		}

		return cs[i].displayName < cs[j].displayName
	})

	for _, c := range cs {
		if c.isGitHubUser {
			fmt.Printf("| [%s](https://github.com/%s) ||\n", c.displayName, c.nickname)
		} else {
			fmt.Printf("| %s ||\n", c.displayName)
		}
	}
}
