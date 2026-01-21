package util

func NullableString(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

func IsValidServer(server string) bool {
	validServers := map[string]bool{
		"americas": true,
		"europe":   true,
		"asia":     true,
	}
	return validServers[server]
}
