package main

import (
	"math"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandSeq(n int) string {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  b := make([]rune, n)
  for i := range b {
    b[i] = letters[r.Intn(len(letters))]
  }
  return string(b)
}

func CalculateDistance(point1 Location, point2 Location) float64 {
  const R = 6371 //Radius of the earth in km
  dLat := deg2rad(float64(point1.Lat) - float64(point2.Lat))
  dLong := deg2rad(float64(point1.Long) - float64(point2.Long))
  a := math.Sin(dLat / 2) * math.Sin(dLat / 2) + math.Cos(
    deg2rad(float64(point1.Lat)) * math.Cos(deg2rad(float64(point2.Lat))),
  ) * math.Sin(dLong / 2) * math.Sin(dLong / 2)
  c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1 - a))
  d := R * c;
  return d
}

func deg2rad(deg float64) float64 {
  return deg * (math.Pi / 180)
}
