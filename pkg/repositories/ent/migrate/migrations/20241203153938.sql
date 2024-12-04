-- Create index "repository_name_registry_repositories" to table: "repositories"
CREATE UNIQUE INDEX "repository_name_registry_repositories" ON "repositories" ("name", "registry_repositories");
