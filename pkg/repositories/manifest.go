package repositories

import (
	"context"

	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifest"
	ent_manifest "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifest"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifesttagreference"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
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

func (mr *ManifestRepository) getTagOrReferencePredicate(reference string) predicate.Manifest {
	manifestPredicate := ent_manifest.HasTagsWith(manifesttagreference.Tag(reference))
	if helpers.IsSHA256Digest(reference) {
		manifestPredicate = ent_manifest.Digest(reference)
	}

	return manifestPredicate
}

func (mr *ManifestRepository) GetManifestsByReferenceAndMediaType(ctx context.Context, reference string, mediaTypes []string, repository *ent.Repository) ([]*ent.Manifest, error) {
	return mr.dbClient.Manifest.Query().Where(
		ent_manifest.And(
			mr.getTagOrReferencePredicate(reference),
			ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
			ent_manifest.MediaTypeIn(mediaTypes...),
		),
	).All(ctx)
}

func (mr *ManifestRepository) GetManifestByReferenceAndMediaType(ctx context.Context, reference string, mediaType string, repository *ent.Repository) (*ent.Manifest, bool, error) {
	manifest, err := mr.dbClient.Manifest.Query().Where(
		ent_manifest.And(
			mr.getTagOrReferencePredicate(reference),
			ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
			ent_manifest.MediaType(mediaType),
		),
	).First(ctx)

	if err != nil && ent.IsNotFound(err) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return manifest, true, nil
}

func (mr *ManifestRepository) CreateManifest(ctx context.Context, digest string, mediaType string, s3Path string, subjectManifest *ent.Manifest, repository *ent.Repository) (*ent.Manifest, error) {
	client := mr.getClient(ctx)

	manifest := client.Manifest.
		Create().
		SetDigest(digest).
		SetMediaType(mediaType).
		SetS3Path(s3Path).
		SetRepository(repository)

	if subjectManifest != nil {
		manifest = manifest.AddSubject(subjectManifest)
	}

	return manifest.Save(ctx)
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

	tagReferenceUpdate := client.ManifestTagReference.UpdateOneID(tagReference.ID)
	_, err = tagReferenceUpdate.SetManifestsID(manifest.ID).Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (mr *ManifestRepository) UpsertManifestWithSubjectAndTag(ctx context.Context, reference string, digest string, mediaType string, s3Path string, manifestSubject *ent.Manifest, repository *ent.Repository) (*ent.Manifest, error) {
	tx, err := mr.dbClient.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ctxV := context.WithValue(ctx, "tx", tx)

	if manifestSubject != nil && manifestSubject.ID == 0 {
		manifestSubject, err = mr.CreateManifest(ctxV, manifestSubject.Digest, manifestSubject.MediaType, manifestSubject.S3Path, nil, repository)
		if err != nil {
			return nil, err
		}
	}

	mfst, found, err := mr.GetManifestByReferenceAndMediaType(ctxV, digest, mediaType, repository)
	if err != nil {
		return nil, err
	}

	if found {
		manifestUpdate := tx.Manifest.UpdateOne(mfst)
		manifestUpdate = manifestUpdate.SetDigest(digest).SetMediaType(mediaType).SetS3Path(s3Path)
		if manifestSubject != nil {
			manifestUpdate = manifestUpdate.AddSubject(manifestSubject)
		}
		mfst, err = manifestUpdate.Save(ctxV)
	} else {
		mfst, err = mr.CreateManifest(ctxV, digest, mediaType, s3Path, manifestSubject, repository)
	}

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

func (mr *ManifestRepository) GetManifestReferrers(ctx context.Context, digest string, repository *ent.Repository) ([]*ent.Manifest, error) {
	return mr.dbClient.Manifest.Query().Where(
		manifest.And(
			manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
			manifest.HasRefererWith(
				manifest.Digest(digest),
			),
		),
	).All(ctx)
}

func (mr *ManifestRepository) getClient(ctx context.Context) *ent.Client {
	client := mr.dbClient
	if ctx.Value("tx") != nil {
		tx := ctx.Value("tx").(*ent.Tx)
		client = tx.Client()
	}
	return client
}
