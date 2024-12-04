-- Create index "registry_slug_organization_registries" to table: "registries"
CREATE UNIQUE INDEX "registry_slug_organization_registries" ON "registries" ("slug", "organization_registries");
