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
	dbClient *ent.Client
}

func NewManifestTagRepository(dbClient *ent.Client) *ManifestTagRepository {
	return &ManifestTagRepository{
		dbClient: dbClient,
	}
}

func (mr *ManifestTagRepository) ListTagsForRepository(ctx context.Context, repository *ent.Repository, limit int, lastTagName string) ([]*ent.ManifestTagReference, error) {
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

	return mr.dbClient.ManifestTagReference.Query().Where(
		predicate,
	).Order(
		manifesttagreference.ByTag(
			sql.OrderAsc(),
		),
	).Limit(limit).All(ctx)
}

func (mr *ManifestTagRepository) GetTagByName(ctx context.Context, repository *ent.Repository, tagName string) (*ent.ManifestTagReference, bool, error) {
	tag, err := mr.dbClient.ManifestTagReference.Query().Where(
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
