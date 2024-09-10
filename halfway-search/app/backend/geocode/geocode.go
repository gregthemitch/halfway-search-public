package geocode

import (
	"encoding/json"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/paulmach/orb"
	"github.com/useinsider/go-pkg/insrequester"
)

// Move to main.go file when ready
func getEnvVars(key string) string {
	err := godotenv.Load("../../../../.env.go")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func Geocode(addresses *[]string) []orb.Point {
	return postGoogle(addresses)
}

func constructURLs(address string) url.Values {
	params := url.Values{}
	params.Add("address", address)
	params.Add("key", getEnvVars("GOOGLE_API_KEY"))

	return params
}

func postGoogle(addresses *[]string) []orb.Point {
	requester := insrequester.NewRequester().Load()
	ch := make(chan []byte, len(*addresses))

	base_url := "https://maps.googleapis.com/maps/api/geocode/json?"

	for _, a := range *addresses {
		go func() {
			res, _ := requester.Get(insrequester.RequestEntity{Endpoint: base_url + constructURLs(a).Encode()})
			// if err == nil {
			// 	return
			// }
			bytes, _ := io.ReadAll(res.Body)
			// Closing response body to prevent memory leak
			defer res.Body.Close()
			ch <- bytes
		}()
	}

	var jsonRes map[string]interface{}
	responses := make([]orb.Point, 0, len(*addresses))

	for range *addresses {
		// Decode json response and put into jsonRes
		if err := json.Unmarshal(<-ch, &jsonRes); err != nil {
			panic(err)
		}

		if jsonRes["status"] == "OK" { //Only continue if there are results for the query
			// Parse JSON to find lat and lng coords
			// Resources reviewed: https://rakaar.github.io/posts/2021-04-23-go-json-res-parse/
			results_map := jsonRes["results"].([]interface{})[0].(map[string]interface{})["geometry"].(map[string]interface{})["location"].(map[string]interface{})

			lat, lng := results_map["lat"].(float64), results_map["lng"].(float64)
			// Add coords to responses
			responses = append(responses, orb.Point{lat, lng})
		}

	}

	return responses
}
