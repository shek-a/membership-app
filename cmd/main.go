package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"members.com/membership/internal/database"
	"members.com/membership/internal/routes"
	"members.com/membership/pkg/handler"
	"members.com/membership/pkg/repository"
	"members.com/membership/pkg/service"
)

func main() {
	mongoConnection, err := database.ConnectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	server := gin.Default()

	memberService := service.NewMemberService(repository.NewMembershipRepository(mongoConnection))
	MemberHandler := handler.NewMemberHandler(server, memberService)

	routes.RegisterRoutes(server, MemberHandler)

	server.Run(":8080")
}
