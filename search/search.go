package main

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"log"
	"fmt"
)


type Item struct {
	Title string
	Link string
	Snippet string	
}

type Response struct {
	Items []Item
}

func main() {
	const key = "AIzaSyC7mgdbc3yXCocyHSydGNnUIByw1MIPGFY"
	const cx = "003210436456761614684:ytytx_0dwrc"
	const num = "3"

	base := "https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&num=%s"
	query := "chess"
	
	url := fmt.Sprintf(base, key, cx, query, num)
	fmt.Println("URL is:", url)

	resp, err := http.Get(url)	
	if err != nil {
		log.Fatal("HTTP GET Error: ", err)	
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("IOUtil.ReadAll Error: ", err)
	}


	var v Response	
	err = json.Unmarshal(body, &v)
	if err != nil {
		log.Fatal("json.Unmarshal Error: ", err)
	}

	//fmt.Println(string(body))

	res := v.Items
	for _, val := range res {
		fmt.Println(val.Title)
	}

	//fmt.Println(v.Items)
}
