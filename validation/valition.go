package validation

import (
	"goravel/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValideParamId(idParams string, ctx *gin.Context) (int, bool) {
	if idParams == "" {
		response := model.Response{
			Status:  http.StatusBadRequest,
			Message: "ID is required",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return 0, false
	}

	id, err := strconv.Atoi(idParams)
	if err != nil {
		response := model.Response{
			Status:  http.StatusBadRequest,
			Message: "ID must be a number",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return 0, false
	}
	return id, true
}
