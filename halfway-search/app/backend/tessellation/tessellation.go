package tessellation

import (
	dest "app/halfway-search/app/backend/tessellation/destination"
	"math"
	"slices"
	"strconv"

	"github.com/engelsjk/polygol"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

// func main() {
// 	fmt.Println(draw_hex(orb.Point{41.7841599, -87.5905214}, miles2meters(2), math.Pow(3, .5)*miles2meters(2)/2))
// 	// fmt.Println(dest.FindDestination(orb.Point{41.7841599, -87.5905214}, 2*math.Pow(3, .5)*miles2meters(2)/2, 0))

// 	// fmt.Println(dest.FindDestination(orb.Point{41.83429639109297, -87.59052139999999}, 2*math.Pow(3, .5)*miles2meters(2)/2, 180))
// 	radius := miles2meters(2)
// 	apothem := math.Pow(3, .5) * radius / 2

// 	fmt.Println(dest.FindDestination(orb.Point{41.7841599, -87.5905214}, 2*apothem, 0))
// 	fmt.Println(dest.FindDestination(orb.Point{41.7841599, -87.5905214}, 2*apothem, 60))
// 	fmt.Println(dest.FindDestination(orb.Point{41.7841599, -87.5905214}, 2*apothem, 120))
// 	fmt.Println(dest.FindDestination(orb.Point{41.7841599, -87.5905214}, 2*apothem, 180))
// 	fmt.Println(dest.FindDestination(orb.Point{41.7841599, -87.5905214}, 2*apothem, -120))
// 	fmt.Println(dest.FindDestination(orb.Point{41.7841599, -87.5905214}, 2*apothem, -60))

// 	fmt.Println("")
// 	newpt := dest.FindDestination(orb.Point{41.7841599, -87.5905214}, 2*apothem, 0)
// 	fmt.Println(dest.FindDestination(newpt, 2*apothem, 0))
// 	fmt.Println(dest.FindDestination(newpt, 2*apothem, 60))
// 	fmt.Println(dest.FindDestination(newpt, 2*apothem, 120))
// 	fmt.Println(dest.FindDestination(newpt, 2*apothem, 180))
// 	fmt.Println(dest.FindDestination(newpt, 2*apothem, -120))
// 	fmt.Println(dest.FindDestination(newpt, 2*apothem, -60))

// }

func Tessellation(addresses []orb.Point, radius_distance_miles float64) (orb.MultiPoint, orb.Point) {
	addressPolygon, centroid := giftWrap(addresses)
	return get_hexes(centroid, radius_distance_miles, len(addresses), addressPolygon), centroid
}

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

// Convert from polygol geometries to orb geometries
func p2g(p [][][][]float64) orb.Geometry {

	g := make(orb.MultiPolygon, len(p))

	for i := range p {
		g[i] = make([]orb.Ring, len(p[i]))
		for j := range p[i] {
			g[i][j] = make([]orb.Point, len(p[i][j]))
			for k := range p[i][j] {
				pt := p[i][j][k]
				point := orb.Point{pt[0], pt[1]}
				g[i][j][k] = point
			}
		}
	}
	return g
}

// Convert miles to meters
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

// Checks a single point against a slice of points to determine if it already exists in the slice
// func contains(running_list *orb.MultiPoint, point orb.Point) bool {
// 	for _, p := range *running_list {
// 		if x, y := math.Abs(round(p.X(), 1e-4)-round(point.X(), 1e-4)), math.Abs(round(p.Y(), 1e-4)-round(point.Y(), 1e-4)); x < 1e-4 && y < 1e-4 {
// 			return true
// 		}
// 	}
// 	return false
// }

// Checks a single point against a slice of points to determine if it already exists in the slice.
// Distance of the new point relative to all other points is calculated. Points shouldn't be at least the apothem's distance away.
// With enough hexagons, the code runs into floating point issues where points should be the same, but end up slightly different. Checking distance is a way of handling that.
func contains(running_list *orb.MultiPoint, point orb.Point, apothem float64) bool {
	for _, p := range *running_list {
		// Points are stored in the opposite order
		pt1 := orb.Point{p.Y(), p.X()}
		pt2 := orb.Point{point.Y(), point.X()}

		if dist := geo.Distance(pt1, pt2); dist < apothem {
			return true
		}
	}
	return false
}

// Recursively draw centroids for hexagons stemming from an original point
func draw_points(centroid orb.Point, radius float64, apothem float64, points_list *orb.MultiPoint, hex_list *[]orb.Polygon, address_poly orb.Polygon) {

	// STOP if centroid has already been checked
	if contains(points_list, centroid, apothem) {
		return
	}

	hex := draw_hex(centroid, radius, apothem)

	// STOP if new hexagon does not intersect the address polygon
	intersection, _ := polygol.Intersection(g2p(address_poly), g2p(hex))

	// address_area := planar.Area(address_poly)
	// intersection_area := planar.Area(p2g(intersection))

	if len(intersection) == 0 {
		return
	}
	// fmt.Println(intersection_area / address_area)

	// Append point only if it doesn't exist in list and intersects address polygon
	*points_list = append(*points_list, centroid)
	// Add hexagon to list of hexagons
	*hex_list = append(*hex_list, hex)

	directions := map[string]float64{
		"top":          0,
		"top-right":    60,
		"bottom-right": 120,
		"bottom":       180,
		"bottom-left":  -120,
		"top-left":     -60}

	for _, degree := range directions {
		new_centroid := dest.FindDestination(centroid, 2*apothem, degree)
		draw_points(new_centroid, radius, apothem, points_list, hex_list, address_poly)
	}

}

// Draw hexagonal polygons around a given point
func draw_hex(centroid orb.Point, radius float64, apothem float64) orb.Polygon {

	var points orb.Ring

	directions_keys := []string{"top", "top-right", "right", "bottom-right", "bottom", "bottom-left", "left", "top-left"}
	directions := map[string]float64{
		"top":          0,
		"top-right":    30,
		"right":        90,
		"bottom-right": 150,
		"bottom":       180,
		"bottom-left":  -150,
		"left":         -90,
		"top-left":     -30}

	for _, dir := range directions_keys {
		if dir == "top" || dir == "bottom" {
			points = append(points, dest.FindDestination(centroid, apothem, directions[dir]))
		} else {
			points = append(points, dest.FindDestination(centroid, radius, directions[dir]))
		}
	}

	return orb.Polygon{append(points, dest.FindDestination(centroid, apothem, directions["top"]))}
}

func get_hexes(centroid orb.Point, radius float64, num_addresses int, address_poly orb.Polygon) orb.MultiPoint {
	// Get centroids and hexagons originating from a point

	radius_meters := miles2meters(radius)
	apothem := math.Pow(3, .5) * radius_meters / 2

	centroids := make(orb.MultiPoint, 0, num_addresses)
	hex_list := make([]orb.Polygon, 0, 20)

	draw_points(centroid, radius_meters, apothem, &centroids, &hex_list, address_poly)

	// fmt.Println(hex_list)

	return centroids
}

// Find the average point across a series of points
func FindCentroid(addresses []orb.Point) orb.Point {

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
		new_point1 := dest.FindDestination(addresses[0], 10, 90)
		new_point2 := dest.FindDestination(addresses[1], 10, 90)

		addresses = append(append(addresses, new_point1), new_point2)
	}

	// Sort addresses ascending on Y-coordinate
	slices.SortFunc(addresses, func(a, b orb.Point) int {
		if n := a.Y() - b.Y(); n > 1e-7 {
			return 1
		} else if n < 1e-7 {
			return -1
		} else {
			return 0
		}
	})

	starting_point := addresses[0]
	centroid := FindCentroid(addresses)

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
