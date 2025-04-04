-- Create "organizations" table
CREATE TABLE "organizations" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "name" character varying NOT NULL, "slug" character varying NOT NULL, PRIMARY KEY ("id"));
-- Create index "organization_slug" to table: "organizations"
CREATE UNIQUE INDEX "organization_slug" ON "organizations" ("slug");
-- Create "organization_memberships" table
CREATE TABLE "organization_memberships" ("role" bigint NOT NULL, "user_id" bigint NOT NULL, "organization_id" bigint NOT NULL, PRIMARY KEY ("user_id", "organization_id"), CONSTRAINT "organization_memberships_organizations_organization" FOREIGN KEY ("organization_id") REFERENCES "organizations" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "organization_memberships_users_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create "registries" table
CREATE TABLE "registries" ("id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY, "name" character varying NOT NULL, "slug" character varying NOT NULL, "organization_registries" bigint NULL, PRIMARY KEY ("id"), CONSTRAINT "registries_organizations_registries" FOREIGN KEY ("organization_registries") REFERENCES "organizations" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
