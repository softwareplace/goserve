package model

import "time"

type ApiKeyEntryData struct {
	ClientName string        `json:"clientName"` // Client information for which the public key is generated (required)
	Expiration time.Duration `json:"expiration"` // Expiration specifies the duration until the API key expires (optional).
	ClientId   string        `json:"clientId"`   // ClientId represents the unique identifier for the client associated with the API key entry (required).
}

type Entry struct {
	Key        string
	Expiration time.Duration
	Roles      []string
	PublicKey  *string
}
