package main

import (
	"MGS/api"
	"MGS/client"
	"MGS/osmdata"
	"fmt"

	"github.com/golang/geo/s2"
)

func initData(nw, se s2.LatLng) (*osmdata.OsmData, error) {

	settings := fmt.Sprintf("[timeout:120][bbox:%v,%v,%v,%v];", se.Lat.Degrees(), nw.Lng.Degrees(), nw.Lat.Degrees(), se.Lng.Degrees())

	cli := client.NewClient()
	cli.SetSettings(settings)

	obstacles, err := cli.GetObstacles()
	if err != nil {
		return nil, err
	}
	fmt.Println("got obstacles")
	_ = obstacles
	roads, err := cli.GetRoads()
	if err != nil {
		return nil, err
	}
	fmt.Println("got roads")

	if err = osmdata.AddCamerasFromFile("data/cache/cameras.json", roads, obstacles); err != nil {
		return nil, err
	}

	roads.CompileGraph()

	return roads, nil
}
func main() {
	nw := s2.LatLngFromDegrees(55.767944, 37.600997)
	se := s2.LatLngFromDegrees(55.763467, 37.615685)

	roads, err := initData(nw, se)
	if err != nil {
		panic(err)
	}
	fmt.Println("init done")

	api.Roads = roads
	api.Serve()
}
