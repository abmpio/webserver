package controller

import (
	"fmt"
	"net/http"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/abmp/pkg/model"
	"github.com/abmpio/webserver/options"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

const (
	log_prefix_webrequest = "(web request) --->"
)

func infoRequestLog(ctx iris.Context, m string, a ...any) {
	info := fmt.Sprintf("%s %s --> url:%s --> requestId:%v ",
		log_prefix_webrequest,
		ctx.Method(),
		getRequestEndpoint(ctx),
		ctx.GetID())
	log.Logger.Info(fmt.Sprintf(info+m, a...),
		zap.String("userId", getUserId(ctx)))
}

func errorRequestLog(ctx iris.Context, m string, a ...any) {
	info := fmt.Sprintf("%s %s ---> url:%s --> requestId:%v ",
		log_prefix_webrequest,
		ctx.Method(),
		getRequestEndpoint(ctx),
		ctx.GetID())
	log.Logger.Warn(fmt.Sprintf(info+m, a...),
		zap.String("userId", getUserId(ctx)))
}

func getUserId(ctx iris.Context) string {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return ""
	}
	return userId
}

func getRequestEndpoint(ctx iris.Context) string {
	request := ctx.Request()
	if request == nil {
		return ""
	}
	return request.RequestURI
}

func HandleError(statusCode int, ctx iris.Context, err error) {
	if options.GetOptions().Log.EnableLogRequest {
		errorRequestLog(ctx, "statusCode:%d,err:%s",
			statusCode,
			err.Error())
	}
	ctx.StopWithError(statusCode, err)
}

// stop with status code
func HandleStopWithStatusCode(statusCode int, ctx iris.Context) {
	if options.GetOptions().Log.EnableLogRequest {
		errorRequestLog(ctx, "statusCode:%d",
			statusCode)
	}
	ctx.StopWithStatus(statusCode)
}

// handle StatusBadRequest
func HandleErrorBadRequest(ctx iris.Context, err error) {
	HandleError(http.StatusBadRequest, ctx, err)
}

func HandleErrorUnauthorized(ctx iris.Context, err error) {
	HandleError(http.StatusUnauthorized, ctx, err)
}

func HandleErrorForbidden(ctx iris.Context) {
	HandleError(http.StatusForbidden, ctx, fmt.Errorf("没有权限访问"))
}

func HandleErrorNotFound(ctx iris.Context, err error) {
	HandleError(http.StatusNotFound, ctx, err)
}

func HandleErrorInternalServerError(ctx iris.Context, err error) {
	HandleError(http.StatusInternalServerError, ctx, err)
}

func HandleResponseWith(ctx iris.Context, opts ...func(*model.BaseResponse)) {
	statusCode := http.StatusOK
	if options.GetOptions().Log.EnableLogRequest {
		infoRequestLog(ctx, "statusCode:%d",
			statusCode)
	}
	ctx.StopWithJSON(statusCode, model.NewSuccessResponse(opts...))
}

func HandleNotSuccess(ctx iris.Context, opts ...func(*model.BaseResponse)) {
	statusCode := http.StatusOK
	if options.GetOptions().Log.EnableLogRequest {
		infoRequestLog(ctx, "statusCode:%d",
			statusCode)
	}

	ctx.StopWithJSON(statusCode, model.NewErrorResponse(opts...))
}

func HandleSuccess(ctx iris.Context) {
	statusCode := http.StatusOK
	if options.GetOptions().Log.EnableLogRequest {
		infoRequestLog(ctx, "statusCode:%d",
			statusCode)
	}

	ctx.StopWithJSON(statusCode, model.NewSuccessResponse())
}

func HandleSuccessWithData(ctx iris.Context, data interface{}) {
	statusCode := http.StatusOK
	if options.GetOptions().Log.EnableLogRequest {
		infoRequestLog(ctx, "statusCode:%d",
			statusCode)
	}
	ctx.StopWithJSON(statusCode, model.NewSuccessResponse(func(br *model.BaseResponse) {
		br.SetData(data)
	}))
}

func HandleSuccessWithListData(ctx iris.Context, data interface{}, total int64) {
	statusCode := http.StatusOK
	if options.GetOptions().Log.EnableLogRequest {
		infoRequestLog(ctx, "statusCode:%d",
			statusCode)
	}
	ctx.StopWithJSON(statusCode, model.NewSuccessListResponse(data, total))
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
