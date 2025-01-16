package api_context

type DefaultContext struct{}

func (d DefaultContext) SetAuthorizationClaims(map[string]interface{}) {

}

func (d DefaultContext) SetApiKeyId(string) {

}

func (d DefaultContext) SetAccessId(string) {

}

func (d DefaultContext) Data(ApiContextData) {

}

func (d DefaultContext) Salt() string {
	return ""
}

func (d DefaultContext) Roles() []string {
	return []string{}
}
