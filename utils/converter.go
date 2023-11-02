package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/jftuga/geodist"
	"strconv"
)

func LocationToKM(c *gin.Context, userLatitude, userLongitude, postLatitude, postLongitude string) string {
	postLat, err := StringToFloat64(postLatitude)
	if err != nil {
		HttpRespFailed(c, 400, err.Error())
		return err.Error()
	}
	postLng, err := StringToFloat64(postLongitude)
	if err != nil {
		HttpRespFailed(c, 400, err.Error())
		return err.Error()
	}
	post := geodist.Coord{postLat, postLng}

	userLat, err := StringToFloat64(userLatitude)
	if err != nil {
		HttpRespFailed(c, 400, err.Error())
		return err.Error()
	}

	userLng, err := StringToFloat64(userLongitude)
	if err != nil {
		HttpRespFailed(c, 400, err.Error())
		return err.Error()
	}

	user := geodist.Coord{userLat, userLng}

	_, km, _ := geodist.VincentyDistance(user, post)
	kmStr := Float64ToString(km)

	return kmStr
}

func StringToFloat64(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
