package repositories

import (
	"context"
	"slices"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifest"
	ent_manifest "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifest"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifestlayer"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifestmisconfiguration"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifesttagreference"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/misconfiguration"
	ent_organization "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organization"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/predicate"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/registry"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/repository"
	ent_repository "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/repository"
	ent_vulnerability "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/vulnerability"
)

type ManifestRepository struct{}

func NewManifestRepository() *ManifestRepository {
	return &ManifestRepository{}
}

// Build an entgo query predicte to query either by digest or tag depending on the passed in reference
func (mr *ManifestRepository) getTagOrReferencePredicate(reference string) predicate.Manifest {
	manifestPredicate := ent_manifest.HasTagsWith(manifesttagreference.Tag(reference))
	if helpers.IsSHA256Digest(reference) {
		manifestPredicate = ent_manifest.Digest(reference)
	}

	return manifestPredicate
}

// Get all manifests by reference and media types that belong to the specified repository
func (mr *ManifestRepository) GetManifestsByReferenceAndMediaType(ctx context.Context, reference string, mediaTypes []string, repository *ent.Repository) ([]*ent.Manifest, error) {
	dbClient := getClient(ctx)
	predicates := []predicate.Manifest{
		mr.getTagOrReferencePredicate(reference),
		ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
	}

	if len(mediaTypes) > 0 && !slices.Contains(mediaTypes, "*/*") {
		predicates = append(predicates, ent_manifest.MediaTypeIn(mediaTypes...))
	}

	return dbClient.Manifest.Query().Where(
		ent_manifest.And(
			predicates...,
		),
	).All(ctx)
}

// Get a manifest by reference and media type that belong to the specified repository
func (mr *ManifestRepository) GetManifestByReferenceAndMediaType(ctx context.Context, reference string, mediaType string, repository *ent.Repository) (*ent.Manifest, bool, error) {
	dbClient := getClient(ctx)

	manifest, err := dbClient.Manifest.Query().Where(
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

// Get a specific manifest by reference that belong to the specified repository
func (mr *ManifestRepository) GetManifestByReference(ctx context.Context, reference string, repository *ent.Repository, withTags bool) (*ent.Manifest, bool, error) {
	dbClient := getClient(ctx)

	query := dbClient.Manifest.Query().Where(
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

// Get a manifest vulnerabilities by its reference that balong to the specified repository
func (mr *ManifestRepository) GetManifestVulnerabilitiesByReference(ctx context.Context, reference string, repository *ent.Repository) (ent.Vulnerabilities, error) {
	dbClient := getClient(ctx)

	return dbClient.Vulnerability.Query().Where(
		ent_vulnerability.HasManifestsWith(
			ent_manifest.And(
				mr.getTagOrReferencePredicate(reference),
				ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
			),
		),
	).All(ctx)
}

// Get a manifest misconfigurations by its reference that balong to the specified repository
func (mr *ManifestRepository) GetManifestMisconfigurationsByReference(ctx context.Context, reference string, repository *ent.Repository) (ent.ManifestMisconfigurations, error) {
	dbClient := getClient(ctx)

	manifestId, err := dbClient.Manifest.Query().Where(
		ent_manifest.And(
			mr.getTagOrReferencePredicate(reference),
			ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
		),
	).FirstID(ctx)
	if err != nil {
		return nil, err
	}

	return dbClient.ManifestMisconfiguration.Query().Where(
		manifestmisconfiguration.ManifestID(manifestId),
	).WithMisconfiguration().All(ctx)
}

// Get all manifests by type with tags that belong to the specified repository
func (mr *ManifestRepository) GetAllByTypeWithTags(ctx context.Context, artifactType string, repository *ent.Repository) ([]*ent.Manifest, error) {
	dbClient := getClient(ctx)

	manifests, err := dbClient.Manifest.Query().Where(
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

// Get all manifests with tags in the specified repository
func (mr *ManifestRepository) GetAllWithTags(ctx context.Context, repository *ent.Repository) ([]*ent.Manifest, error) {
	dbClient := getClient(ctx)

	manifests, err := dbClient.Manifest.Query().Where(
		ent_manifest.And(
			ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
		),
	).WithTags().WithRepository(func(rq *ent.RepositoryQuery) {
		rq.WithRegistry()
	},
	).All(ctx)

	if err != nil {
		return nil, err
	}

	return manifests, nil
}

// Get all manifests with tags in the specified organization
func (mr *ManifestRepository) GetAllForOrgWithTags(ctx context.Context, organization *ent.Organization) ([]*ent.Manifest, error) {
	dbClient := getClient(ctx)

	manifests, err := dbClient.Manifest.Query().Where(
		ent_manifest.And(
			ent_manifest.HasRepositoryWith(
				ent_repository.HasRegistryWith(
					registry.HasOrganizationWith(
						ent_organization.ID(organization.ID),
					),
				),
			),
		),
	).WithTags().WithRepository(func(rq *ent.RepositoryQuery) {
		rq.WithRegistry(
			func(rq *ent.RegistryQuery) {
				rq.WithOrganization()
			},
		)
	}).Order(ent_manifest.ByUploadedAt(sql.OrderDesc())).All(ctx)

	if err != nil {
		return nil, err
	}

	return manifests, nil
}

// Create a manifest in the database
func (mr *ManifestRepository) CreateManifest(ctx context.Context, digest string, mediaType string, artifactType *string, s3Path string, subjectManifest *ent.Manifest, repository *ent.Repository) (*ent.Manifest, error) {
	dbClient := getClient(ctx)

	manifest := dbClient.Manifest.
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

// Create a manifest's layers in bulk
func (mr *ManifestRepository) CreateManifestLayers(ctx context.Context, layers []*ent.ManifestLayer, manifest *ent.Manifest) error {
	dbClient := getClient(ctx)

	return dbClient.ManifestLayer.MapCreateBulk(layers, func(mlc *ent.ManifestLayerCreate, i int) {
		mlc.
			SetMediaType(layers[i].MediaType).
			SetDigest(layers[i].Digest).
			SetSize(layers[i].Size).
			SetAnnotations(layers[i].Annotations).
			SetManifestID(manifest.ID)
	}).OnConflictColumns("digest", "manifest_manifest_layers").UpdateNewValues().Exec(ctx)
}

// Upsert a manifest by tag or digest, this check if the manifest already exists and only creates the tag reference
func (mr *ManifestRepository) UpsertManifestTagReference(ctx context.Context, reference string, manifest *ent.Manifest, repository *ent.Repository) error {
	if helpers.IsSHA256Digest(reference) {
		return nil
	}

	dbClient := getClient(ctx)

	var tagReference *ent.ManifestTagReference

	tagReference, err := dbClient.ManifestTagReference.Query().Where(
		manifesttagreference.And(
			manifesttagreference.HasManifestsWith(
				ent_manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
			),
			manifesttagreference.Tag(reference),
		),
	).First(ctx)

	tagReferenceNotFound := (err != nil && ent.IsNotFound(err))

	if tagReferenceNotFound {
		tagReference, err = dbClient.ManifestTagReference.Create().
			SetManifests(manifest).
			SetTag(reference).
			Save(ctx)
	}

	if err != nil {
		return err
	}

	tagReferenceUpdate := dbClient.ManifestTagReference.UpdateOneID(tagReference.ID)
	_, err = tagReferenceUpdate.SetManifestsID(manifest.ID).Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Create vulnerabilities and relationship to manifest in bulk
func (mr *ManifestRepository) CreateVulnerabilitiesInBulkAndMarkAsScanned(ctx context.Context, vulnerabilities ent.Vulnerabilities, manifest *ent.Manifest) error {
	dbClient := getClient(ctx)

	for _, vulnerability := range vulnerabilities {
		err := dbClient.Vulnerability.Create().
			SetPackageName(vulnerability.PackageName).
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
			Ignore().Exec(ctx)
		if err != nil {
			return err
		}
	}

	if err := mr.MarkAsScanned(ctx, manifest); err != nil {
		return err
	}

	return nil
}

// Create misconfigurations in bulk
func (mr *ManifestRepository) CreateMisconfigurationsInBulk(ctx context.Context, misconfigurations ent.Misconfigurations) error {
	dbClient := getClient(ctx)
	return dbClient.Misconfiguration.MapCreateBulk(misconfigurations, func(mc *ent.MisconfigurationCreate, i int) {
		mc.SetMisconfigurationID(misconfigurations[i].MisconfigurationID).
			SetTitle(misconfigurations[i].Title).
			SetMisconfigurationURLDetails(misconfigurations[i].MisconfigurationURLDetails).
			SetSeverity(misconfigurations[i].Severity)
	}).OnConflictColumns("misconfiguration_id").DoNothing().Exec(ctx)
}

// Create misconfigurations and relationship to manifest in bulk
func (mr *ManifestRepository) CreateManifestMisconfigurationsInBulk(ctx context.Context, manifestMisconfigurations ent.ManifestMisconfigurations) error {
	dbClient := getClient(ctx)
	return dbClient.ManifestMisconfiguration.MapCreateBulk(manifestMisconfigurations, func(mmc *ent.ManifestMisconfigurationCreate, i int) {
		mmc.SetTargetFile(manifestMisconfigurations[i].TargetFile).
			SetMessage(manifestMisconfigurations[i].Message).
			SetResolution(manifestMisconfigurations[i].Resolution).
			SetManifestID(manifestMisconfigurations[i].ManifestID).
			SetMisconfigurationID(manifestMisconfigurations[i].MisconfigurationID)
	}).Exec(ctx)
}

// Get a list of misconfigurations by ID
func (mr *ManifestRepository) GetMisconfigurationsByIDs(ctx context.Context, ids []string) (ent.Misconfigurations, error) {
	dbClient := getClient(ctx)
	return dbClient.Misconfiguration.Query().Where(
		misconfiguration.MisconfigurationIDIn(ids...),
	).All(ctx)
}

// Upsert a manifest by tag or digest, this check if the manifest already exists and only creates the tag reference
// it also creates the manifest being sent as a subject
func (mr *ManifestRepository) UpsertManifestWithSubjectAndTag(ctx context.Context, layers []*ent.ManifestLayer, reference string, digest string, mediaType string, artifactType *string, s3Path string, manifestSubject *ent.Manifest, repository *ent.Repository) (*ent.Manifest, error) {
	var err error
	dbClient := getClient(ctx)

	if manifestSubject != nil && manifestSubject.ID == 0 {
		manifestSubject, err = mr.CreateManifest(ctx, manifestSubject.Digest, manifestSubject.MediaType, nil, manifestSubject.S3Path, nil, repository)
		if err != nil {
			return nil, err
		}
	}

	mfst, found, err := mr.GetManifestByReferenceAndMediaType(ctx, digest, mediaType, repository)
	if err != nil {
		return nil, err
	}

	if found {
		manifestUpdate := dbClient.Manifest.UpdateOne(mfst)
		manifestUpdate = manifestUpdate.SetDigest(digest).SetMediaType(mediaType).SetS3Path(s3Path)
		if manifestSubject != nil {
			manifestUpdate = manifestUpdate.AddSubject(manifestSubject)
		}
		mfst, err = manifestUpdate.Save(ctx)
	} else {
		mfst, err = mr.CreateManifest(ctx, digest, mediaType, artifactType, s3Path, manifestSubject, repository)
	}

	if err != nil {
		return nil, err
	}

	err = mr.CreateManifestLayers(ctx, layers, mfst)
	if err != nil {
		return nil, err
	}

	err = mr.UpsertManifestTagReference(ctx, reference, mfst, repository)
	if err != nil {
		return nil, err
	}

	return mfst, nil
}

// Get all of a manifest's referrers
func (mr *ManifestRepository) GetManifestReferrers(ctx context.Context, digest string, artifactType string, repository *ent.Repository) ([]*ent.Manifest, error) {
	dbClient := getClient(ctx)

	manifestPredicate := []predicate.Manifest{
		manifest.HasRepositoryWith(ent_repository.ID(repository.ID)),
		manifest.HasRefererWith(
			manifest.Digest(digest),
		),
	}
	if artifactType != "" {
		manifestPredicate = append(manifestPredicate, manifest.ArtifactType(artifactType))
	}

	return dbClient.Manifest.Query().Where(
		manifest.And(
			manifestPredicate...,
		),
	).All(ctx)
}

// Get all unscanned manifests
func (mr *ManifestRepository) GetAllUnscanned(ctx context.Context) ([]*ent.Manifest, error) {
	dbClient := getClient(ctx)

	manifests, err := dbClient.Manifest.Query().Where(
		ent_manifest.ScannedAtIsNil(),
	).WithRepository(
		func(rq *ent.RepositoryQuery) {
			rq.WithRegistry(
				func(rq *ent.RegistryQuery) {
					rq.WithOrganization()
				},
			)
		},
	).WithManifestLayers().All(ctx)

	if err != nil {
		return nil, err
	}

	return manifests, nil
}

// Get all scanned manifests
func (mr *ManifestRepository) MarkAsScanned(ctx context.Context, manifest *ent.Manifest) error {
	dbClient := getClient(ctx)
	_, err := dbClient.Manifest.UpdateOne(manifest).SetScannedAt(time.Now()).Save(ctx)
	return err
}

// Delete the specified manifest
func (mr *ManifestRepository) DeleteManifest(ctx context.Context, manifest *ent.Manifest) error {
	dbClient := getClient(ctx)
	return dbClient.Manifest.DeleteOne(manifest).Exec(ctx)
}

// Get all the unique layers for a manifest. This means layers that aren't being shared by other manifests
func (mr *ManifestRepository) GetUniqueManifestLayers(ctx context.Context, manifest *ent.Manifest) (ent.ManifestLayers, error) {
	dbClient := getClient(ctx)

	return dbClient.ManifestLayer.
		Query().
		Where(manifestlayer.HasManifestWith(
			ent_manifest.ID(manifest.ID),
		)).
		Where(func(s *sql.Selector) {
			// Column reference for the digest in the outer query.
			digestCol := s.C(manifestlayer.FieldDigest)
			// Build a subquery to look for any rows with the same digest but a different manifest_id.
			subq := sql.
				SelectExpr(sql.Expr("1")).
				From(sql.Table(manifestlayer.Table)).
				Where(sql.EQ(manifestlayer.FieldDigest, digestCol)).
				Where(sql.NEQ(manifestlayer.ManifestColumn, manifest.ID))
			// Only include rows for which the subquery returns no rows.
			s.Where(sql.NotExists(subq))
		}).All(ctx)
}

// Get the sum of storage used by manifests for a specific organization in bytes
func (mr *ManifestRepository) GetStorageUsedInBytesForOrganization(ctx context.Context, organization *ent.Organization) (int, error) {
	dbClient := getClient(ctx)

	return dbClient.ManifestLayer.Query().Aggregate(
		func(s *sql.Selector) string {
			manifestLayerTable := sql.Table(manifestlayer.Table)
			manifestTable := sql.Table(manifest.Table)
			repositoryTable := sql.Table(repository.Table)
			registryTable := sql.Table(registry.Table)
			organizationTable := sql.Table(ent_organization.Table)

			aggregationLayerSize := sql.Select(sql.As(sql.Max(manifestLayerTable.C(manifestlayer.FieldSize)), "size_per_layer")).
				From(manifestLayerTable).
				Join(manifestTable).On(manifestLayerTable.C(manifestlayer.ManifestColumn), manifestTable.C(manifest.FieldID)).
				Join(repositoryTable).On(manifestTable.C(manifest.RepositoryColumn), repositoryTable.C(repository.FieldID)).
				Join(registryTable).On(repositoryTable.C(repository.RegistryColumn), registryTable.C(registry.FieldID)).
				Join(organizationTable).On(registryTable.C(registry.OrganizationColumn), organizationTable.C(ent_organization.FieldID)).
				Where(sql.EQ(organizationTable.C(ent_organization.FieldID), organization.ID)).
				GroupBy(manifestLayerTable.C(manifestlayer.FieldDigest)).As("aggregation_layer_size")

			s.From(aggregationLayerSize)
			q, _ := sql.ExprFunc(func(b *sql.Builder) {
				b.WriteString("COALESCE(").Ident(sql.Sum(aggregationLayerSize.C("size_per_layer")))
				b.Comma().WriteString("0)")
			}).Query()
			return sql.As(q, "total_storage_used")
		},
	).Int(ctx)
}

// Get the manifest count for an organization
func (mr *ManifestRepository) GetCountForOrg(ctx context.Context, organization *ent.Organization) (int, error) {
	dbClient := getClient(ctx)

	return dbClient.Manifest.Query().Where(
		ent_manifest.And(
			ent_manifest.HasRepositoryWith(
				ent_repository.HasRegistryWith(
					registry.HasOrganizationWith(
						ent_organization.ID(organization.ID),
					),
				),
			),
		),
	).Count(ctx)
}

// Get the count of manifests with vulnerabilities for an organization
func (mr *ManifestRepository) GetCountWithVulnerabilitiesForOrg(ctx context.Context, organization *ent.Organization) (int, error) {
	dbClient := getClient(ctx)

	return dbClient.Manifest.Query().Where(
		ent_manifest.And(
			ent_manifest.HasRepositoryWith(
				ent_repository.HasRegistryWith(
					registry.HasOrganizationWith(
						ent_organization.ID(organization.ID),
					),
				),
			),
			ent_manifest.HasVulnerabilities(),
		),
	).Count(ctx)
}
