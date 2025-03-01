-- Modify "manifest_layers" table
ALTER TABLE "manifest_layers" DROP CONSTRAINT "manifest_layers_manifests_manifest_layers", ADD CONSTRAINT "manifest_layers_manifests_manifest_layers" FOREIGN KEY ("manifest_manifest_layers") REFERENCES "manifests" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
