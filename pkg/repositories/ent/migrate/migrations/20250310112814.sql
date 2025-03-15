-- Create index "organizationinvite_email_organization_id" to table: "organization_invites"
CREATE UNIQUE INDEX "organizationinvite_email_organization_id" ON "organization_invites" ("email", "organization_id");
-- Create index "organizationinvite_user_id_organization_id" to table: "organization_invites"
CREATE UNIQUE INDEX "organizationinvite_user_id_organization_id" ON "organization_invites" ("user_id", "organization_id");
