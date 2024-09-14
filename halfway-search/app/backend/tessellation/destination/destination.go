package dest

import (
	"math"

	"github.com/paulmach/orb"
)

// Relied on the math from Turf.js package's destination function: https://turfjs.org/docs/api/destination
// https://github.com/Turfjs/turf/blob/6889237200a088fa5539f92f361fc96135e09ac1/packages/turf-destination/index.ts#L40

var earthRadius float64 = 6371008.8

// func main() {
// 	// r := -180.0
// 	radius := 500.0
// 	// flats := []float64{-180.0, 180.0, 0}
// 	apothem := math.Pow(3, .5) * radius / 2
// 	fmt.Println(apothem)

// 	fmt.Println(FindDestination(orb.Point{41.7841599, -87.5905214}, apothem, -180))
// 	fmt.Println(FindDestination(orb.Point{41.7841599, -87.5905214}, radius, -150))
// 	fmt.Println(FindDestination(orb.Point{41.7841599, -87.5905214}, radius, -90))
// 	fmt.Println(FindDestination(orb.Point{41.7841599, -87.5905214}, radius, -30))
// 	fmt.Println(FindDestination(orb.Point{41.7841599, -87.5905214}, apothem, 0))
// 	fmt.Println(FindDestination(orb.Point{41.7841599, -87.5905214}, radius, 30))
// 	fmt.Println(FindDestination(orb.Point{41.7841599, -87.5905214}, radius, 90))
// 	fmt.Println(FindDestination(orb.Point{41.7841599, -87.5905214}, radius, 150))
// 	fmt.Println(FindDestination(orb.Point{41.7841599, -87.5905214}, apothem, 180))

// 	newpt := FindDestination(orb.Point{41.7841599, -87.5905214}, 2*apothem, -60)
// 	fmt.Println("")
// 	fmt.Println(newpt)
// 	fmt.Println("")
// 	fmt.Println(FindDestination(newpt, apothem, -180))
// 	fmt.Println(FindDestination(newpt, radius, -150))
// 	fmt.Println(FindDestination(newpt, radius, -90))
// 	fmt.Println(FindDestination(newpt, radius, -30))
// 	fmt.Println(FindDestination(newpt, apothem, 0))
// 	fmt.Println(FindDestination(newpt, radius, 30))
// 	fmt.Println(FindDestination(newpt, radius, 90))
// 	fmt.Println(FindDestination(newpt, radius, 150))
// 	fmt.Println(FindDestination(newpt, apothem, 180))

// 	newpt = FindDestination(orb.Point{41.7841599, -87.5905214}, 2*apothem, 60)
// 	fmt.Println("")
// 	fmt.Println(newpt)
// 	fmt.Println("")
// 	fmt.Println(FindDestination(newpt, apothem, -180))
// 	fmt.Println(FindDestination(newpt, radius, -150))
// 	fmt.Println(FindDestination(newpt, radius, -90))
// 	fmt.Println(FindDestination(newpt, radius, -30))
// 	fmt.Println(FindDestination(newpt, apothem, 0))
// 	fmt.Println(FindDestination(newpt, radius, 30))
// 	fmt.Println(FindDestination(newpt, radius, 90))
// 	fmt.Println(FindDestination(newpt, radius, 150))
// 	fmt.Println(FindDestination(newpt, apothem, 180))

// 	newpt = FindDestination(newpt, 2*apothem, 60)
// 	fmt.Println("")
// 	fmt.Println(newpt)
// 	fmt.Println("")
// 	fmt.Println(FindDestination(newpt, apothem, -180))
// 	fmt.Println(FindDestination(newpt, radius, -150))
// 	fmt.Println(FindDestination(newpt, radius, -90))
// 	fmt.Println(FindDestination(newpt, radius, -30))
// 	fmt.Println(FindDestination(newpt, apothem, 0))
// 	fmt.Println(FindDestination(newpt, radius, 30))
// 	fmt.Println(FindDestination(newpt, radius, 90))
// 	fmt.Println(FindDestination(newpt, radius, 150))
// 	fmt.Println(FindDestination(newpt, apothem, 180))

// 	// for r < 180 {
// 	// 	if contains(&flats, r) {
// 	// 		fmt.Println(FindDestination(orb.Point{41.7841599, -87.5905214}, apothem, r))
// 	// 	} else {
// 	// 		fmt.Println(FindDestination(orb.Point{41.7841599, -87.5905214}, radius, r))
// 	// 	}

// 	// 	r += 60
// 	// }

// }

// func contains(slice *[]float64, num float64) bool {
// 	for _, n := range *slice {
// 		if num == n {
// 			return true
// 		}
// 	}
// 	return false
// }

// Find the coordinates of a point given an origin point, distance, and bearing (degrees measured clockwise from north)
func FindDestination(origin_point orb.Point, distance float64, bearing float64) orb.Point {
	lat1 := degree2rad(origin_point.X())
	lng1 := degree2rad(origin_point.Y())
	distanceRad := lengthToRadians(distance)
	bearingRad := degree2rad(bearing)

	lat2 := math.Asin(
		math.Sin(lat1)*math.Cos(distanceRad) +
			math.Cos(lat1)*math.Sin(distanceRad)*math.Cos(bearingRad))

	lng2 := lng1 + math.Atan2(
		math.Sin(bearingRad)*math.Sin(distanceRad)*math.Cos(lat1),
		math.Cos(distanceRad)-math.Sin(distanceRad)*math.Cos(distanceRad))

	return orb.Point{rad2degree(lat2), rad2degree(lng2)}
}

// Covert degrees to radians
func degree2rad(degree float64) float64 {
	return degree * math.Pi / 180
}

// Convert radians to degrees
func rad2degree(radians float64) float64 {
	return radians * 180 / math.Pi
}

// Covert distance (meters) into radians assuming a spherical Earth
func lengthToRadians(distance float64) float64 {
	return distance / earthRadius
}
