package main

import (
	"encoding/json"
	"github.com/awprice/mtbsearch/internal/product"
	"github.com/blevesearch/bleve"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	index, err := bleve.Open("index.bleve")
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile("./data.json")
	if err != nil {
		log.Fatal(err)
	}

	products := make([]*product.ProductInfo, 0)
	if err := json.Unmarshal(data, &products); err != nil {
		log.Fatal(err)
	}

	productMapping := make(map[string]*product.ProductInfo)
	for i := range products {
		productMapping[products[i].ID] = products[i]
	}

	stringQuery := strings.Join([]string{
		"+wheel_size_inches:29",
		"+model_year_float:>=2019",
		"+brakes:shimano",
		"+shifters:shimano",
		"+rear_travel_mm:>=130",
		"+rear_travel_mm:<=150",
		"+fork_travel_mm:>=150",
		"+fork_travel_mm:<=160",
		"+bottle_cage_mounts:yes",
		"+frame_material:carbon",
		"-e_bike_class:class",
	}, " ")

	query := bleve.NewQueryStringQuery(stringQuery)
	search := bleve.NewSearchRequest(query)
	search.Size = len(productMapping)
	searchResults, err := index.Search(search)
	if err != nil {
		log.Fatal(err)
	}

	for _, hit := range searchResults.Hits {
		printProductInfo(productMapping[hit.ID])
	}

	if err := index.Close(); err != nil {
		log.Fatal(err)
	}
}

func printProductInfo(p *product.ProductInfo) {
	log.Println(p.Spec["model_year"], p.Spec["product"], "-", p.Spec["price"])
}
