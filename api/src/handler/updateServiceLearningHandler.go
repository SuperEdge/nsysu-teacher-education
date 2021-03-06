package handler

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/kataras/iris"
	"github.com/nsysu/teacher-education/src/errors"
	"github.com/nsysu/teacher-education/src/service"
	t "github.com/nsysu/teacher-education/src/utils/time"
	"github.com/nsysu/teacher-education/src/utils/typecast"
)

// UpdateServiceLearningHandler update service-learning
func UpdateServiceLearningHandler(ctx iris.Context) {
	type rule struct {
		ServiceLearningID string `valid:"required"`
		Type              string `valid:"required, in(internship|volunteer|both)"`
		Content           string `valid:"required, length(0|150)"`
		Session           string `valid:"required, length(0|36)"`
		Hours             uint   `valid:"required, int"`
	}

	loc, _ := time.LoadLocation("Asia/Taipei")
	startTime, err := time.ParseInLocation(t.Date, ctx.FormValue("Start"), loc)
	if err != nil {
		failed(ctx, errors.ValidateError("Start: "+ctx.FormValue("Start")+" does not validate as timestamp"))
		return
	}
	endTime, err := time.ParseInLocation(t.Date, ctx.FormValue("End"), loc)
	if err != nil {
		failed(ctx, errors.ValidateError("End: "+ctx.FormValue("Start")+" does not validate as timestamp"))
		return
	}
	if !startTime.Before(endTime) {
		failed(ctx, errors.ValidateError("Start: "+ctx.FormValue("Start")+" does not before "+ctx.FormValue("End")))
		return
	}

	params := &rule{
		ServiceLearningID: ctx.FormValue("ServiceLearningID"),
		Type:              ctx.FormValue("Type"),
		Content:           ctx.FormValue("Content"),
		Session:           ctx.FormValue("Session"),
		Hours:             typecast.StringToUint(ctx.FormValue("Hours")),
	}

	if _, err := govalidator.ValidateStruct(params); err != nil {
		failed(ctx, errors.ValidateError(err.Error()))
		return
	}

	result, err := service.UpdateServieLearning(params.ServiceLearningID, params.Type, params.Content, params.Session, params.Hours, startTime, endTime)

	if err != (*errors.Error)(nil) {
		failed(ctx, err)
		return
	}

	success(ctx, result)
	return
}
