package attendance

import (
	"math"
)

type GeoFence struct{}

func NewGeoFence() *GeoFence {
	return &GeoFence{}
}

func (g *GeoFence) IsWithinRadius(lat1, lon1, lat2, lon2, radiusMeters float64) bool {
	distance := g.HaversineDistance(lat1, lon1, lat2, lon2)
	return distance <= radiusMeters
}

func (g *GeoFence) HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (g *GeoFence) ValidatePunch(lat, lon *float64, locations []AttendanceLocation, requireGPS bool) (bool, string) {
	if !requireGPS {
		return true, ""
	}
	if lat == nil || lon == nil {
		return false, "GPS required"
	}
	for _, loc := range locations {
		if !loc.IsActive {
			continue
		}
		if g.IsWithinRadius(*lat, *lon, loc.Latitude, loc.Longitude, float64(loc.RadiusMeters)) {
			return true, ""
		}
	}
	return false, "outside authorized location"
}
