package main

import (
	"context"
	"fmt"
	"github.com/passionde/user-segmentation-service/config"
	"github.com/passionde/user-segmentation-service/internal/repo"
	"github.com/passionde/user-segmentation-service/internal/service"
	"github.com/passionde/user-segmentation-service/pkg/postgres"
	"github.com/passionde/user-segmentation-service/pkg/secure"
	log "github.com/sirupsen/logrus"
	"os"
)

const configPath = "config/config.yaml"

func main() {
	// Args
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli exist <KeyApi>|generate")
	}

	cmd := os.Args[1]
	if cmd != "exist" && cmd != "generate" {
		log.Fatal("Usage: cli exist <KeyApi>|generate")
	}
	if cmd == "exist" && len(os.Args) != 3 {
		log.Fatal("Usage: cli exist <KeyApi>")
	}

	// Config
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Postgresql
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}
	defer pg.Close()

	repositories := repo.NewRepositories(pg)

	// Services
	deps := service.ServicesDependencies{
		Repos:     repositories,
		APISecure: secure.NewSecure(cfg.Secure.Salt),
	}
	services := service.NewServices(deps)

	// Handlers
	switch cmd {
	case "exist":
		existCommand(services, os.Args[2])
	case "generate":
		generateCommand(services)
	}
}

func existCommand(services *service.Services, token string) {
	id, err := services.Auth.TokenExist(context.TODO(), token)
	if err != nil {
		fmt.Printf("ApiKey <%s> - does not exist\n", token)
		return
	}
	fmt.Printf("ApiKey <%s> - exist, ID = %d\n", token, id)
}

func generateCommand(services *service.Services) {
	_, key, err := services.Auth.GenerateToken(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Api key: Bearer %s\n", key)
}
