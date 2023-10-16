package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cmd/main.go/internal/service/user_request"

	"cmd/main.go/internal/repo"

	"cmd/main.go/internal/config"
	"cmd/main.go/internal/database"
	"cmd/main.go/internal/server"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

var (
	batchSize uint = 2
)

const grpsServerMainLogTag = "GrpsServerMain"

func main() {
	ctx := context.Background()

	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal(ctx, fmt.Sprintf("%s: failed init configuration", grpsServerMainLogTag),
			"err", err,
		)
	}
	cfg := config.GetConfigInstance()

	initCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := database.NewPostgres(initCtx, dsn, cfg.Database.Driver)
	if err != nil {
		return
	}
	defer db.Close()

	requestRepository := repo.NewUserRequestRepo(db, batchSize)

	userRequestService := user_request.New(db, requestRepository)

	if err := server.NewGrpcServer(userRequestService).Start(ctx, &cfg); err != nil {
		log.Print(ctx, fmt.Sprintf("%s: failed creating gRPC server", grpsServerMainLogTag), "err", err)

		return
	}

}
