package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Apiconfig - Importing access token from config.json
type Apiconfig struct {
	Token string
}

func main() {
	// Read in the access token from an untracked local file
	file, err := os.Open("apiconfig.json")
	if err != nil {
		fmt.Println("File failed: ", err)
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var config Apiconfig
	json.Unmarshal(byteValue, &config)

	client := &http.Client{}

	query := `{
		device(id: 212014918250766) {
		  groupId
		  product {
			shortName
		  }
		}
	  }`

	postdata := []byte(`{"query":` + query + `, "variables": {}}`)
	pd, _ := json.Marshal(string(postdata))
	fmt.Println(string(pd))

	// Generate the API query
	url := "https://api.samsara.com/v1/admin/graphql"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pd))
	if err != nil {
		fmt.Printf("Error generating request: %s", err)
	}

	req.Header.Add("X-Access-Token", config.Token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error getting response: %s", err)
	}
	fmt.Println(resp)
}
