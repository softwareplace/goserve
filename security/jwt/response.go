package jwt

type Response struct {
	JWT      string `json:"jwt"`
	Expires  int    `json:"expires"`
	IssuedAt int    `json:"issuedAt"`
}
