package api

import (
	"MGS/osmdata"
	"MGS/routing"
	"MGS/shared"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	Roads *osmdata.OsmData
)

func Serve() {
	router := gin.Default()

	router.GET("/mapping/find_route", routingHandler)
	router.GET("/cameras/get", cameraHandler)

	router.Run()
}

type RouteRequest struct {
	Targets []routing.LatLonPair `json:"targets"`
}

func routingHandler(ctx *gin.Context) {
	var request []routing.LatLonPair
	if err := ctx.BindJSON(&request); err != nil {
		return
	}
	if len(request) < 2 {
		ctx.JSON(http.StatusBadRequest, "need at least 2 points")
	}

	path := routing.GetRouteFromLatLon(request[0], request[1], Roads.Nodes, Roads.Shi, Roads.Shapes, Roads.NodeIndex)
	ctx.JSON(http.StatusOK, shared.MustMarshallToJSON(path))
}
func cameraHandler(ctx *gin.Context) {
	b, err := os.ReadFile("data/cache/cameras.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
	}
	ctx.JSON(http.StatusOK, b)
}
