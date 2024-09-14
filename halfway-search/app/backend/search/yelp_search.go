package search

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"slices"
	"strconv"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/useinsider/go-pkg/insrequester"
)

func YelpSearch(points orb.MultiPoint, query_term string, true_centroid orb.Point) []YelpResponse {
	return postYelp(&points, query_term, true_centroid)
}

type YelpResponse struct {
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Coordinates orb.Point `json:"coordinates"`
	Price       string    `json:"price"`
	URL         string    `json:"url"`
	Distance    float64   `json:"distance"`
}

func constructURLs(point orb.Point, query string) url.Values {

	params := url.Values{}
	params.Add("term", query)
	params.Add("latitude", strconv.FormatFloat(point.X(), 'f', -1, 64))
	params.Add("longitude", strconv.FormatFloat(point.Y(), 'f', -1, 64))
	params.Add("sort_by", "best_match")
	params.Add("limit", "5")

	return params
}

func contains(slice []YelpResponse, response YelpResponse) bool {
	for _, r := range slice {
		if r == response {
			return true
		}
	}
	return false
}

func removeDuplicates(responses *[]YelpResponse) []YelpResponse {
	unique := make([]YelpResponse, 0, 100)
	for _, r := range *responses {
		if contains(unique, r) {
			continue
		} else {
			unique = append(unique, r)
		}
	}

	return unique
}

func postYelp(points *orb.MultiPoint, query string, centroid orb.Point) []YelpResponse {
	requester := insrequester.NewRequester().Load()
	ch := make(chan []byte, len(*points))

	base_url := "https://api.yelp.com/v3/businesses/search?"
	for _, p := range *points {
		go func() {
			headers := insrequester.Headers{{"Authorization": "Bearer " + os.Getenv("YELP_API_KEY")}, {"accept": "application/json"}, {"content-type": "application/json"}}
			res, err := requester.Get(insrequester.RequestEntity{Headers: headers, Endpoint: base_url + constructURLs(p, query).Encode()})
			if err != nil {
				return
			}
			bytes, _ := io.ReadAll(res.Body)
			// Closing response body to prevent memory leak
			defer res.Body.Close()
			ch <- bytes
		}()
	}

	var jsonRes map[string]interface{}
	responses := make([]YelpResponse, 0, len(*points)*20)

	for range *points {
		// Decode json response and put into jsonRes
		if err := json.Unmarshal(<-ch, &jsonRes); err != nil {
			panic(err)
		}

		//Parse JSON to business details
		//Resources reviewed: https://rakaar.github.io/posts/2021-04-23-go-json-res-parse/
		results_map, ok := jsonRes["businesses"]
		if !ok {
			continue
		}

		for _, r := range results_map.([]interface{}) {
			var result YelpResponse

			// Name
			result.Name = r.(map[string]interface{})["name"].(string)

			// Address
			address := r.(map[string]interface{})["location"].(map[string]interface{})["display_address"].([]interface{})
			var full_address string
			for i, a := range address {
				full_address += a.(string)
				if i == 0 {
					full_address += ", "
				}
			}
			result.Address = full_address

			// Coordinates
			lat := r.(map[string]interface{})["coordinates"].(map[string]interface{})["latitude"].(float64)
			lng := r.(map[string]interface{})["coordinates"].(map[string]interface{})["longitude"].(float64)
			result.Coordinates = orb.Point{lat, lng}

			// Price
			if val, ok := r.(map[string]interface{})["price"]; ok {
				result.Price = val.(string)
			} else {
				result.Price = "Unknown"
			}

			// URL
			result.URL = r.(map[string]interface{})["url"].(string)

			// Distance
			// Points are saved as lat/lng internally, so flipping to get accurate distance calculation
			result.Distance = geo.Distance(orb.Point{result.Coordinates.Y(), result.Coordinates.X()}, orb.Point{centroid.Y(), centroid.X()})

			responses = append(responses, result)
		}
	}

	unique_responses := removeDuplicates(&responses)
	sortResponses(unique_responses)

	return unique_responses
}

func sortResponses(responses []YelpResponse) []YelpResponse {
	slices.SortFunc(responses, func(a, b YelpResponse) int {
		if a.Distance > b.Distance {
			return 1
		} else if a.Distance < b.Distance {
			return -1
		} else {
			return 0
		}
	})

	return responses
}
