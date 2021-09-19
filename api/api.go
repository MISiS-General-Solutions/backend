package api

import (
	"MGS/osmdata"
	"MGS/routing"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	Roads *osmdata.OsmData
)

func Serve() {
	router := gin.Default()

	router.POST("/mapping/find_route", routingHandler)
	router.POST("/cameras/get", cameraHandler)

	router.Run(":80")
}

func routingHandler(ctx *gin.Context) {
	var request []routing.LatLonPair
	if err := ctx.BindJSON(&request); err != nil {
		return
	}
	if len(request) < 2 {
		ctx.JSON(http.StatusBadRequest, "need at least 2 points")
	}
	var fullPath routing.Path
	for i := 1; i < len(request); i++ {
		path := routing.GetRouteFromLatLon(request[i-1], request[i], Roads.Nodes, Roads.Shi, Roads.Shapes, Roads.NodeIndex)
		fullPath.Coords = append(fullPath.Coords, path.Coords...)
		fullPath.Distance += path.Distance
	}

	ctx.JSON(http.StatusOK, fullPath)
}
func cameraHandler(ctx *gin.Context) {
	b, err := os.ReadFile("data/cache/cameras.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
	}
	ctx.Data(http.StatusOK, "application/json", b)
}
