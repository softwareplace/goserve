package string

type ReplaceEntry struct {
	Key   string
	Value string
}

func Replacement(key string, value string) ReplaceEntry {
	return ReplaceEntry{
		Key:   key,
		Value: value,
	}
}
