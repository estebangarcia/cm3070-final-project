package repositories

import (
	"context"
	"time"

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
	predicates := []predicate.Manifest{
		mr.getTagOrReferencePredicate(reference),
		ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
	}

	if len(mediaTypes) > 0 {
		predicates = append(predicates, ent_manifest.MediaTypeIn(mediaTypes...))
	}

	return mr.dbClient.Manifest.Query().Where(
		ent_manifest.And(
			predicates...,
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

func (mr *ManifestRepository) GetManifestByReference(ctx context.Context, reference string, repository *ent.Repository, withTags bool) (*ent.Manifest, bool, error) {
	query := mr.dbClient.Manifest.Query().Where(
		ent_manifest.And(
			mr.getTagOrReferencePredicate(reference),
			ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
		),
	)
	if withTags {
		query = query.WithTags()
	}

	manifest, err := query.First(ctx)

	if err != nil && ent.IsNotFound(err) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}

	return manifest, true, nil
}

func (mr *ManifestRepository) GetAllByTypeWithTags(ctx context.Context, artifactType string, repository *ent.Repository) ([]*ent.Manifest, error) {
	manifests, err := mr.dbClient.Manifest.Query().Where(
		ent_manifest.And(
			ent_manifest.Or(
				ent_manifest.ArtifactType(artifactType),
				ent_manifest.MediaType(artifactType),
			),
			ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
			ent_manifest.HasTags(),
		),
	).WithTags().WithManifestLayers().All(ctx)

	if err != nil {
		return nil, err
	}

	return manifests, nil
}

func (mr *ManifestRepository) GetAllWithTags(ctx context.Context, repository *ent.Repository) ([]*ent.Manifest, error) {
	manifests, err := mr.dbClient.Manifest.Query().Where(
		ent_manifest.And(
			ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
		),
	).WithTags().All(ctx)

	if err != nil {
		return nil, err
	}

	return manifests, nil
}

func (mr *ManifestRepository) CreateManifest(ctx context.Context, digest string, mediaType string, artifactType *string, s3Path string, subjectManifest *ent.Manifest, repository *ent.Repository) (*ent.Manifest, error) {
	client := mr.getClient(ctx)

	manifest := client.Manifest.
		Create().
		SetDigest(digest).
		SetMediaType(mediaType).
		SetS3Path(s3Path).
		SetNillableArtifactType(artifactType).
		SetRepository(repository)

	if subjectManifest != nil {
		manifest = manifest.AddSubject(subjectManifest)
	}

	return manifest.Save(ctx)
}

func (mr *ManifestRepository) CreateManifestLayers(ctx context.Context, layers []*ent.ManifestLayer, manifest *ent.Manifest) error {
	client := mr.getClient(ctx)

	return client.ManifestLayer.MapCreateBulk(layers, func(mlc *ent.ManifestLayerCreate, i int) {
		mlc.
			SetMediaType(layers[i].MediaType).
			SetDigest(layers[i].Digest).
			SetSize(layers[i].Size).
			SetAnnotations(layers[i].Annotations).
			SetManifestID(manifest.ID)
	}).OnConflictColumns("digest", "manifest_manifest_layers").UpdateNewValues().Exec(ctx)
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

func (mr *ManifestRepository) CreateVulnerabilitiesInBulkAndMarkAsScanned(ctx context.Context, vulnerabilities ent.Vulnerabilities, manifest *ent.Manifest) error {
	tx, err := mr.dbClient.Debug().Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctxV := context.WithValue(ctx, "tx", tx)

	for _, vulnerability := range vulnerabilities {
		err := tx.Vulnerability.Create().
			SetFixedVersion(vulnerability.FixedVersion).
			SetInstalledVersion(vulnerability.InstalledVersion).
			SetSeverity(vulnerability.Severity).
			SetStatus(vulnerability.Status).
			SetTitle(vulnerability.Title).
			SetV3Score(vulnerability.V3Score).
			SetVulnerabilityID(vulnerability.VulnerabilityID).
			SetVulnerabilityURLDetails(vulnerability.VulnerabilityURLDetails).
			AddManifests(manifest).
			OnConflictColumns("vulnerability_id").
			Ignore().Exec(ctxV)
		if err != nil {
			return err
		}
	}

	if err = mr.MarkAsScanned(ctxV, manifest); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (mr *ManifestRepository) UpsertManifestWithSubjectAndTag(ctx context.Context, layers []*ent.ManifestLayer, reference string, digest string, mediaType string, artifactType *string, s3Path string, manifestSubject *ent.Manifest, repository *ent.Repository) (*ent.Manifest, error) {
	tx, err := mr.dbClient.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ctxV := context.WithValue(ctx, "tx", tx)

	if manifestSubject != nil && manifestSubject.ID == 0 {
		manifestSubject, err = mr.CreateManifest(ctxV, manifestSubject.Digest, manifestSubject.MediaType, nil, manifestSubject.S3Path, nil, repository)
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
		mfst, err = mr.CreateManifest(ctxV, digest, mediaType, artifactType, s3Path, manifestSubject, repository)
	}

	if err != nil {
		return nil, err
	}

	err = mr.CreateManifestLayers(ctxV, layers, mfst)
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

func (mr *ManifestRepository) GetManifestReferrers(ctx context.Context, digest string, artifactType string, repository *ent.Repository) ([]*ent.Manifest, error) {
	manifestPredicate := []predicate.Manifest{
		manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
		manifest.HasRefererWith(
			manifest.Digest(digest),
		),
	}
	if artifactType != "" {
		manifestPredicate = append(manifestPredicate, manifest.ArtifactType(artifactType))
	}

	return mr.dbClient.Manifest.Query().Where(
		manifest.And(
			manifestPredicate...,
		),
	).All(ctx)
}

func (mr *ManifestRepository) GetAllUnscanned(ctx context.Context) ([]*ent.Manifest, error) {
	manifests, err := mr.dbClient.Manifest.Query().Where(
		ent_manifest.ScannedAtIsNil(),
	).WithRepository(
		func(rq *ent.RepositoryQuery) {
			rq.WithRegistry(
				func(rq *ent.RegistryQuery) {
					rq.WithOrganization()
				},
			)
		},
	).All(ctx)

	if err != nil {
		return nil, err
	}

	return manifests, nil
}

func (mr *ManifestRepository) MarkAsScanned(ctx context.Context, manifest *ent.Manifest) error {
	dbClient := mr.getClient(ctx)
	_, err := dbClient.Manifest.UpdateOne(manifest).SetScannedAt(time.Now()).Save(ctx)
	return err
}

func (mr *ManifestRepository) DeleteManifest(ctx context.Context, manifest *ent.Manifest) error {
	return mr.dbClient.Manifest.DeleteOne(manifest).Exec(ctx)
}

func (mr *ManifestRepository) getClient(ctx context.Context) *ent.Client {
	client := mr.dbClient
	if ctx.Value("tx") != nil {
		tx := ctx.Value("tx").(*ent.Tx)
		client = tx.Client()
	}
	return client
}
