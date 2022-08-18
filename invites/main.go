package main

import (
	"github.com/micro/distributed-api/invites/handler"
	pb "github.com/micro/distributed-api/invites/proto"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dbAddress = "postgresql://postgres:postgres@localhost:5432/invites?sslmode=disable"

func main() {
	// Create service
	srv := service.New(
		service.Name("invites"),
		service.Version("latest"),
	)

	// Connect to the database
	cfg, err := config.Get("invites.database")
	if err != nil {
		logger.Fatalf("Error loading config: %v", err)
	}
	addr := cfg.String(dbAddress)
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "invites_",
		},
	})
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	if err := db.AutoMigrate(&handler.Invite{}); err != nil {
		logger.Fatalf("Error migrating database: %v", err)
	}

	// Register handler
	pb.RegisterInvitesHandler(srv.Server(), &handler.Invites{DB: db})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
