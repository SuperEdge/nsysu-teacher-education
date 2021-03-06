package route

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	"github.com/nsysu/teacher-education/src/handler"
	"github.com/nsysu/teacher-education/src/middleware"
	"github.com/nsysu/teacher-education/src/utils/config"
)

var addr = fmt.Sprintf("%v:%v", config.Get("server.host"), config.Get("server.port"))

// Run maps the routing path and keeps listening for request
func Run() {
	app := iris.New()

	// CORS
	app.AllowMethods(iris.MethodOptions)
	app.Use(middleware.CorsMiddleware)

	app.Get("/", hero.Handler(handler.HelloHandler))

	v1 := app.Party("/v1", hero.Handler(middleware.AuthMiddleware))
	{
		v1.Post("/login", hero.Handler(handler.LoginHandler))            // 登入
		v1.Post("/logout", hero.Handler(handler.LogoutHandler))          // 登出
		v1.Post("/renew-token", hero.Handler(handler.RenewTokenHandler)) // 取得新的 access token

		user := v1.Party("/user")
		{
			user.Get("/", hero.Handler(handler.GetStudentInformationHandler))
		}

		users := v1.Party("/users")
		{
			users.Get("/", hero.Handler(handler.GetStudentsHandler))     // 取得學生列表
			users.Post("/", hero.Handler(handler.CreateStudentsHandler)) // 新增學生帳號
		}

		course := v1.Party("/course")
		{
			course.Get("/", hero.Handler(handler.GetCourseHandler))                         // 取得講座列表
			course.Get("/{courseID}", hero.Handler(handler.GetCourseInformationHandler))    // 取得講座資訊
			course.Post("/", hero.Handler(handler.CreateCourseHandler))                     // 新增講座
			course.Post("/sign-up", hero.Handler(handler.SignUpCourseHandler))              // 報名講座
			course.Patch("/{courseID}", hero.Handler(handler.UpdateCourseHandler))          // 修改講座
			course.Delete("/{courseID}", hero.Handler(handler.DeleteCourseHandler))         // 刪除講座
			course.Patch("/show/{courseID}", hero.Handler(handler.UpdateCourseShowHandler)) // 修改講座顯示狀態

			student := course.Party("/student")
			{
				student.Get("/", hero.Handler(handler.GetStudentCourseHandler))                  // 取得報名講座列表
				student.Patch("/review", hero.Handler(handler.UpdateStudentCourseReviewHandler)) // 上傳心得
				student.Patch("/status", hero.Handler(handler.UpdateStudentCourseStatusHandler)) // 審核
			}
		}

		serviceLearning := v1.Party("/service-learning")
		{
			serviceLearning.Get("/", hero.Handler(handler.GetServiceLearningHandler))                                  // 取得服務學習列表
			serviceLearning.Post("/", hero.Handler(handler.CreateServiceLearningHandler))                              // 新增服務學習
			serviceLearning.Post("/sign-up", hero.Handler(handler.SignUpServiceLearningHandler))                       // 報名服務學習
			serviceLearning.Patch("/", hero.Handler(handler.UpdateServiceLearningHandler))                             // 修改服務學習
			serviceLearning.Delete("/{serviceLearnginID}", hero.Handler(handler.DeleteServiceLearningHandler))         // 刪除服務學習
			serviceLearning.Patch("/show/{serviceLearnginID}", hero.Handler(handler.UpdateServiceLearningShowHandler)) // 修改服務學習顯示狀態

			student := serviceLearning.Party("/student")
			{
				student.Get("/", hero.Handler(handler.GetStudentsServiceLearningHandler))                 // 報名服務學習列表
				student.Get("/{file}", hero.Handler(handler.GetStudentServiceLearningFileHandler))        // 取得佐證資料或心得
				student.Patch("/", hero.Handler(handler.UpdateStudentServiceLearningHandler))             // 上傳資料
				student.Patch("/status", hero.Handler(handler.UpdateStudentServiceLearningStatusHandler)) // 審核
			}
		}

		lecture := v1.Party("/lecture")
		{
			lecture.Get("/", hero.Handler(handler.GetLecturesHandler))                 // 取得課程列表
			lecture.Get("/{lectureID}", hero.Handler(handler.GetLectureDetailHandler)) // 取得課程詳細資料
			lecture.Post("/sign-up", hero.Handler(handler.SignUpLectureHandler))       // 報名課程

			student := lecture.Party("/student")
			{
				student.Get("/", hero.Handler(handler.GetStudentLectureHandler))                                // 取得報名課程列表
				student.Get("/detail/{studentLectureID}", hero.Handler(handler.GetStudentLectureDetailHandler)) // 取得學生課程詳細資料
				student.Patch("/subject", hero.Handler(handler.UpdateStudentSubjectHandler))                    // 更新學生成績
				student.Patch("/pass", hero.Handler(handler.UpdateStudentLectureStatusHandler))                 // 更新狀態
			}
		}
	}

	app.Run(iris.Addr(addr))
}
