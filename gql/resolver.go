package gql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import "github.com/mosiur404/goserver/db"

type Resolver struct {
	Db db.Store
}
