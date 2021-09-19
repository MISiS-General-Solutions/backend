package osmdata

import (
	"MGS/routing"
	"context"
	"encoding/gob"
	"io"
	"os"

	"github.com/golang/geo/s2"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmxml"
)

type OsmData struct {
	Nodes     map[int64]*routing.Node
	Ways      map[int64]*osm.Way
	Relations map[int64]*osm.Relation

	NodeIndex map[s2.Point]int64
	Shi       *s2.ShapeIndex
	Shapes    map[int32]s2.Polyline
}

func NewOsmData() *OsmData {
	res := &OsmData{
		Nodes:     make(map[int64]*routing.Node),
		Ways:      make(map[int64]*osm.Way),
		Relations: make(map[int64]*osm.Relation),

		NodeIndex: make(map[s2.Point]int64),
		Shapes:    make(map[int32]s2.Polyline),
	}
	res.Shi = s2.NewShapeIndex()
	return res
}
func ReadOSM(r io.Reader) (*OsmData, error) {
	scanner := osmxml.New(context.Background(), r)
	defer scanner.Close()

	res := NewOsmData()
	for scanner.Scan() {
		o := scanner.Object()
		switch obj := o.(type) {
		case *osm.Node:
			res.Nodes[int64(obj.ID)] = &routing.Node{Node: obj}
		case *osm.Way:
			res.Ways[int64(obj.ID)] = obj
		case *osm.Relation:
			res.Relations[int64(obj.ID)] = obj
		}
	}

	err := scanner.Err()
	if err != nil {
		return nil, err
	}
	res.Compile()
	return res, nil
}
func (r *OsmData) Compile() {
	r.Shapes = make(map[int32]s2.Polyline, len(r.Ways))

	for _, way := range r.Ways {
		points := make([]s2.Point, len(way.Nodes))
		for i, nodeShallow := range way.Nodes {
			node := r.Nodes[int64(nodeShallow.ID)]
			node.Tags = append(node.Tags, way.Tags...)
			points[i] = s2.PointFromLatLng(s2.LatLngFromDegrees(node.Lat, node.Lon))
			r.NodeIndex[points[i]] = int64(node.ID)
		}
		poly := s2.Polyline(points)
		idx := r.Shi.Add(&poly)
		r.Shapes[idx] = poly
	}
}
func ReadOSMFromFile(path string) (*OsmData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadOSM(f)
}
func (r *OsmData) CompileGraph() {
	for i := range r.Ways {
		for j := 1; j < len(r.Ways[i].Nodes); j++ {
			prev := r.Nodes[int64(r.Ways[i].Nodes[j-1].ID)]
			cur := r.Nodes[int64(r.Ways[i].Nodes[j].ID)]

			for k := range cur.Neighbours {
				if cur.Neighbours[k] == prev {
					goto skipAdd1
				}
			}
			cur.Neighbours = append(cur.Neighbours, prev)
		skipAdd1:

			for k := range prev.Neighbours {
				if prev.Neighbours[k] == cur {
					goto skipAdd2
				}
			}
			prev.Neighbours = append(prev.Neighbours, cur)
		skipAdd2:
		}
	}
}
func (r *OsmData) SaveGOB(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	if err = enc.Encode(*r); err != nil {
		return err
	}
	return nil
}

func LoadGOBFromFile(path string) (*OsmData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(file)
	var res OsmData
	if err = dec.Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}
