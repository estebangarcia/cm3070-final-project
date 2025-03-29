migrate:
	docker run --rm --net=host -v $(pwd)/pkg/repositories/ent/migrate/migrations:/migrations arigaio/atlas migrate apply --url "postgres://registry:registry@localhost:5432/registry?search_path=public&sslmode=disable"

migrate-status:
	docker run --rm --net=host -v $(pwd)/pkg/repositories/ent/migrate/migrations:/migrations arigaio/atlas migrate status --url "postgres://registry:registry@localhost:5432/registry?search_path=public&sslmode=disable"