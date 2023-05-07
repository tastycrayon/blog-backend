package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tastycrayon/blog-backend/db"
	"github.com/tastycrayon/blog-backend/gql"
	"github.com/tastycrayon/blog-backend/gql/generated"
	"github.com/tastycrayon/blog-backend/middleware"
	"github.com/tastycrayon/blog-backend/storage"
	"github.com/tastycrayon/blog-backend/util"
)

func Service(config util.Config) (http.Handler, error) {
	router := chi.NewRouter()
	conn, err := db.InitDB(config)
	if err != nil {
		return router, err
	}
	var dataBase db.Store = *db.NewStore(conn)
	// Start GraphQL
	OriginFunc := func(r *http.Request, origin string) bool {
		return origin == r.Header.Get("Origin")
		// return true
	}
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   strings.Split(config.AllowedServer, "|"),
		AllowOriginFunc:  OriginFunc,
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
	} else {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("HUH ¯\\_(ツ)_/¯"))
		})
	}
	return router, nil
}
