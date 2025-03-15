// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/manifest"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/organizationinvite"
	"github.com/estebangarcia/cm3070-final-project/pkg/repositories/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	manifestFields := schema.Manifest{}.Fields()
	_ = manifestFields
	// manifestDescUploadedAt is the schema descriptor for uploaded_at field.
	manifestDescUploadedAt := manifestFields[5].Descriptor()
	// manifest.DefaultUploadedAt holds the default value on creation for the uploaded_at field.
	manifest.DefaultUploadedAt = manifestDescUploadedAt.Default.(func() time.Time)
	organizationinviteFields := schema.OrganizationInvite{}.Fields()
	_ = organizationinviteFields
	// organizationinviteDescInviteID is the schema descriptor for invite_id field.
	organizationinviteDescInviteID := organizationinviteFields[0].Descriptor()
	// organizationinvite.DefaultInviteID holds the default value on creation for the invite_id field.
	organizationinvite.DefaultInviteID = organizationinviteDescInviteID.Default.(func() string)
}
