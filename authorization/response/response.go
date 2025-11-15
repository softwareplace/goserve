package response

type AuthorizationResponse struct {
	Jwt      string `json:"jwt" validate:"required,gt=19"`
	Expires  int64  `json:"expires" validate:"required,gt=1757353373"`
	IssuedAt int64  `json:"issuedAt" validate:"required,gt=1757353373"`
}
