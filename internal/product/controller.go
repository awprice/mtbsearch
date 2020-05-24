package product

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"strconv"
	"strings"
	"sync"
)

type controller struct {
	queue     *queue.Queue
	collector *colly.Collector
	products  *sync.Map
}

type ProductInfo struct {
	ID   string
	URL  string
	Spec map[string]interface{}
}

func New(urls []string) (*controller, error) {
	c := colly.NewCollector()

	q, err := queue.New(
		10,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)
	if err != nil {
		return nil, err
	}

	pController := &controller{
		queue:     q,
		collector: c,
		products:  &sync.Map{},
	}
	c.OnHTML("div.b-product-specs#specs", pController.handleSpecsHTML)

	for _, u := range urls {
		if err := q.AddURL(u); err != nil {
			return nil, err
		}
	}

	return pController, nil
}

func (c *controller) handleSpecsHTML(html *colly.HTMLElement) {
	productId := html.Attr("data-id")
	productSpecMap := map[string]interface{}{}
	headers := make([]string, 0)
	values := make([]string, 0)
	html.ForEach("table.specs > tbody > tr > th", func(i int, element *colly.HTMLElement) {
		headers = append(headers, strings.TrimSpace(element.Text))
	})
	html.ForEach("table.specs > tbody > tr > td", func(i int, element *colly.HTMLElement) {
		values = append(values, strings.TrimSpace(element.Text))
	})
	for i := range headers {
		cleanedHeader := cleanHeader(headers[i])
		if cleanedHeader == "sizes_and_geometry" {
			continue
		}
		productSpecMap[cleanedHeader] = values[i]

		if v, ok := extractMillimetres(values[i]); ok {
			productSpecMap[cleanedHeader + "_mm"] = v
		}
		if v, ok := extractInches(values[i]); ok {
			productSpecMap[cleanedHeader + "_inches"] = v
		}

		if v, err := strconv.ParseFloat(values[i], 64); err == nil {
			productSpecMap[cleanedHeader + "_float"] = v
		}

	}
	info := &ProductInfo{
		ID:   productId,
		URL:  html.Request.URL.String(),
		Spec: productSpecMap,
	}
	c.products.Store(productId, info)
}

func (c *controller) Products() []*ProductInfo {
	products := make([]*ProductInfo, 0)
	c.products.Range(func(key, value interface{}) bool {
		products = append(products, value.(*ProductInfo))
		return true
	})
	return products
}

func (c *controller) Done() bool {
	return c.queue.IsEmpty()
}

func (c *controller) QueueLength() int {
	s, _ := c.queue.Size()
	return s
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
