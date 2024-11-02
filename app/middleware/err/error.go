package err

import (
	"encoding/json"

	"github.com/kataras/iris/v12/context"

	"github.com/abmpio/abmp/pkg/model"
)

type errWrapperMiddleware struct {
}

func New() context.Handler {
	v := &errWrapperMiddleware{}
	return v.ServeHTTP
}

func (v *errWrapperMiddleware) ServeHTTP(ctx *context.Context) {
	ctx.Record()
	ctx.Next()

	statusCode := ctx.GetStatusCode()
	if context.StatusCodeNotSuccessful(statusCode) {
		responseData := ctx.Recorder().Body()
		if !v.responseIsIgnore(responseData) {
			ctx.Recorder().ResetBody()
			err := ctx.GetErr()
			ctx.StopWithJSON(statusCode, model.NewErrorResponse(func(br *model.BaseResponse) {
				br.Code = statusCode
				br.Status = model.HttpResponseMessageError
				var message string
				if err != nil {
					message = err.Error()
				} else {
					message = string(responseData)
				}
				if len(message) > 0 {
					br.SetMessage(message)
				}
			}))
		}
	}
}

func (v *errWrapperMiddleware) responseIsIgnore(responseData []byte) bool {
	if len(responseData) <= 0 {
		return false
	}
	baseResponse := &model.BaseResponse{}
	marshalErr := json.Unmarshal(responseData, baseResponse)
	return marshalErr == nil
}
