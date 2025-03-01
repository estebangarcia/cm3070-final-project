package repositories

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent"
	ent_manifest "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifest"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifesttagreference"
	ent_repository "github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/repository"
)

type ManifestTagRepository struct {
}

func NewManifestTagRepository() *ManifestTagRepository {
	return &ManifestTagRepository{}
}

func (mr *ManifestTagRepository) ListTagsForRepository(ctx context.Context, repository *ent.Repository, limit int, lastTagName string) ([]*ent.ManifestTagReference, error) {
	dbClient := getClient(ctx)

	predicate := manifesttagreference.HasManifestsWith(
		ent_manifest.HasRepositoryWith(
			ent_repository.ID(repository.ID),
		),
	)

	if lastTagName != "" {
		predicate = manifesttagreference.And(
			predicate,
			manifesttagreference.TagGT(lastTagName),
		)
	}

	return dbClient.ManifestTagReference.Query().Where(
		predicate,
	).Order(
		manifesttagreference.ByTag(
			sql.OrderAsc(),
		),
	).Limit(limit).All(ctx)
}

func (mr *ManifestTagRepository) GetTagByName(ctx context.Context, repository *ent.Repository, tagName string) (*ent.ManifestTagReference, bool, error) {
	dbClient := getClient(ctx)

	tag, err := dbClient.ManifestTagReference.Query().Where(
		manifesttagreference.And(
			manifesttagreference.HasManifestsWith(
				ent_manifest.HasRepositoryWith(
					ent_repository.ID(repository.ID),
				),
			),
			manifesttagreference.Tag(tagName),
		),
	).First(ctx)

	if err != nil && ent.IsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return tag, true, nil
}

func (mr *ManifestTagRepository) DeleteTag(ctx context.Context, tag *ent.ManifestTagReference) error {
	dbClient := getClient(ctx)
	return dbClient.ManifestTagReference.DeleteOne(tag).Exec(ctx)
}
