package routes

import (
	"CfBE/controllers"
	middleware "CfBE/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	r := gin.Default()
	// for auth
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.POST("/refresh-token", controllers.RefreshToken)
	// ------------------------------------------------------------------------
	// for in-user
	r.GET("/users", controllers.Users)
	r.GET("/user/:id1", controllers.User)
	// ------------------------------------------------------------------------
	// for contest
	r.POST("/contest", middleware.AdminRoleRequired(), controllers.CreateContest)        // create a new contest
	r.DELETE("/contest/:id2", middleware.AdminRoleRequired(), controllers.DeleteContest) // delete a contest
	r.GET("/contests", controllers.Contests)                                             // return all contest
	r.GET("/contest/:id2", controllers.Contest)                                          // return a contest
	//r.POST("/contest/:id/participant", middleware.ParticipantRoleRequired(), controllers.Participant)

	// ------------------------------------------------------------------------
	// for problem
	r.POST("/contest/:id2/problem", middleware.AdminRoleRequired(), controllers.CreateProblem)        // add a new problem to a contest                                        // create a new problem
	r.DELETE("/contest/:id2/problem/:id3", middleware.AdminRoleRequired(), controllers.DeleteProblem) // delete a problem
	r.GET("/contest/:id2/problems", controllers.Problems)                                             // return all problems
	r.GET("/contest/:id2/problem/:id3", controllers.Problem)                                          // return a problem

	// ------------------------------------------------------------------------
	// for submission
	r.POST("/contest/:id2/problem/:id3/submit/:id1", middleware.SubmitAuth(), controllers.SubmitCode)
	r.GET("/contest/:id2/problem/:id3/submission/:id4", controllers.GetSubmission)
	/*
		r.GET("/contest/:id/leaderboard", controllers.Leaderboard)
		// ----------------------------------------------------------------------
		// for submission
		r.GET("/contest/problem/:id/submission", controllers.Submission)
		// ------------------------------------------------------------------------
		// for admin
		r.GET("/admin", controllers.admin)
	*/
	return r
}
