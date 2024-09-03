package main

import (
	"fmt"
	"math"
	"slices"
)

func main() {
	// centroids, hexes := get_hexes(Point{0, 0}, 2, 0, 2)
	// fmt.Println("Centroid list", centroids)
	// fmt.Println("Hex polygon list", hexes)
	fmt.Println(giftWrap([]Point{{6, 8}, {0, 5}, {8, 6}, {1, 1}, {7, 7}, {-3, 2}, {-4, -4}, {13, -4}, {-10, 6}, {3, 0}, {2, 3}, {8, 9}, {6, -2}}))
}

type Point struct {
	x float64
	y float64
}

type Hexagon struct {
	shape [7]Point
}

type Polygon struct {
	shape []Point
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

func get_hexes(centroid Point, radius float64, hex_num int, max_hex int) ([]Point, []Hexagon) {
	// Get centroids and hexagons originating from a point

	apothem := math.Pow((3*radius), .5) / 2
	offset := radius * 1.5

	centroids := make([]Point, 0, 10)
	draw_points(centroid, radius, apothem, offset, hex_num, max_hex, &centroids)

	hex_polygons := make([]Hexagon, 0, len(centroids))

	for _, c := range centroids {
		hex_polygons = append(hex_polygons, Hexagon{draw_hex(c, radius, apothem)})
	}

	return centroids, hex_polygons

}

func averagePoint(addresses []Point) float64 {
	n := float64(len(addresses))

	x := 0.0

	for _, p := range addresses {
		x += p.x
	}

	return x / n
}

func giftWrap(addresses []Point) Polygon {

	slices.SortFunc(addresses, func(a, b Point) int {
		if n := a.y - b.y; n > 1e9 {
			return 1
		} else if n < 0 {
			return -1
		} else {
			return 0
		}

		// if y-coordinates are the same, order by x-coordinate
		// if a.x > b.x {
		// 	return 1
		// } else if a.x < b.x {
		// 	return -1
		// } else {
		// 	return 0
		// }
	})

	starting_point := addresses[0]
	average_x := averagePoint(addresses)

	left_side := make([]Point, 0, len(addresses))
	right_side := make([]Point, 0, len(addresses))

	for _, a := range addresses[1:] {
		if average_x > a.x {
			left_side = append(left_side, a)
		} else if average_x < a.x {
			right_side = append(right_side, a)
		}
	}

	// Find smallest x value on right side
	smallest_x := 0.0
	biggest_y := 0.0
	for i, a := range right_side {
		if i == 0 || a.x < smallest_x {
			smallest_x = a.x
		}
		if i == 0 || a.y > biggest_y {
			biggest_y = a.y
		}
	}

	fmt.Println(addresses)

	fmt.Println(right_side)
	slices.Reverse(right_side)
	fmt.Println(right_side)

	return Polygon{append(append([]Point{starting_point}, append(append(left_side, Point{smallest_x, biggest_y}), right_side...)...), starting_point)}
}
