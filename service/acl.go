package service

import (
	"time"

	"bitbucket.org/pqstudio/go-acl/datastore"
	"bitbucket.org/pqstudio/go-acl/model"

	groupModel "bitbucket.org/pqstudio/go-user-group/model"
	groupS "bitbucket.org/pqstudio/go-user-group/service"

	"bitbucket.org/pqstudio/go-webutils"

	//. "bitbucket.org/pqstudio/go-webutils/logger"
)

func GetOne(uid string) (*model.ACL, error) {
	r, err := datastore.GetOne(uid)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func Check(userUID string, object string, permission string, action string) (bool, error) {
	gr, err := groupS.GetByUserUID(userUID, 100, 0)
	if err != nil {
		return false, err
	}

	allowed, err := CheckGroups(gr, object, permission, action)

	return allowed, err
}

func CheckGroup(userUID string, group string) (bool, error) {
	grs, err := groupS.GetByUserUID(userUID, 100, 0)
	if err != nil {
		return false, err
	}

	for _, gr := range grs {
		if group == gr.Name {
			return true, nil
		}
	}
	return false, nil
}

func CheckGroups(groups []groupModel.Group, object string, permission string, action string) (bool, error) {
	var e error
	for _, group := range groups {
		_, err := datastore.GetByGroupIDPermissionAndAction(group.Name, object, permission, action)

		if err != nil {
			e = err
			continue
		}

		return true, nil
	}

	return false, e
}

func CreateFromModel(m *model.ACL) error {
	m.UID = utils.NewUUID()
	m.CreatedAt = time.Now().UTC()

	err := datastore.Create(m)
	if err != nil {
		return err
	}

	return nil
}

func Create(userID string, objectID string, permission string, actions []string) error {

	var e error

	for _, action := range actions {
		acl := &model.ACL{
			GroupID:    userID,
			Object:     objectID,
			Permission: permission,
			Action:     action,
		}

		err := CreateFromModel(acl)
		if err != nil {
			e = err
			continue
		}
	}

	return e
}

func Update(m *model.ACL) error {
	err := datastore.Update(m)
	if err != nil {
		return err
	}

	return nil
}

func Delete(uid string) error {
	err := datastore.Delete(uid)
	if err != nil {
		return err
	}

	return nil
}
