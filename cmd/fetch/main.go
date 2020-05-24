package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/awprice/mtbsearch/internal/product"
	"github.com/awprice/mtbsearch/internal/productpages"
)

const (
	baseURL      = "https://www.vitalmtb.com"
	startingPath = "/product/category/Bikes,2"
)

func main() {
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	startingPathURL, err := url.Parse(startingPath)
	if err != nil {
		log.Fatal(err)
	}

	startingURL := baseURL.ResolveReference(startingPathURL)
	pagesController, err := productpages.New(baseURL, startingURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Fetching product URLs")
	wg := &sync.WaitGroup{}
	pagesController.Run(wg)
	wg.Wait()

	productPageURLs := pagesController.ProductPageURLs()
	log.Println("Retrieved", len(productPageURLs), "product URLs")

	productsController, err := product.New(productPageURLs)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Fetching product specs")
	wg = &sync.WaitGroup{}

	go func() {
		for {
			if productsController.Done() {
				return
			}
			log.Println("Products retrieved:", len(productsController.Products()), "queue length:", productsController.QueueLength())
			time.Sleep(5 * time.Second)
		}
	}()

	productsController.Run(wg)
	wg.Wait()

	data, err := json.Marshal(productsController.Products())
	if err != nil {
		log.Fatal(err)
	}

	if err := writeResults(data); err != nil {
		log.Fatal(err)
	}
}

func writeResults(data []byte) error {
	f, err := os.Create("./data.json")
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}
