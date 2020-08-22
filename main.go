
package main

import (
	"fmt"
	"os"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Product struct {
	Product ProductBody
}
type ProductBody struct {
	Variants []Variant
}
type Variant struct {
	ID uint
	Title string
}
func main() {
	url := os.Args[1];
	if (!strings.Contains(url, "products")) {
		fmt.Println("URL is not a shopify product url.")
		return;
	}
	fmt.Println("Getting variants for", url);
	client := &http.Client{};
	req, err := http.NewRequest("GET", url + ".json", nil)
	resp, err := client.Do(req);
	if (err != nil) {
		fmt.Println(err);
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var parsedJson Product
	json.Unmarshal([]byte(body), &(parsedJson))
	for i := 0; i < len(parsedJson.Product.Variants); i++ {
		fmt.Println(parsedJson.Product.Variants[i].ID, "-" , parsedJson.Product.Variants[i].Title);
	}
}