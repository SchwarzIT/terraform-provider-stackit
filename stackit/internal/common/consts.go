package common

const (
	// Acceptance tests constants
	// https://portal.stackit.cloud/projects/8a2d2862-ac85-4084-8144-4c72d92ddcdd/dashboard
	ACC_TEST_PROJECT_ID = "8a2d2862-ac85-4084-8144-4c72d92ddcdd"

	// errors
	ERR_UNEXPECTED_EOF = "unexpected EOF"
)

// KnownRanges are the known ranges of IP addresses used by STACKIT
var KnownRanges = []string{
	"193.148.160.0/19",
	"45.129.40.0/21",
	"45.135.244.0/22",
}
