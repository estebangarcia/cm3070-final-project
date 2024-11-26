package helpers

import "strings"

const sha256Prefix = "sha256:"

func IsSHA256Digest(digest string) bool {
	return strings.HasPrefix(digest, sha256Prefix)
}

func GetDigestAsNestedFolder(digest string) string {
	// Remove the "sha256:" prefix if it exists
	digest = strings.TrimPrefix(digest, sha256Prefix)

	// Split the digest into chunks of 2 characters
	var folders []string
	for i := 0; i < len(digest); i += 2 {
		if i+2 > len(digest) {
			// Handle any remainder (unlikely with SHA256 as it's 64 characters)
			folders = append(folders, digest[i:])
		} else {
			folders = append(folders, digest[i:i+2])
		}
	}

	// Join the folders with a "/" to simulate the S3 folder structure
	return strings.Join(folders, "/")
}

func IsVendorSpecificContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "application/vnd")
}
