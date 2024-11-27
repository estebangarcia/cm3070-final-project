package repositories

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/blobchunk"
)

type BlobChunkRepository struct {
	dbClient *ent.Client
}

func NewBlobChunkRepository(dbClient *ent.Client) *BlobChunkRepository {
	return &BlobChunkRepository{
		dbClient: dbClient,
	}
}

func (bcr *BlobChunkRepository) GetForRange(ctx context.Context, sessionId string, uploadId string, rangeFrom uint64, rangeTo uint64) (*ent.BlobChunk, error) {
	return nil, nil
}

func (bcr *BlobChunkRepository) CreateBlobChunk(ctx context.Context, sessionId string, uploadId string, rangeFrom uint64, rangeTo uint64, partNumber uint64) (*ent.BlobChunk, error) {
	return bcr.dbClient.BlobChunk.Create().
		SetSessionID(sessionId).
		SetUploadID(uploadId).
		SetRangeFrom(rangeFrom).
		SetRangeTo(rangeTo).
		SetPartNumber(partNumber).
		Save(ctx)
}

func (bcr *BlobChunkRepository) GetBlobChunkCount(ctx context.Context, sessionId string, uploadId string) (int, error) {
	return bcr.dbClient.BlobChunk.Query().Where(
		blobchunk.And(
			blobchunk.SessionID(sessionId),
			blobchunk.UploadID(uploadId),
		),
	).Count(ctx)
}

func (bcr *BlobChunkRepository) GetLatestBlobChunk(ctx context.Context, sessionId string, uploadId string) (*ent.BlobChunk, error) {
	return bcr.dbClient.BlobChunk.Query().Where(
		blobchunk.And(
			blobchunk.SessionID(sessionId),
			blobchunk.UploadID(uploadId),
		),
	).Order(blobchunk.ByPartNumber(sql.OrderDesc())).First(ctx)
}

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

func (bcr *BlobChunkRepository) DeleteAllForUploadID(ctx context.Context, uploadId string) error {
	_, err := bcr.dbClient.BlobChunk.Delete().Where(
		blobchunk.UploadID(uploadId),
	).Exec(ctx)
	return err
}

func (bcr *BlobChunkRepository) GetByUploadID(ctx context.Context, uploadId string) ([]*ent.BlobChunk, error) {
	return bcr.dbClient.BlobChunk.Query().Where(
		blobchunk.UploadID(uploadId),
	).All(ctx)
}
