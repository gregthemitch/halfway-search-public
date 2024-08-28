package main

import (
	"fmt"
	"math"
)

func main() {
	centroids, hexes := get_hexes(Point{0, 0}, 2, 0, 2)
	fmt.Println("Centroid list", centroids)
	fmt.Println("Hex polygon list", hexes)
}

type Point struct {
	x float64
	y float64
}

type Polygon struct {
	shape [7]Point
}

func contains(running_list *[]Point, point Point) bool {
	for _, p := range *running_list {
		if p == point {
			return true
		}
	}
	return false
}

func draw_points(centroid Point, radius float64, apothem float64, offset float64, hex_num int, max_hex int, points_list *[]Point) {

	// Recursively draw centroids for hexagons stemming from an original point. The stopping condition is the number of hexagon "steps" away from the original point.
	// Update to use spatial intersection

	// Append point only if it doesn't exist in list
	if !contains(points_list, centroid) {
		*points_list = append(*points_list, centroid)
	}
	// Recursive stopping condition
	if hex_num == max_hex {
		return
	}

	directions := []string{"top", "top-left", "top-right", "bottom", "bottom-left", "bottom-right"}

	for _, direction := range directions {
		var new_centroid Point

		if direction == "top" {
			new_centroid = Point{centroid.x + 0, centroid.y + apothem*2}
		}
		if direction == "bottom" {
			new_centroid = Point{centroid.x + 0, centroid.y - apothem*2}
		}
		if direction == "top-left" {
			new_centroid = Point{centroid.x - offset, centroid.y + apothem}
		}
		if direction == "top-right" {
			new_centroid = Point{centroid.x + offset, centroid.y + apothem}
		}
		if direction == "bottom-left" {
			new_centroid = Point{centroid.x - offset, centroid.y - apothem}
		}
		if direction == "bottom-right" {
			new_centroid = Point{centroid.x + offset, centroid.y - apothem}
		}

		draw_points(new_centroid, radius, apothem, offset, hex_num+1, max_hex, points_list)
	}
}

func draw_hex(centroid Point, radius float64, apothem float64) [7]Point {
	// Draw hexagonal polygons around a given point

	top_left := Point{centroid.x - radius*.5, centroid.y + apothem}
	top_right := Point{centroid.x + radius*.5, centroid.y + apothem}
	right := Point{centroid.x + radius, centroid.y}
	bottom_right := Point{centroid.x + radius*.5, centroid.y - apothem}
	bottom_left := Point{centroid.x - radius*.5, centroid.y - apothem}
	left := Point{centroid.x - radius, centroid.y}

	return [7]Point{top_left, top_right, right, bottom_right, bottom_left, left, top_left}
}

func get_hexes(centroid Point, radius float64, hex_num int, max_hex int) ([]Point, []Polygon) {
	// Get centroids and hexagons originating from a point

	apothem := math.Pow((3*radius), .5) / 2
	offset := radius * 1.5

	centroids := make([]Point, 0, 10)
	draw_points(centroid, radius, apothem, offset, hex_num, max_hex, &centroids)

	hex_polygons := make([]Polygon, 0, len(centroids))

	for _, c := range centroids {
		hex_polygons = append(hex_polygons, Polygon{draw_hex(c, radius, apothem)})
	}

	return centroids, hex_polygons

}
