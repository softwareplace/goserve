package api_context

type DefaultContext struct {
	roles []string
}

func NewDefaultCtx() *DefaultContext {
	return &DefaultContext{}
}

func (d *DefaultContext) GetSalt() string {
	return "081162586c7f4f77b877fbca0f09cb7f"
}

func (d *DefaultContext) GetRoles() []string {
	return d.roles
}

func (d *DefaultContext) SetRoles(roles ...string) {
	if roles == nil {
		d.roles = []string{}
	}
	d.roles = append(d.roles, roles...)
}
