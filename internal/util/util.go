package util

func NullableString(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

func NullableInt64(val int64) *int64 {
	return &val
}

func NullableFloat64(val float64) *float64 {
	return &val
}
