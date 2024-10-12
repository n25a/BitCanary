package config

var defaultConfig = Config{
	PrimaryAddress:   "127.0.0.1:8080",
	SecondaryAddress: "127.0.0.1:8081",
	SharedURLs:       []string{},
	Canary: Canary{
		Enabled:   false,
		Bucket:    Bucket{},
		Whitelist: []uint64{},
	},
	UserIDHeaderKey: "",
	UserNestedKey:   "",
}
