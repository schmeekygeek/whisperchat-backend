package main

import "math"

func CalculateDistance(point1 Location, point2 Location) float64 {
  const R = 6371 //Radius of the earth in km
  dLat := deg2rad(float64(point1.lat) - float64(point2.lat))
  dLong := deg2rad(float64(point1.long) - float64(point2.long))
  a := math.Sin(dLat / 2) * math.Sin(dLat / 2) + math.Cos(deg2rad(float64(point1.lat)) * math.Cos(deg2rad(float64(point2.lat)))) * math.Sin(dLong / 2) * math.Sin(dLong / 2)
  c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1 - a))
  d := R * c;
  return d
}

func deg2rad(deg float64) float64 {
  return deg * (math.Pi / 180)
}
