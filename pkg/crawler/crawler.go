package crawler

import (
	"fmt"
	"github.com/wikidistance/wikidist/pkg/db"
	"sync"
)

type Crawler struct {
	nWorkers int
	startURL string

	results  chan db.Article
	database db.DB

	toSee map[string]struct{}
	l     sync.Mutex

	seen  map[string]struct{}
	graph map[string]db.Article
}

func NewCrawler(nWorkers int, startURL string, database db.DB) *Crawler {
	c := Crawler{}

	c.database = database

	c.nWorkers = nWorkers
	c.startURL = startURL

	c.results = make(chan db.Article, nWorkers)
	c.seen = make(map[string]struct{})
	c.toSee = make(map[string]struct{})
	c.graph = make(map[string]db.Article)

	return &c
}

func (c *Crawler) Run() {
	nQueued := 1
	c.toSee[c.startURL] = struct{}{}
	c.seen[c.startURL] = struct{}{}

	for i := 1; i <= c.nWorkers; i++ {
		go c.addWorker()
	}

	for nCrawled := 0; nQueued > nCrawled; nCrawled++ {
		result := <-c.results
		fmt.Println("got result", result.Title, len(result.LinkedArticles))
		resultCopy := result

		c.database.AddVisited(&resultCopy)

		c.graph[result.URL] = result
		for _, neighbour := range result.LinkedArticles {
			fmt.Println(neighbour.URL)
			if _, ok := c.seen[neighbour.URL]; !ok {
				nQueued++

				c.l.Lock()
				c.toSee[neighbour.URL] = struct{}{}
				c.l.Unlock()

				c.seen[neighbour.URL] = struct{}{}
			}
		}

		fmt.Println(nQueued, "queued,", nCrawled, "crawled")
	}
}

func (c *Crawler) addWorker() {
	for {
		var url string
		c.l.Lock()
		for link := range c.toSee {
			url = link
			break
		}
		delete(c.toSee, url)
		c.l.Unlock()

		if url == "" {
			continue
		}

		fmt.Println("getting", url)
		c.results <- CrawlArticle(url)
	}
}