package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	ent_manifest "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifest"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifesttagreference"
	ent_repository "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/repository"
)

type ManifestRepository struct {
	dbClient *ent.Client
}

func NewManifestRepository(dbClient *ent.Client) *ManifestRepository {
	return &ManifestRepository{
		dbClient: dbClient,
	}
}

func (mr *ManifestRepository) GetManifestsByReferenceAndMediaType(ctx context.Context, reference string, mediaTypes []string, repository *ent.Repository) ([]*ent.Manifest, error) {
	manifestPredicate := ent_manifest.HasTagsWith(manifesttagreference.Tag(reference))
	if helpers.IsSHA256Digest(reference) {
		manifestPredicate = ent_manifest.Digest(reference)
	}

	return mr.dbClient.Manifest.Query().Where(
		ent_manifest.And(
			manifestPredicate,
			ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
			ent_manifest.MediaTypeIn(mediaTypes...),
		),
	).All(ctx)
}

func (mr *ManifestRepository) CreateManifest(ctx context.Context, digest string, mediaType string, s3Path string, repository *ent.Repository) (*ent.Manifest, error) {
	client := mr.getClient(ctx)

	return client.Manifest.
		Create().
		SetDigest(digest).
		SetMediaType(mediaType).
		SetS3Path(s3Path).
		SetRepository(repository).
		Save(ctx)
}

func (mr *ManifestRepository) UpsertManifestTagReference(ctx context.Context, reference string, manifest *ent.Manifest, repository *ent.Repository) error {
	if helpers.IsSHA256Digest(reference) {
		return nil
	}

	client := mr.getClient(ctx)

	var tagReference *ent.ManifestTagReference

	tagReference, err := client.ManifestTagReference.Query().Where(
		manifesttagreference.And(
			manifesttagreference.HasManifestsWith(
				ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
			),
			manifesttagreference.Tag(reference),
		),
	).First(ctx)

	tagReferenceNotFound := (err != nil && ent.IsNotFound(err))

	if tagReferenceNotFound {
		tagReference, err = client.ManifestTagReference.Create().
			SetManifests(manifest).
			SetTag(reference).
			Save(ctx)
	}

	if err != nil {
		return err
	}

	_, err = tagReference.Update().SetManifestsID(manifest.ID).Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (mr *ManifestRepository) CreateManifestAndUpsertTag(ctx context.Context, reference string, digest string, mediaType string, s3Path string, repository *ent.Repository) (*ent.Manifest, error) {
	tx, err := mr.dbClient.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ctxV := context.WithValue(ctx, "tx", tx)

	mfst, err := mr.CreateManifest(ctxV, digest, mediaType, s3Path, repository)
	if err != nil {
		return nil, err
	}

	err = mr.UpsertManifestTagReference(ctxV, reference, mfst, repository)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return mfst, nil
}

func (mr *ManifestRepository) getClient(ctx context.Context) *ent.Client {
	client := mr.dbClient
	if ctx.Value("tx") != nil {
		tx := ctx.Value("tx").(*ent.Tx)
		client = tx.Client()
	}
	return client
}
