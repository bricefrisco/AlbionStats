package util

func NullableString(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}
