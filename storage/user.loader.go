package storage

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/graph-gophers/dataloader"
	"github.com/tastycrayon/blog-backend/db"
	"github.com/tastycrayon/blog-backend/util"
)

// GetUser wraps the User dataloader for efficient retrieval by user ID
func (i *DataLoader) GetUser(ctx context.Context, userID int64) (*db.User, error) {
	strId := strconv.FormatInt(userID, 10)

	thunk := i.userLoader.Load(ctx, dataloader.StringKey(strId))
	result, err := thunk()
	if err != nil {
		log.Fatalf("Failed at loader %v", err)
		return nil, err
	}
	user, ok := result.(db.User)
	if ok {
		return &user, nil
	}

	log.Fatalf("Failed at loader %v", err)
	return nil, err

}

// userBatcher wraps storage and provides a "get" method for the user dataloader
type userBatcher struct {
	db db.Store
}

// get implements the dataloader for finding many users by Id and returns
// them in the order requested
func (u *userBatcher) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// collect the keys to search for
	var userIDs []string
	for index, key := range keys {
		userIDs = append(userIDs, key.String())
		keyOrder[key.String()] = index
	}
	// search for those users
	dbRecords, err := u.db.Queries.GetUsersByIds(ctx, userIDs)
	// if DB error, return
	if err != nil {
		return []*dataloader.Result{{Data: nil, Error: err}}
	}
	// construct an output array of dataloader results
	results := make([]*dataloader.Result, len(keys))
	// // enumerate records, put into output
	for _, record := range dbRecords {
		index, ok := keyOrder[util.IDtoString(record.ID)]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[index] = &dataloader.Result{Data: record, Error: nil}
			delete(keyOrder, util.IDtoString(record.ID))
		}
	}
	// fill array positions with errors where not found in DB
	for userID, ix := range keyOrder {
		err := fmt.Errorf("user not found %s", userID)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}
	return results
}
