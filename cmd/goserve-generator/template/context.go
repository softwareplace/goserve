package template

const Context = `package application

import goservectx "github.com/softwareplace/goserve/context"

type Principal struct {
	roles             []string
	encryptedPassword string
	RequesterId       string
}

func New(requesterId string) goservectx.Principal {
	return &Principal{
		RequesterId: requesterId,
	}
}

func (c *Principal) GetId() string {
	return c.RequesterId
}

func (c *Principal) GetRoles() []string {
	return c.roles
}

func (c *Principal) SetEncryptedPassword(encryptedPassword string) {
	c.encryptedPassword = encryptedPassword
}

func (c *Principal) EncryptedPassword() string {
	return c.encryptedPassword
}

func (c *Principal) AddRoles(roles ...string) {
	if roles == nil {
		c.roles = []string{}
	}
	c.roles = append(c.roles, roles...)
}
`
