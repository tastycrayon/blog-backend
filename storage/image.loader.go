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

// GetImage wraps the Image dataloader for efficient retrieval by image ID
func (i *DataLoader) GetImage(ctx context.Context, imageID int64) (*db.Image, error) {
	strId := strconv.FormatInt(imageID, 10)

	thunk := i.imageLoader.Load(ctx, dataloader.StringKey(strId))
	result, err := thunk()
	if err != nil {
		log.Fatalf("Failed at loader %v", err)
		return nil, err
	}
	image, ok := result.(db.Image)
	if ok {
		return &image, nil
	}

	log.Fatalf("Failed at loader %v", err)
	return nil, err

}

// imageBatcher wraps storage and provides a "get" method for the image dataloader
type imageBatcher struct {
	db db.Store
}

// get implements the dataloader for finding many images by Id and returns
// them in the order requested
func (img *imageBatcher) get(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	// collect the keys to search for
	var imageIDs []string
	for index, key := range keys {
		imageIDs = append(imageIDs, key.String())
		keyOrder[key.String()] = index
	}
	// search for those images
	dbRecords, err := img.db.Queries.GetImagesByIds(ctx, imageIDs)
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
	for imageID, ix := range keyOrder {
		err := fmt.Errorf("image not found %s", imageID)
		results[ix] = &dataloader.Result{Data: nil, Error: err}
	}
	return results
}
