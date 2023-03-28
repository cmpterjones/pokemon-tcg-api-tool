package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	defer timer("main", time.Now())

	// Define and then parse flags with defaults
	// The API does all the hard work for us, so we'll take advantage of that and build queries and set the sort per requirements, freeing us to focus on displaying the data
	baseURL := flag.String("base-url", "https://api.pokemontcg.io/v2/cards", "Base URL including protocol without a trailing slash that queries will be made against")
	limit := flag.Int("limit", 10, "Number of results to return from the API.")
	query := flag.String("query", "!rarity:Rare hp:[90 TO *] (types:Grass OR types:Fire)", "Valid query string for searching the PokemonTCG Cards API")
	orderBy := flag.String("order-by", "id", "Valid sort parameter for the PokemonTCG API")
	flag.Parse()

	// Log the raw values we're using to build the request
	log.Printf("Building request for API with the following parameters:\n URL: %q\n Limit: %d\n Unescaped Query: %q\n OrderBy: %q", *baseURL, *limit, *query, *orderBy)
	request, err := buildRequest(*baseURL, *limit, *query, *orderBy)
	if err != nil {
		log.Fatalf("ERROR building HTTP request URL:\n%s", err) // log.Fatalf automatically exits with an non-zero return code after execution
	}
	// Here we'll display the final URL of the request after escaping the query params, nothing secret should show up here so just print it
	log.Printf("Making HTTP request: GET %s", request.URL.String())
	response, err := doRequest(request)
	if err != nil {
		log.Fatalf("ERROR making HTTP request: %s", err)
	}
	// Now process the response and read the data
	result, err := processData(response)
	if err != nil {
		log.Fatalf("ERROR processing data: %s", err)
	}
	// Make it JSON, 2-space indent
	rawBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("ERROR marshalling ResultData to JSON: %s", err)
	}
	log.Printf("Results:\n%s", string(rawBytes))
}

// buildRequest handles escaping query params and is encapsulated like this to take advantage of the timer function to track execution time for each step of the program
func buildRequest(baseURL string, limit int, query string, orderBy string) (*http.Request, error) {
	defer timer("buildRequest", time.Now()) // Inspired by StackOverflow, since params to defered calls are frozen when the defer is pushed onto the stack, we can use this to generically track exec time
	sanitizedURL := fmt.Sprintf("%s?pageSize=%d&q=%s&orderBy=%s", baseURL, limit, url.QueryEscape(query), url.QueryEscape(orderBy))
	return http.NewRequest(http.MethodGet, sanitizedURL, nil)
}

// doRequest simply makes the HTTP request with timing logged
func doRequest(request *http.Request) (*http.Response, error) {
	defer timer("doRequest", time.Now())
	return http.DefaultClient.Do(request)
}

// This type simply represents the schema of the response from the API.
// Unmarshalling the API response will fill in the necessary struct data using the tags, which will also set the json format on the way back to json
type ResultData struct {
	Data []CardData `json:"data"`
}

// This type represents the card data that we care to output
type CardData struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Types  []string `json:"types"`
	HP     string   `json:"hp"`
	Rarity string   `json:"rarity"`
}

// processData reads the response and parses the data into our custom types above
func processData(response *http.Response) (ResultData, error) {
	defer timer("processData", time.Now())
	var rd ResultData
	// read the response, we'll need it for anything we do from this point on
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return rd, err
	}
	// check the response code, if not 200 we have to stop here
	if response.StatusCode != 200 {
		return rd, errors.New(fmt.Sprintf("API returned %d and response: %s", response.StatusCode, string(body)))
	}
	// Unmarshal the response to our custom type, no need for path expressions or anything complicated, go will match the json keys to the struct tags for us
	err = json.Unmarshal(body, &rd)
	return rd, err
}

// timer is just a fun little way to capture and log function execution timing to provide more execution data at runtime
func timer(function string, now time.Time) {
	log.Printf("Function %q took %s", function, time.Since(now))
}
