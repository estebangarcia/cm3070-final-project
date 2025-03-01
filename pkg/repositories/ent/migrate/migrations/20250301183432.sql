-- Drop index "manifest_misconfigurations_manifest_id_key" from table: "manifest_misconfigurations"
DROP INDEX "manifest_misconfigurations_manifest_id_key";
-- Drop index "manifest_misconfigurations_misconfiguration_id_key" from table: "manifest_misconfigurations"
DROP INDEX "manifest_misconfigurations_misconfiguration_id_key";
-- Modify "manifest_misconfigurations" table
ALTER TABLE "manifest_misconfigurations" ALTER COLUMN "misconfiguration_id" DROP NOT NULL, ADD CONSTRAINT "manifest_misconfigurations_mis_94dc736889d0188c890b7a9691cce076" FOREIGN KEY ("misconfiguration_id") REFERENCES "misconfigurations" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
