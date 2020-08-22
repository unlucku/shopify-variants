package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Product struct {
	Product ProductBody
}
type ProductBody struct {
	Variants []Variant
}
type Variant struct {
	ID    uint
	Title string
}

var (
	client = &http.Client{
		Timeout: 30 * time.Second,
	}
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify a product url")
	}

	rawurl := os.Args[1]

	url, err := url.Parse(rawurl)

	if err != nil && !strings.Contains(url.Path, "products") {
		log.Fatalf("URL is not valid - %s: %v\n", rawurl, err)
	}

	fmt.Println("Getting variants for", url.String())

	req, err := http.NewRequest("GET", fmt.Sprintf("%s.json", url.String()), nil)

	if err != nil {
		log.Fatalf("Failed to create request: %v\n", err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Failed to carry out request: %v\n", err)
	}

	defer resp.Body.Close()

	isShopify := false
	for _, cookie := range resp.Cookies() {
		if strings.Contains(cookie.Name, "shopify") {
			isShopify = true
		}
	}

	if !isShopify {
		log.Println("[WARN] Probably not a shopify store")
	}

	switch resp.StatusCode {
	case 200:
		var parsedJSON Product

		err = json.NewDecoder(resp.Body).Decode(&parsedJSON)

		if err != nil {
			log.Fatalf("Could not decode JSON: %v\n", err)
		}

		for _, product := range parsedJSON.Product.Variants {
			fmt.Printf("%v - %s\n", product.ID, product.Title)
		}
	default:
		log.Fatalf("Invalid status code: %v", resp.StatusCode)
	}

}
