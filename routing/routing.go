package routing

import (
	"MGS/geometry"
	"errors"

	"github.com/beefsack/go-astar"
	"github.com/golang/geo/s2"
	"github.com/paulmach/osm"
)

type CostProfile struct {
}

var (
	CameraTag = osm.Tag{Key: "camera:type", Value: "fixed"}
)

func PathCost(n1, n2 *Node) float64 {
	koef := 1.0
	for _, tag := range n1.Tags {
		if tag == CameraTag {
			koef *= 5
		}
	}
	for _, tag := range n2.Tags {
		if tag == CameraTag {
			koef *= 5
		}
	}
	return koef * geometry.EuclidianDistanceApprox(n1.Lat, n2.Lat, n1.Lon, n2.Lon, geometry.RegionMoscow)
}

type Node struct {
	*osm.Node
	Neighbours []astar.Pather
}

func (r *Node) PathNeighbors() []astar.Pather {
	return r.Neighbours
}

func (r *Node) PathNeighborCost(to astar.Pather) float64 {
	n2, ok := to.(*Node)
	if !ok {
		panic(errors.New("wrong type"))
	}
	return PathCost(r, n2)
}

func (r *Node) PathEstimatedCost(to astar.Pather) float64 {
	n2, ok := to.(*Node)
	if !ok {
		panic(errors.New("wrong type"))
	}
	return PathCost(r, n2)
}

type LatLonPair struct {
	Lat, Lon float64
}
type Path struct {
	Coords   []LatLonPair
	Distance float64
}

func GetRouteFromLatLon(a, b LatLonPair, nodes map[int64]*Node, shi *s2.ShapeIndex, shapes map[int32]s2.Polyline, nodeIndex map[s2.Point]int64) Path {
	aInd, err := geometry.GetClosestNode(shi, shapes, nodeIndex, s2.PointFromLatLng(s2.LatLngFromDegrees(a.Lat, a.Lon)))
	if err != nil {
		return Path{}
	}
	bInd, err := geometry.GetClosestNode(shi, shapes, nodeIndex, s2.PointFromLatLng(s2.LatLngFromDegrees(b.Lat, b.Lon)))
	if err != nil {
		return Path{}
	}
	return GetRoute(aInd, bInd, nodes)
}
func GetRoute(a, b int64, nodes map[int64]*Node) Path {
	path, distance, found := astar.Path(nodes[a], nodes[b])
	if !found {
		return Path{}
	}
	res := make([]LatLonPair, len(path))
	for i := range path {
		node, ok := path[i].(*Node)
		if !ok {
			panic(errors.New("wrong type"))
		}
		res[i] = LatLonPair{node.Lat, node.Lon}
	}
	return Path{Coords: res, Distance: distance}
}

// func ReadPythonLatLons(path string) ([]s2.LatLng, error) {
// 	b, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	bb := bytes.Split(b, []byte{'('})[1:]
// 	res := make([]s2.LatLng, len(bb))
// 	for i, coords := range bb {
// 		s := string(coords)
// 		_ = s
// 		sepPos := bytes.IndexByte(coords, ',')
// 		if sepPos == -1 {
// 			return nil, errors.New("incorrect format")
// 		}
// 		lat, err := strconv.ParseFloat(string(coords[:sepPos]), 64)
// 		if err != nil {
// 			return nil, err
// 		}
// 		lon, err := strconv.ParseFloat(string(coords[sepPos+2:len(coords)-1]), 64)
// 		if err != nil {
// 			return nil, err
// 		}
// 		res[i] = s2.LatLngFromDegrees(lat, lon)
// 	}
// 	return res, nil
// }
