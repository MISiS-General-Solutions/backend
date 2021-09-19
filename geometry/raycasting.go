package geometry

import (
	"errors"
	"math"
)

type SegmentType int64

const (
	Obstacle SegmentType = iota
	Road
)

type Vector2d struct {
	V0, V1 Point2D
}

type Segment struct {
	Type SegmentType
	Edge
}

type Intersection struct {
	Coords   Point2D
	Distance float64
}

func GetDistance(a, b Point2D) float64 {
	return math.Sqrt(math.Pow((a.X-b.X), 2) + math.Pow((a.Y-b.Y), 2))
}

func GetAngleOffsetPoint(point Point2D, angle float64) Point2D {
	return Point2D{
		X: point.X + math.Sin(math.Pi/180*angle),
		Y: point.Y + math.Cos(math.Pi/180*angle),
	}
}

func GetIntersectionPoint(ray Vector2d, segment Segment) (Intersection, error) {
	A := segment.V0
	B := segment.V1
	C := ray.V0
	D := ray.V1

	denominator := (D.X-C.X)*(B.Y-A.Y) - (B.X-A.X)*(D.Y-C.Y)
	r := ((B.X-A.X)*(C.Y-A.Y) - (C.X-A.X)*(B.Y-A.Y)) / denominator

	if r < 0 {
		return Intersection{}, errors.New("no intersection")
	}

	s := ((A.X-C.X)*(D.Y-C.Y) - (D.X-C.X)*(A.Y-C.Y)) / denominator

	if s < 0 || s > 1 {
		return Intersection{}, errors.New("no intersection")
	}

	return Intersection{
		Coords:   Point2D{X: s*(B.X-A.X) + A.X, Y: s*(B.Y-A.Y) + A.Y},
		Distance: r,
	}, nil
}

func FindIntersectingRoads(ray Vector2d, segments []Segment) []int64 {
	closest := 100000.0
	for _, segment := range segments {
		if segment.Type != Obstacle {
			continue
		}
		inter, err := GetIntersectionPoint(ray, segment)
		if err != nil {
			continue
		}
		if inter.Distance < closest {
			closest = inter.Distance
		}
	}

	var result []int64
	for _, segment := range segments {
		if segment.Type != Road {
			continue
		}
		inter, err := GetIntersectionPoint(ray, segment)
		if err != nil {
			continue
		}
		if inter.Distance < closest {
			if GetDistance(ray.V0, segment.V0) < GetDistance(ray.V1, segment.V1) {
				result = append(result, segment.ID0)
			} else {
				result = append(result, segment.ID1)
			}
		}
	}

	return result
}
func RayCastFromSlices(center Point2D, obstacles, roads []Edge) []int64 {
	segments := make([]Segment, 0, len(roads)+len(obstacles))
	for i := range obstacles {
		segments = append(segments, Segment{Edge: obstacles[i], Type: Obstacle})
	}
	for i := range roads {
		segments = append(segments, Segment{Edge: roads[i], Type: Road})
	}
	nodes := RayCast(center, segments)
	nodeSlice := make([]int64, 0, len(nodes))
	for k := range nodes {
		nodeSlice = append(nodeSlice, k)
	}
	return nodeSlice
}
func RayCast(center Point2D, segments []Segment) map[int64]bool {
	visible := make(map[int64]bool)
	for angle := 0.; angle < 360-0.00001; angle += 6 {
		offsetPoint := GetAngleOffsetPoint(center, angle)
		roads := FindIntersectingRoads(
			Vector2d{V0: center, V1: offsetPoint},
			segments,
		)
		for _, roadID := range roads {
			visible[roadID] = true
		}
	}
	return visible
}
