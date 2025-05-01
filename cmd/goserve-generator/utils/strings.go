package utils

func Replacement(key string, value string) ReplaceEntry {
	return ReplaceEntry{
		Key:   key,
		Value: value,
	}
}
