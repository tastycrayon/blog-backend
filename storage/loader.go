package storage

import (
	"context"
	"net/http"

	"github.com/graph-gophers/dataloader"
	"github.com/mosiur404/goserver/db"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// DataLoader offers data loaders scoped to a context
type DataLoader struct {
	userLoader  *dataloader.Loader
	imageLoader *dataloader.Loader
}

// NewDataLoader returns the instantiated Loaders struct for use in a request
func NewDataLoader(db db.Store) *DataLoader {
	// instantiate the dataloaders
	users := &userBatcher{db: db}
	images := &imageBatcher{db: db}
	// return the DataLoader
	return &DataLoader{
		userLoader:  dataloader.NewBatchedLoader(users.get),
		imageLoader: dataloader.NewBatchedLoader(images.get),
	}
}

// Middleware injects a DataLoader into the request context so it can be
// used later in the schema resolvers
func Middleware(loader *DataLoader, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCtx := context.WithValue(r.Context(), loadersKey, loader)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *DataLoader {
	return ctx.Value(loadersKey).(*DataLoader)
}
