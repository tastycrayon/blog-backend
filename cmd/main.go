package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mosiur404/goserver/db"
	"github.com/mosiur404/goserver/gql"
	"github.com/mosiur404/goserver/gql/generated"
	"github.com/mosiur404/goserver/middleware"
	"github.com/mosiur404/goserver/storage"
	"github.com/mosiur404/goserver/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		panic("❗could not load .env file")
	}

	dataBase := *db.NewStore(db.InitDB(config))

	// Start GraphQL
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(config.AllowedServer, "|"),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.AuthMiddleware(dataBase)) // access token
	router.Post("/refresh", gql.RefreshTokenRoute)  // refresh token

	c := generated.Config{Resolvers: &gql.Resolver{
		Db: dataBase,
	}}

	// make a data loader
	loader := storage.NewDataLoader(dataBase)
	// create the query handler
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	// whenever a panic happens, gracefully return a message to the user before stopping parsing
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		return util.GqlError(ctx, fmt.Sprintf("SetRecoverFunc: %v", err), util.InternalServerError)
	})
	// wrap the query handler with middleware to inject dataloader
	dataloaderSrv := storage.Middleware(loader, srv)

	router.Handle("/query", dataloaderSrv)
	// playground only available on developement mode
	if config.Environment != util.PRODUCTION {
		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	log.Printf("🚀 Server ready at: http://localhost:%s", config.HTTPServerPort)
	log.Fatal(http.ListenAndServe(":"+config.HTTPServerPort, router))
}
