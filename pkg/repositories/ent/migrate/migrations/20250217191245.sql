-- Modify "manifests" table
ALTER TABLE "manifests" ADD COLUMN "scanned_at" timestamptz NULL;
