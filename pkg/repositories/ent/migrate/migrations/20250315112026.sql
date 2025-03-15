-- Modify "organization_invites" table
ALTER TABLE "organization_invites" ADD COLUMN "invite_id" character varying NOT NULL;
-- Create index "organizationinvite_invite_id" to table: "organization_invites"
CREATE UNIQUE INDEX "organizationinvite_invite_id" ON "organization_invites" ("invite_id");
