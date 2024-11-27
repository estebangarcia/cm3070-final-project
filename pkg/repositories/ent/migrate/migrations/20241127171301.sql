-- Drop index "blobchunk_session_id" from table: "blob_chunks"
DROP INDEX "blobchunk_session_id";
-- Drop index "blobchunk_upload_id" from table: "blob_chunks"
DROP INDEX "blobchunk_upload_id";
-- Create index "blobchunk_session_id" to table: "blob_chunks"
CREATE INDEX "blobchunk_session_id" ON "blob_chunks" ("session_id");
-- Create index "blobchunk_upload_id" to table: "blob_chunks"
CREATE INDEX "blobchunk_upload_id" ON "blob_chunks" ("upload_id");
-- Create index "blobchunk_upload_id_session_id_part_number" to table: "blob_chunks"
CREATE UNIQUE INDEX "blobchunk_upload_id_session_id_part_number" ON "blob_chunks" ("upload_id", "session_id", "part_number");
