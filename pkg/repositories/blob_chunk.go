package repositories

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/blobchunk"
)

type BlobChunkRepository struct {
}

func NewBlobChunkRepository() *BlobChunkRepository {
	return &BlobChunkRepository{}
}

// Create a blob chunk in the database
func (bcr *BlobChunkRepository) CreateBlobChunk(ctx context.Context, sessionId string, uploadId string, rangeFrom uint64, rangeTo uint64, partNumber uint64) (*ent.BlobChunk, error) {
	dbClient := getClient(ctx)
	return dbClient.BlobChunk.Create().
		SetSessionID(sessionId).
		SetUploadID(uploadId).
		SetRangeFrom(rangeFrom).
		SetRangeTo(rangeTo).
		SetPartNumber(partNumber).
		Save(ctx)
}

// Get the blob chunk count for an upload session
func (bcr *BlobChunkRepository) GetBlobChunkCount(ctx context.Context, sessionId string, uploadId string) (int, error) {
	dbClient := getClient(ctx)
	return dbClient.BlobChunk.Query().Where(
		blobchunk.And(
			blobchunk.SessionID(sessionId),
			blobchunk.UploadID(uploadId),
		),
	).Count(ctx)
}

// Get the latest uploaded blob chunk for an upload session
func (bcr *BlobChunkRepository) GetLatestBlobChunk(ctx context.Context, sessionId string, uploadId string) (*ent.BlobChunk, error) {
	dbClient := getClient(ctx)
	return dbClient.BlobChunk.Query().Where(
		blobchunk.And(
			blobchunk.SessionID(sessionId),
			blobchunk.UploadID(uploadId),
		),
	).Order(blobchunk.ByPartNumber(sql.OrderDesc())).First(ctx)
}

// Check if a blob chunk upload is out of order by checking the order of ranges from uploaded chunks
func (bcr *BlobChunkRepository) IsOutOfOrder(ctx context.Context, sessionId string, uploadId string, rangeFrom uint64, rangeTo uint64) (bool, error) {
	blob_chunk_count, err := bcr.GetBlobChunkCount(ctx, sessionId, uploadId)
	if err != nil {
		return false, err
	}

	if blob_chunk_count == 0 && rangeFrom == 0 {
		return false, nil
	} else if blob_chunk_count == 0 && rangeFrom > 0 {
		return true, nil
	}

	latestChunk, err := bcr.GetLatestBlobChunk(ctx, sessionId, uploadId)
	if err != nil {
		return false, err
	}

	return (latestChunk.RangeTo+1 != rangeFrom), nil
}

// Get the next blob chunk for an upload session
func (bcr *BlobChunkRepository) GetNext(ctx context.Context, sessionId string, uploadId string, rangeFrom uint64, rangeTo uint64) (*ent.BlobChunk, error) {
	blob_chunk_count, err := bcr.GetBlobChunkCount(ctx, sessionId, uploadId)
	if err != nil {
		return nil, err
	}

	if blob_chunk_count == 0 {
		return bcr.CreateBlobChunk(ctx, sessionId, uploadId, rangeFrom, rangeTo, 1)
	}

	latestChunk, err := bcr.GetLatestBlobChunk(ctx, sessionId, uploadId)
	if err != nil {
		return nil, err
	}

	return bcr.CreateBlobChunk(ctx, sessionId, uploadId, rangeFrom, rangeTo, latestChunk.PartNumber+1)
}

// Delete all blob chunks for an upload session
func (bcr *BlobChunkRepository) DeleteAllForUploadID(ctx context.Context, uploadId string) error {
	dbClient := getClient(ctx)
	_, err := dbClient.BlobChunk.Delete().Where(
		blobchunk.UploadID(uploadId),
	).Exec(ctx)
	return err
}

// Get all blob chunks for an upload session
func (bcr *BlobChunkRepository) GetByUploadID(ctx context.Context, uploadId string) ([]*ent.BlobChunk, error) {
	dbClient := getClient(ctx)
	return dbClient.BlobChunk.Query().Where(
		blobchunk.UploadID(uploadId),
	).All(ctx)
}
