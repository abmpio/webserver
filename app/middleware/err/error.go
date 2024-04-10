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

	responseData := ctx.Recorder().Body()
	statusCode := ctx.GetStatusCode()
	if context.StatusCodeNotSuccessful(statusCode) && !v.responseIsIgnore(responseData) {
		ctx.Recorder().ResetBody()
		err := ctx.GetErr()
		// responseDataString := string(responseData)
		ctx.StopWithJSON(statusCode, model.NewErrorResponse(func(br *model.BaseResponse) {
			if err != nil {
				br.SetMessage(err.Error())
			}
		}))
		// log.Logger.Error(responseDataString)
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
