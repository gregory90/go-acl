package service

import (
	"database/sql"
	"time"

	"bitbucket.org/pqstudio/go-acl/datastore"
	"bitbucket.org/pqstudio/go-acl/model"

	groupModel "bitbucket.org/pqstudio/go-user-group/model"
	groupS "bitbucket.org/pqstudio/go-user-group/service"

	"bitbucket.org/pqstudio/go-webutils"
)

func GetOne(tx *sql.Tx, uid string) (*model.ACL, error) {
	r, err := datastore.GetOne(tx, uid)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func Check(tx *sql.Tx, userUID string, object string, permission string, action string) (bool, error) {
	grs, err := groupS.GetByUserUID(tx, userUID, 100, 0)
	if err != nil {
		return false, err
	}

	allowed, err := CheckGroups(tx, grs, object, permission, action)

	return allowed, err
}

func CheckGroup(tx *sql.Tx, userUID string, group string) (bool, error) {
	grs, err := groupS.GetByUserUID(tx, userUID, 100, 0)
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

func CheckGroups(tx *sql.Tx, groups []groupModel.Group, object string, permission string, action string) (bool, error) {
	var e error
	for _, group := range groups {
		_, err := datastore.GetByGroupIDPermissionAndAction(tx, group.Name, object, permission, action)

		if err != nil {
			e = err
			continue
		}

		return true, nil
	}

	return false, e
}

func CreateFromModel(tx *sql.Tx, m *model.ACL) error {
	m.UID = utils.NewUUID()
	m.CreatedAt = time.Now().UTC()

	err := datastore.Create(tx, m)
	if err != nil {
		return err
	}

	return nil
}

func Create(tx *sql.Tx, userID string, objectID string, permission string, actions []string) error {

	var e error

	for _, action := range actions {
		acl := &model.ACL{
			GroupID:    userID,
			Object:     objectID,
			Permission: permission,
			Action:     action,
		}

		err := CreateFromModel(tx, acl)
		if err != nil {
			e = err
			continue
		}
	}

	return e
}

func Update(tx *sql.Tx, m *model.ACL) error {
	err := datastore.Update(tx, m)
	if err != nil {
		return err
	}

	return nil
}

func Delete(tx *sql.Tx, uid string) error {
	err := datastore.Delete(tx, uid)
	if err != nil {
		return err
	}

	return nil
}
