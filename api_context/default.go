package api_context

type DefaultContext struct {
	roles []string
}

func (d *DefaultContext) GetSalt() string {
	return ""
}

func (d *DefaultContext) GetRoles() []string {
	return d.roles
}

func (d *DefaultContext) SetRoles(roles []string) {
	if roles == nil {
		d.roles = []string{}
	}
	d.roles = append(d.roles, roles...)
}
