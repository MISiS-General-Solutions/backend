package geometry

import (
	"encoding/json"
	"errors"
	"math"
	"os"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

const (
	RegionMoscow = 1
	Meter        = 7.7446506001348e-13
	cosineMoscow = 0.562576257
)

type Point2D struct {
	X, Y float64
}
type Edge struct {
	V0, V1   Point2D
	ID0, ID1 int64
}
type S2EdgeWithID struct {
	Edge s2.Edge
}
type FlatCoords struct {
	Center    s2.LatLng
	Obstacles []Edge
	Roads     []Edge
}

func MeasureMeridian() float64 {
	p1 := s2.LatLngFromDegrees(0, 0)
	p2 := s2.LatLngFromDegrees(0, 180)
	return float64(s1.ChordAngleFromAngle(p1.Distance(p2)))
}
func SaveFlatJSON(path string, fc FlatCoords) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	b, err := json.Marshal(fc)
	if err != nil {
		return err
	}
	_, err = file.Write(b)
	return err
}
func FlatApproxRoadsAndObstacles(center s2.LatLng, obstacles []s2.Edge, roads []s2.Edge, region int, nodeIndexO, nodeIndexR map[s2.Point]int64) FlatCoords {
	var cosine float64
	if region != RegionMoscow {
		cosine = cosineMoscow
	} else {
		cosine = math.Cos(center.Lat.Radians())
	}
	return FlatCoords{
		Center:    center,
		Obstacles: FlatApprox(center, obstacles, cosine, nodeIndexO),
		Roads:     FlatApprox(center, roads, cosine, nodeIndexR),
	}
}
func FlatApprox(center s2.LatLng, edges []s2.Edge, cosine float64, nodeIndex map[s2.Point]int64) []Edge {
	res := make([]Edge, len(edges))

	i := 0
	for _, edge := range edges {
		v1 := s2.LatLngFromPoint(edge.V0)
		v2 := s2.LatLngFromPoint(edge.V1)

		res[i] = Edge{
			V0: Point2D{
				Y: 111111 * (v1.Lat.Degrees() - center.Lat.Degrees()),
				X: 111111 * (v1.Lng.Degrees() - center.Lng.Degrees()) * cosine,
			},
			V1: Point2D{
				Y: 111111 * (v2.Lat.Degrees() - center.Lat.Degrees()),
				X: 111111 * (v2.Lng.Degrees() - center.Lng.Degrees()) * cosine,
			},
			ID0: nodeIndex[edge.V0],
			ID1: nodeIndex[edge.V1],
		}

		i += 1
	}
	return res
}
func EuclidianDistanceApprox(aLat, bLat, aLng, bLng float64, region int) float64 {
	var cosine float64
	if region != RegionMoscow {
		cosine = cosineMoscow
	} else {
		cosine = math.Cos(aLat / 180)
	}
	dy := 111111 * (aLat - bLat)
	dx := 111111 * (aLng - bLng) * cosine
	return math.Sqrt(dx*dx + dy*dy)
}
func GetEdgesInDistance(shi *s2.ShapeIndex, shapes map[int32]s2.Polyline, center s2.LatLng, dist s1.ChordAngle) []s2.Edge {
	opts := s2.NewClosestEdgeQueryOptions().DistanceLimit(dist)
	edgeQuery := s2.NewClosestEdgeQuery(shi, opts)
	edgeResult := edgeQuery.FindEdges(s2.NewMinDistanceToPointTarget(s2.PointFromLatLng(center)))

	closePoints := make([]s2.Edge, len(edgeResult))
	for i := range edgeResult {
		poly := shapes[edgeResult[i].ShapeID()]
		edge := poly.Edge(int(edgeResult[i].EdgeID()))
		closePoints[i] = edge
	}
	return closePoints
}
func GetClosestNode(shi *s2.ShapeIndex, shapes map[int32]s2.Polyline, nodeIndex map[s2.Point]int64, origin s2.Point) (int64, error) {

	opts := s2.NewClosestEdgeQueryOptions().MaxResults(1)
	edgeQuery := s2.NewClosestEdgeQuery(shi, opts)
	edgeResult := edgeQuery.FindEdges(s2.NewMinDistanceToPointTarget(origin))
	if len(edgeResult) != 1 {
		return 0, errors.New("no nodes found")
	}
	poly := shapes[edgeResult[0].ShapeID()]
	edge := poly.Edge(int(edgeResult[0].EdgeID()))
	if edge.V0.Distance(origin) < edge.V1.Distance(origin) {
		return nodeIndex[edge.V0], nil
	}
	return nodeIndex[edge.V1], nil
}
