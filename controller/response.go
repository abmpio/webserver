package controller

import (
	"net/http"

	"github.com/abmpio/abmp/pkg/model"
	"github.com/kataras/iris/v12"
)

func HandleError(statusCode int, ctx iris.Context, err error) {
	ctx.StopWithError(statusCode, err)
}

// handle StatusBadRequest
func HandleErrorBadRequest(ctx iris.Context, err error) {
	HandleError(http.StatusBadRequest, ctx, err)
}

func HandleErrorUnauthorized(ctx iris.Context, err error) {
	HandleError(http.StatusUnauthorized, ctx, err)
}

func HandleErrorNotFound(ctx iris.Context, err error) {
	HandleError(http.StatusNotFound, ctx, err)
}

func HandleErrorInternalServerError(ctx iris.Context, err error) {
	HandleError(http.StatusInternalServerError, ctx, err)
}

func HandleSuccess(ctx iris.Context) {
	ctx.StopWithJSON(http.StatusOK, model.NewSuccessResponse())
}

func HandleSuccessWithData(ctx iris.Context, data interface{}) {
	ctx.StopWithJSON(http.StatusOK, model.NewSuccessResponse(func(br *model.BaseResponse) {
		br.SetData(data)
	}))
}

func HandleSuccessWithListData(ctx iris.Context, data interface{}, total int64) {
	ctx.StopWithJSON(http.StatusOK, model.NewSuccessListResponse(data, total))
}

func HandlerBinary(ctx iris.Context, data []byte) (int, error) {
	return ctx.Binary(data)
}

func HandleSuccessWithTableData(ctx iris.Context, list interface{}, total int64, opts ...TableDataOption) {
	d := newDefaultTableData(list, total)
	for _, eachOpt := range opts {
		eachOpt(d)
	}
	HandleSuccessWithData(ctx, d)
}
