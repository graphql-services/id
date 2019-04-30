package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/graphql-services/id"
	"github.com/graphql-services/id/database"
	"github.com/graphql-services/memberships"
)

const defaultPort = "80"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	urlString := os.Getenv("DATABASE_URL")
	if urlString == "" {
		panic(fmt.Errorf("database url must be provided"))
	}

	db := database.NewDBWithString(urlString)
	defer db.Close()
	db.AutoMigrate(&id.User{}, &id.UserActivationRequest{}, &id.ForgotPasswordRequest{})

	gqlHandler := handler.GraphQL(id.NewExecutableSchema(id.Config{Resolvers: &id.Resolver{DB: db}}))
	playgroundHandler := handler.Playground("GraphQL playground", "/graphql")
	http.HandleFunc("/graphql", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			playgroundHandler(res, req)
			return
		}
		ctx := context.WithValue(req.Context(), memberships.DBContextKey, db)
		req = req.WithContext(ctx)
		gqlHandler(res, req)
	})

	http.HandleFunc("/healthcheck", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("OK"))
	})

	log.Printf("connect to http://localhost:%s/graphql for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
