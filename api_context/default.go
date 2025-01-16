package api_context

type DefaultContext struct {
	roles []string
}

func (d *DefaultContext) SetAuthorizationClaims(map[string]interface{}) {

}

func (d *DefaultContext) SetApiKeyId(string) {

}

func (d *DefaultContext) SetAccessId(string) {

}

func (d *DefaultContext) Data(ApiContextData) {

}

func (d *DefaultContext) Salt() string {
	return ""
}

func (d *DefaultContext) Roles() []string {
	return d.roles
}

func (d *DefaultContext) SetRoles(roles []string) {
	if roles == nil {
		d.roles = []string{}
	}
	d.roles = append(d.roles, roles...)
}
