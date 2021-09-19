package osmdata

import (
	"MGS/geometry"
	"MGS/routing"
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang/geo/s2"
)

type Camera struct {
	Coords []float64 `json:"coords"`
	ID     int64     `json:"id"`
	Img    string    `json:"img"`
}

func AddCamerasFromFile(path string, roadData, obstacleData *OsmData) error {

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var cameras []Camera
	if err = json.Unmarshal(b, &cameras); err != nil {
		return err
	}
	for i, camera := range cameras {
		center, err := geometry.StickWithAcceleration(obstacleData.Shi, obstacleData.Shapes, s2.LatLngFromDegrees(camera.Coords[0], camera.Coords[1]))
		if err != nil {
			return err
		}

		closeObstacles := geometry.GetEdgesInDistance(obstacleData.Shi, obstacleData.Shapes, center, geometry.Meter*30)
		closeRoads := geometry.GetEdgesInDistance(roadData.Shi, roadData.Shapes, center, geometry.Meter*60)

		flatCoords := geometry.FlatApproxRoadsAndObstacles(center, closeObstacles, closeRoads, geometry.RegionMoscow, obstacleData.NodeIndex, roadData.NodeIndex)

		affectedRoads := geometry.RayCastFromSlices(geometry.Point2D{}, flatCoords.Obstacles, flatCoords.Roads)
		for _, ID := range affectedRoads {
			roadData.Nodes[ID].Tags = append(roadData.Nodes[ID].Tags, routing.CameraTag)
		}

		if i%25 == 0 {
			fmt.Printf("%v of %v\n", i, len(cameras))
		}
	}
	return nil
}
