-- Create "repositories" table
CREATE TABLE "repositories" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "name" character varying NOT NULL, PRIMARY KEY ("id"));
-- Create "manifests" table
CREATE TABLE "manifests" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "media_type" character varying NOT NULL, "s3_path" character varying NOT NULL, "digest" character varying NOT NULL, "repository_manifests" bigint NULL, PRIMARY KEY ("id"), CONSTRAINT "manifests_repositories_manifests" FOREIGN KEY ("repository_manifests") REFERENCES "repositories" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create "manifest_tag_references" table
CREATE TABLE "manifest_tag_references" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "tag" character varying NOT NULL, "manifest_tag_reference_manifests" bigint NULL, PRIMARY KEY ("id"), CONSTRAINT "manifest_tag_references_manifests_manifests" FOREIGN KEY ("manifest_tag_reference_manifests") REFERENCES "manifests" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
