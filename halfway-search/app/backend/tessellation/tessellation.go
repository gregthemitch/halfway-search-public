package tessellation

import (
	"math"
	"slices"
	"strconv"

	"github.com/engelsjk/polygol"
	"github.com/paulmach/orb"
)

// [41.7841599 -87.5905214]

func Tessellation(addresses []orb.Point) (orb.MultiPoint, orb.Point) {
	addressPolygon, centroid := giftWrap(addresses)
	return get_hexes(centroid, 1, len(addresses), addressPolygon), centroid
}

// func main() {
// 	// centroids, hexes := get_hexes(Point{0, 0}, 2, 0, 2)
// 	// fmt.Println("Centroid list", centroids)
// 	// fmt.Println("Hex polygon list", hexes)
// 	// fmt.Println(giftWrap([]Point{{6, 8}, {0, 5}, {8, 6}, {1, 1}, {7, 7}, {-3, 2}, {-4, -4}, {13, -4}, {-10, 6}, {3, 0}, {2, 3}, {8, 9}, {6, -2}}))
// 	// addresses := []orb.Point{{6, 8}, {0, 5}}
// 	addresses := []orb.Point{{41.7841599, -87.5905214}, {41.7965962, -87.582055}, {41.800011, -87.595481}}
// 	addressPolygon, centroid := giftWrap(addresses)
// 	fmt.Println(get_hexes(centroid, 100, len(addresses), addressPolygon))

// 	// fmt.Println(dist2CoordOffsets(orb.Point{-87.5905214, 41.7841599}, 2*miles2meters(1)*math.Cos(degree2rad(180/6)), 90))
// 	// fmt.Println(draw_hex(orb.Point{-87.5905214, 41.7841599}, miles2meters(1)))
// 	// A := [][][][]float64{{{{-4, -4}, {1, 1}, {-3, 2}, {2, 3}, {0, 5}, {-10, 6}, {3, 9}, {8, 9}, {6, 8}, {7, 7}, {8, 6}, {3, 0}, {6, -2}, {13, -4}, {-4, -4}}}}
// 	// C := [][][][]float64{{{{2.838907469342252, 2.9504110531823837}, {2.8534002229654405, 2.9504110531823837}, {2.8606465997770347, 2.8461538461538463}, {2.8534002229654405, 2.741896639125309}, {2.838907469342252, 2.741896639125309}, {2.831661092530658, 2.8461538461538463}, {2.838907469342252, 2.9504110531823837}}}}

// 	// union, _ := polygol.Union(A, C)
// 	// intersection, _ := polygol.Intersection(A, C)

// 	// fmt.Println("Union", union)
// 	// fmt.Println("Intersection", intersection)
// }

// Convert from orb geometries to polygol geometries
func g2p(g orb.Geometry) [][][][]float64 {

	var p [][][][]float64

	switch v := g.(type) {
	case orb.Polygon:
		p = make([][][][]float64, 1)
		p[0] = make([][][]float64, len(v))
		for i := range v { // rings
			p[0][i] = make([][]float64, len(v[i]))
			for j := range v[i] { // points
				pt := v[i][j]
				p[0][i][j] = []float64{pt.X(), pt.Y()}
			}
		}
	case orb.MultiPolygon:
		p = make([][][][]float64, len(v))
		for i := range v { // polygons
			p[i] = make([][][]float64, len(v[i]))
			for j := range v[i] { // rings
				p[i][j] = make([][]float64, len(v[i][j]))
				for k := range v[i][j] { // points
					pt := v[i][j][k]
					p[i][j][k] = []float64{pt.X(), pt.Y()}
				}
			}
		}
	}

	return p
}

// func p2g(p [][][][]float64) orb.Geometry {

// 	g := make(orb.MultiPolygon, len(p))

// 	for i := range p {
// 		g[i] = make([]orb.Ring, len(p[i]))
// 		for j := range p[i] {
// 			g[i][j] = make([]orb.Point, len(p[i][j]))
// 			for k := range p[i][j] {
// 				pt := p[i][j][k]
// 				point := orb.Point{pt[0], pt[1]}
// 				g[i][j][k] = point
// 			}
// 		}
// 	}
// 	return g
// }

// Function to convert miles to degree changes
// func miles2degrees(miles float64) float64 {

// 	// One degree change is approximately 69 miles (on Earth)
// 	// See reference https://calculator.academy/miles-to-degrees-calculator/
// 	// lat_conversion_factor := 69
// 	// return miles / float64(conversion_factor)
// }

func miles2meters(miles float64) float64 {
	return miles * 1609.344
}

// Round floats to specified units and truncate any trailing floats
func round(num, unit float64) float64 {
	rounded_num := math.Round(num/unit) * unit

	str_float := strconv.FormatFloat(rounded_num, 'f', -1, 64)
	unit_str := strconv.FormatFloat(unit, 'f', -1, 64)

	var str string
	// Collapse strings from slices up until the digit place specified when rounding
	for i := 0; i < len(str_float) && i < len(unit_str); i++ {
		str += string(str_float[i])
	}

	flt, _ := strconv.ParseFloat(str, 64)

	return flt
}

func degree2rad(degree int) float64 {
	return float64(degree) * math.Pi / 180
}

// Method for finding lat/long coordinates from distance and angle from a point
// https://stackoverflow.com/questions/2187657/calculate-second-point-knowing-the-starting-point-and-distance
func dist2CoordOffsets(point orb.Point, dist float64, theta int) (float64, float64) {
	radians := degree2rad(theta)

	dx := dist * math.Cos(radians) // theta measured clockwise from due east
	dy := dist * math.Sin(radians)

	delta_longitude := dx / (111320 * math.Cos(point.X()))
	delta_latitude := dy / 110540

	return round(delta_longitude, 1e-5), round(delta_latitude, 1e-5)
}

func contains(running_list *orb.MultiPoint, point orb.Point) bool {
	for _, p := range *running_list {
		if x, y := math.Abs(p.X()-point.X()), math.Abs(p.Y()-point.Y()); x < 1e-2 && y < 1e-2 {
			return true
		}
	}
	return false
}

func draw_points(centroid orb.Point, radius float64, apothem float64, offset float64, points_list *orb.MultiPoint, hex_list *[]orb.Polygon, address_poly orb.Polygon) {

	// Recursively draw centroids for hexagons stemming from an original point. The stopping condition is the number of hexagon "steps" away from the original point.
	// Update to use spatial intersection

	// STOP if centroid has already been checked
	if contains(points_list, centroid) {
		return
	}

	hex := draw_hex(centroid, radius)

	// STOP if new hexagon does not intersect the address polygon
	intersection, _ := polygol.Intersection(g2p(address_poly), g2p(hex))
	// fmt.Println("Address poly:", g2p(address_poly))
	// fmt.Println("Intersection:", intersection)
	if len(intersection) == 0 {
		return
	}

	// Append point only if it doesn't exist in list and intersects address polygon
	*points_list = append(*points_list, centroid)
	// Add hexagon to list of hexagons
	*hex_list = append(*hex_list, hex)

	directions := []string{"top", "top-left", "top-right", "bottom", "bottom-left", "bottom-right"}

	for _, direction := range directions {
		var x_offset, y_offset float64

		if direction == "top" {
			x_offset, y_offset = dist2CoordOffsets(centroid, apothem*2, 90)
		}
		if direction == "bottom" {
			x_offset, y_offset = dist2CoordOffsets(centroid, apothem*2, 270)
		}
		if direction == "top-right" {
			x_offset, y_offset = dist2CoordOffsets(centroid, apothem*2, 30)
		}
		if direction == "top-left" {
			x_offset, y_offset = dist2CoordOffsets(centroid, apothem*2, 150)
		}
		if direction == "bottom-left" {
			x_offset, y_offset = dist2CoordOffsets(centroid, apothem*2, 210)
		}
		if direction == "bottom-right" {
			x_offset, y_offset = dist2CoordOffsets(centroid, apothem*2, 330)
		}

		new_centroid := orb.Point{round(centroid.X()+x_offset, 1e-7), round(centroid.Y()+y_offset, 1e-7)}

		// fmt.Println(new_centroid)
		// fmt.Println(hex)
		// if contains(points_list, new_centroid) {
		// 	continue
		// }

		// hex := draw_hex(new_centroid, radius)

		// // STOP if new hexagon does not intersect the address polygon
		// intersection, _ := polygol.Intersection(g2p(address_poly), g2p(hex))
		// // fmt.Println("Address poly:", g2p(address_poly))
		// // fmt.Println("Intersection:", intersection)
		// if len(intersection) == 0 {
		// 	continue
		// }

		draw_points(new_centroid, radius, apothem, offset, points_list, hex_list, address_poly)
	}

}

// Draw hexagonal polygons around a given point
func draw_hex(centroid orb.Point, radius float64) orb.Polygon {

	var points orb.Ring

	degrees := 0
	for degrees <= 360 {
		x_offset, y_offset := dist2CoordOffsets(centroid, radius, degrees)

		points = append(points, orb.Point{round(centroid.X()+x_offset, 1e-7), round(centroid.Y()+y_offset, 1e-7)})
		degrees += 60
	}

	return orb.Polygon{points}
}

func get_hexes(centroid orb.Point, radius float64, num_addresses int, address_poly orb.Polygon) orb.MultiPoint {
	// Get centroids and hexagons originating from a point

	radius_meters := miles2meters(radius)
	apothem := radius_meters * math.Cos(degree2rad(180/6))
	offset := radius_meters * 1.5

	centroids := make(orb.MultiPoint, 0, num_addresses)
	hex_list := make([]orb.Polygon, 0, 20)

	draw_points(centroid, radius_meters, apothem, offset, &centroids, &hex_list, address_poly)

	// fmt.Println(hex_list)

	return centroids
	// 41.783578, -87.591028
}

// Find the average point across a series of points
func findCentroid(addresses []orb.Point) orb.Point {

	n := float64(len(addresses))

	x := 0.0
	y := 0.0

	for _, p := range addresses {
		x += p.X()
		y += p.Y()
	}
	return orb.Point{round(x/n, 1e-7), round(y/n, 1e-7)}
}

// Simple algorithm to create a convex hull. Should likely update to something more formal
func giftWrap(addresses []orb.Point) (orb.Polygon, orb.Point) {

	// Turn lines (2 points) into polygons by giving them more points that are offset by a small modifier (250 meters)
	if len(addresses) == 2 {
		x_offset, y_offset := dist2CoordOffsets(findCentroid(addresses), 250, 180)
		new_point1 := orb.Point{addresses[0].X() + x_offset, addresses[0].Y() + y_offset}
		new_point2 := orb.Point{addresses[1].X() + x_offset, addresses[1].Y() + y_offset}

		addresses = append(append(addresses, new_point1), new_point2)
	}

	// Sort addresses ascending on Y-coordinate
	slices.SortFunc(addresses, func(a, b orb.Point) int {
		if n := math.Abs(a.Y() - b.Y()); n > 1e-7 {
			return 1
		} else if n < 0 {
			return -1
		} else {
			return 0
		}
	})

	starting_point := addresses[0]
	centroid := findCentroid(addresses)

	left_side := make([]orb.Point, 0, len(addresses))
	right_side := make([]orb.Point, 0, len(addresses))

	for _, a := range addresses[1:] {
		if centroid.X() > a.X() {
			left_side = append(left_side, a)
		} else if centroid.X() < a.X() {
			right_side = append(right_side, a)
		}
	}

	slices.Reverse(right_side)

	return orb.Polygon{append(append(orb.Ring{starting_point}, append(left_side, right_side...)...), starting_point)}, centroid
}
