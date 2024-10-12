package config

import "time"

var defaultConfig = Config{
	HTTP: HTTP{
		Bind:         "127.0.0.1:8080",
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	},
	PrimaryAddress: "",
	CanaryAddress:  "",
	SharedURLs:     []string{},
	Canary: Canary{
		Enabled:   false,
		Bucket:    Bucket{},
		Whitelist: []uint64{},
	},
	UserIDHeaderKey: "",
	UserNestedKey:   "",
}
