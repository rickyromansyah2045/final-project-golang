package main

import "final-project-golang/routes"

func main() {
	r := routes.StartApp()

	r.Run(":8001")
}
