package api

import (
	"MGS/osmdata"
	"MGS/routing"
	"MGS/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	Roads *osmdata.OsmData
)

func Serve() {
	router := gin.Default()

	router.GET("/mapping/find_route", rougingHandler)

	router.Run()
}

type RouteRequest struct {
	Targets []routing.LatLonPair `json:"targets"`
}

func rougingHandler(ctx *gin.Context) {
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
