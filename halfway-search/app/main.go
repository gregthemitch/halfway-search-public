package main

import (
	"app/halfway-search/app/backend/geocode"
	"app/halfway-search/app/backend/search"
	"app/halfway-search/app/backend/tessellation"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/paulmach/orb"
)

// Load in environmental variables
func load_env_vars() {
	err := godotenv.Load(".env.go")

	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	load_env_vars()

	http.Handle("/", http.FileServer(http.Dir("halfway-search/app/frontend/pages")))
	http.HandleFunc("/submit", readAddresses)

	// http.Handle("/yes", fmt.Fprintf(w, "yes"))
	fmt.Println("Server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Process address and yelp query and return back query results
func readAddresses(w http.ResponseWriter, r *http.Request) {

	// Time server response
	defer timer()()

	bytes, _ := io.ReadAll(r.Body)
	//Closing response body to prevent memory leak
	defer r.Body.Close()

	var address_map map[string]interface{}
	json.Unmarshal(bytes, &address_map)

	var yelp_search_query string
	var do_yelp_search bool = false
	var addresses []string

	for key, val := range address_map {
		if key == "yelp_search" {
			yelp_search_query = val.(string)
		} else if key == "do-yelp-search" {
			do_yelp_search = true
		} else if key != "yelp_search" && val != "" {
			// Skip empty addresses
			addresses = append(addresses, val.(string))
		}
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	coords := geocode.Geocode(&addresses)

	var response map[string]interface{}

	if do_yelp_search {
		// Using 1 mile radius distance
		query_points, centroid := tessellation.Tessellation(coords, 1)
		yelp_results := search.YelpSearch(query_points, yelp_search_query, centroid)
		fmt.Println("Number of queried points: ", len(query_points))

		// Create a simple JSON response
		response = map[string]interface{}{"addresses": coords, "query_points": query_points, "results": yelp_results}
	} else {
		// Create a simple JSON response
		// If Yelp results aren't requested, just map the centroid
		centroid := tessellation.FindCentroid(coords)
		response = map[string]interface{}{"addresses": coords, "query_points": []orb.Point{centroid}}

	}

	fmt.Println(response)
	fmt.Println("Request completed")

	// Encode the response as JSON and send it
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func timer() func() {
	start := time.Now()
	return func() { fmt.Printf("Request took %v\n", time.Since(start)) }
}
