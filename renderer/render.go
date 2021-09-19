package renderer

import (
	"fmt"
	"image/color"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

func Render(path string, nw, se s2.LatLng, camera s2.LatLng, obstacles []s2.Edge, roads []s2.Edge) error {
	ctx := sm.NewContext()
	ctx.SetTileProvider(sm.NewTileProviderStamenTerrain())
	rect, err := sm.CreateBBox(nw.Lat.Degrees(), nw.Lng.Degrees(), se.Lat.Degrees(), se.Lng.Degrees())
	if err != nil {
		return err
	}
	center := rect.Center()
	_ = center

	ctx.SetBoundingBox(*rect)
	// ctx.SetCenter(center)
	// ctx.SetZoom(15)

	for _, edge := range obstacles {
		ctx.AddObject(sm.NewCircle(s2.LatLngFromPoint(edge.V0), color.Black, color.White, 8, 0.5))
		ctx.AddObject(sm.NewCircle(s2.LatLngFromPoint(edge.V1), color.Black, color.White, 8, 0.5))
	}
	for _, edge := range roads {
		ctx.AddObject(sm.NewCircle(s2.LatLngFromPoint(edge.V0), color.RGBA64{0, 0, 65535, 65535}, color.RGBA64{0, 0, 65535, 32767}, 6, 0.2))
		ctx.AddObject(sm.NewCircle(s2.LatLngFromPoint(edge.V1), color.RGBA64{0, 0, 65535, 65535}, color.RGBA64{0, 0, 65535, 32767}, 6, 0.2))
	}
	ctx.AddObject(sm.NewCircle(camera, color.RGBA64{65535, 0, 0, 65535}, color.RGBA64{65535, 0, 0, 65535}, 10, 1))

	ctx.SetTileProvider(sm.NewTileProviderCartoDark())
	img, err := ctx.Render()
	if err != nil {
		return err
	}
	fmt.Println("rendered")

	if err := gg.SavePNG(fmt.Sprint(path), img); err != nil {
		return err
	}
	return nil
}
func RenderAffected(path string, nw, se s2.LatLng, camera s2.LatLng, affected []s2.LatLng) error {
	ctx := sm.NewContext()
	ctx.SetTileProvider(sm.NewTileProviderStamenTerrain())
	rect, err := sm.CreateBBox(nw.Lat.Degrees(), nw.Lng.Degrees(), se.Lat.Degrees(), se.Lng.Degrees())
	if err != nil {
		return err
	}
	center := rect.Center()
	_ = center

	ctx.SetBoundingBox(*rect)
	// ctx.SetCenter(center)
	// ctx.SetZoom(15)

	for _, poi := range affected {
		ctx.AddObject(sm.NewCircle(poi, color.RGBA64{0, 0, 65535, 65535}, color.RGBA64{0, 0, 65535, 32767}, 6, 0.2))
	}
	ctx.AddObject(sm.NewCircle(camera, color.RGBA64{65535, 0, 0, 65535}, color.RGBA64{65535, 0, 0, 65535}, 10, 1))

	ctx.SetTileProvider(sm.NewTileProviderCartoDark())
	img, err := ctx.Render()
	if err != nil {
		return err
	}
	fmt.Println("rendered")

	if err := gg.SavePNG(fmt.Sprint(path), img); err != nil {
		return err
	}
	return nil
}
