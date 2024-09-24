package main

import (
	"CfBE/initializers"
	"CfBE/routes"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectDB()
}
func main() {

	r := routes.SetupRoutes()
	r.Run()
}
