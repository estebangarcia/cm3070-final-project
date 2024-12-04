-- Modify "repositories" table
ALTER TABLE "repositories" ADD COLUMN "registry_repositories" bigint NULL, ADD CONSTRAINT "repositories_registries_repositories" FOREIGN KEY ("registry_repositories") REFERENCES "registries" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
