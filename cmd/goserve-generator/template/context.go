package template

const Context = `package application

import goservectx "github.com/softwareplace/goserve/context"

type Ctx struct {
	roles             []string
	encryptedPassword string
	RequesterId       string
}

func New(requesterId string) goservectx.Principal {
	return &Ctx{
		RequesterId: requesterId,
	}
}

func (c *Ctx) GetId() string {
	return c.RequesterId
}

func (c *Ctx) GetRoles() []string {
	return c.roles
}

func (c *Ctx) SetEncryptedPassword(encryptedPassword string) {
	c.encryptedPassword = encryptedPassword
}

func (c *Ctx) EncryptedPassword() string {
	return c.encryptedPassword
}

func (c *Ctx) AddRoles(roles ...string) {
	if roles == nil {
		c.roles = []string{}
	}
	c.roles = append(c.roles, roles...)
}
`
