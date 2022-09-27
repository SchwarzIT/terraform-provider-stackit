package common

import "time"

const (
	// Common retry durations
	DURATION_1M  = 1 * time.Minute
	DURATION_5M  = 5 * time.Minute
	DURATION_10M = 10 * time.Minute
	DURATION_20M = 20 * time.Minute
	DURATION_30M = 30 * time.Minute
	DURATION_40M = 40 * time.Minute
	DURATION_1H  = 1 * time.Hour
	DURATION_2H  = 2 * time.Hour

	// Common client errors
	ERR_CLIENT_TIMEOUT = "Client.Timeout exceeded"
	ERR_NO_SUCH_HOST   = "dial tcp: lookup"

	// Acceptance tests constants
	ACC_TEST_PROJECT_ID = "2d07c991-6cd4-475e-b00b-8acc2494f73f"
)
