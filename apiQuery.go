package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Query - Used to send query information to graphQL
// Variables is declared as empty rather than as a string of open/close brackets
// GraphQL expects brackets and not a string, otherwise will throw Bad Json request
type Query struct {
	Query     string   `json:"query"`
	Variables struct{} `json:"variables"`
}

// Apiconfig - Importing configs from config.json
type Apiconfig struct {
	Token   string // Access token
	Timeout int    // HTTP timeout
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
	client.Timeout = time.Second * time.Duration(config.Timeout)

	query := `{
		group(id: 3991) {
			name
		}
	}`

	q := Query{
		Query: query,
	}
	b, err := json.Marshal(q)

	// Generate the API query
	url := "https://api.samsara.com/v1/admin/graphql"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		fmt.Printf("Error generating request: %s", err)
	}
	req.Header.Add("X-Access-Token", config.Token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error getting response: %s", err)
	}
	// Check if we get any page errors, this is not caught by err
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
