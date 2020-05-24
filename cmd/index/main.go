package main

import (
	"encoding/json"
	"github.com/awprice/mtbsearch/internal/product"
	"github.com/blevesearch/bleve"
	"io/ioutil"
	"log"
	"sync"
)

func main() {
	data, err := ioutil.ReadFile("./data.json")
	if err != nil {
		log.Fatal(err)
	}

	info := make([]*product.ProductInfo, 0)
	if err := json.Unmarshal(data, &info); err != nil {
		log.Fatal(err)
	}

	m := bleve.NewIndexMapping()
	index, err := bleve.New("index.bleve", m)
	if err != nil {
		log.Fatal(err)
	}

	numJobs := len(info)
	jobs := make(chan int, numJobs)
	wg := &sync.WaitGroup{}
	for w := 1; w <= 100; w++ {
		wg.Add(1)
		go indexWorker(wg, index, info, jobs)
	}

	for j := 0; j < numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	log.Println("Done")
	log.Println(index.DocCount())

	if err := index.Close(); err != nil {
		log.Fatal(err)
	}
}

func indexWorker(wg *sync.WaitGroup, index bleve.Index, info []*product.ProductInfo, jobs <-chan int) {
	b := index.NewBatch()
	for j := range jobs {
		p := info[j]
		if err := b.Index(p.ID, p.Spec); err != nil {
			log.Fatal(err)
		}
	}
	if err := index.Batch(b); err != nil {
		log.Fatal(err)
	}
	wg.Done()
}