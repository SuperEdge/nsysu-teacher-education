package handler

import (
	"github.com/kataras/iris"
	"github.com/nsysu/teacher-education/src/error"
	"github.com/nsysu/teacher-education/src/service"
	"github.com/nsysu/teacher-education/src/utils/auth"
)

// LogoutHandler user logout
func LogoutHandler(ctx iris.Context) {
	result, err := service.Logout(auth.Account(ctx))

	if err != (*error.Error)(nil) {
		failed(ctx, err)
		return
	}

	success(ctx, result)
	return
}
