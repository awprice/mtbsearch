package productpages

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"net/url"
	"sync"
)

type controller struct {
	baseURL      *url.URL
	queue        *queue.Queue
	collector    *colly.Collector
	productPages *sync.Map
	seenURLs     *sync.Map
}


func New(baseURL, startingURL *url.URL) (*controller, error) {
	c := colly.NewCollector()

	q, err := queue.New(
		10,
		&queue.InMemoryQueueStorage{MaxSize: 1000},
	)
	if err != nil {
		return nil, err
	}

	pController := &controller{
		queue:               q,
		collector:           c,
		baseURL:             baseURL,
		productPages:        &sync.Map{},
		seenURLs:            &sync.Map{},
	}
	c.OnHTML("li.product > h6 > a", pController.handleProductLinkHTML)
	c.OnHTML("div.pagination > a", pController.handlePaginationNextPage)

	if err := pController.queue.AddURL(startingURL.String()); err != nil {
		return nil, err
	}

	return pController, nil
}

func (c *controller) ProductPageURLs() []string {
	urls := make([]string, 0)
	c.productPages.Range(func(key, value interface{}) bool {
		urls = append(urls, key.(string))
		return true
	})
	return urls
}

func (c *controller) addURL(url string) error {
	if _, loaded := c.seenURLs.LoadOrStore(url, nil); loaded {
		return nil
	}
	return c.queue.AddURL(url)
}

func (c *controller) handleProductLinkHTML(html *colly.HTMLElement) {
	href := html.Attr("href")
	c.productPages.Store(href, nil)
}

func (c *controller) handlePaginationNextPage(html *colly.HTMLElement) {
	path := html.Attr("href")
	u, err := url.Parse(path)
	if err != nil {
		return
	}

	query := u.Query()
	newQuery := url.Values{}
	newUrl := &url.URL{
		Path: u.Path,
	}
	for key, values := range query {
		added := map[string]interface{}{}
		for _, v := range values {
			if _, ok := added[v]; !ok {
				added[v] = nil
				newQuery.Add(key, v)
			}
		}
	}
	newUrl.RawQuery = newQuery.Encode()

	newBaseURL := c.baseURL.ResolveReference(newUrl)
	if err := c.addURL(newBaseURL.String()); err != nil {
		log.Fatal(err)
	}
}

func (c *controller) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		err := c.queue.Run(c.collector)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()
}