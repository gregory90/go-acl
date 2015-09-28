package model

import (
	"github.com/jamieomatthews/validation"
	"github.com/martini-contrib/binding"
	"time"
)

type ACL struct {
	UID        string    `json:"uid"`
	GroupID    string    `json:"groupID"`
	Object     string    `json:"object"`
	Permission string    `json:"permission"`
	Action     string    `json:"action"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (acl ACL) Validate() binding.Errors {
	var errors binding.Errors

	v := validation.NewValidation(&errors, acl)
	v.KeyTag("json")

	//v.Validate(&acl.Email).Key("email").Message("required").Required()
	//v.Validate(&acl.Email).Message("incorrect").Email()

	//v.Validate(&acl.Password).Key("password").Message("required").Required()
	//v.Validate(&acl.Password).Message("range").Range(6, 60)

	return *v.Errors.(*binding.Errors)
}
