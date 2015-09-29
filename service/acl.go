package service

import (
	"database/sql"
	"time"

	"bitbucket.org/pqstudio/go-acl/datastore"
	"bitbucket.org/pqstudio/go-acl/model"

	groupModel "bitbucket.org/pqstudio/go-user-group/model"
	groupS "bitbucket.org/pqstudio/go-user-group/service"

	"bitbucket.org/pqstudio/go-webutils"
	. "bitbucket.org/pqstudio/go-webutils/db"
)

func GetOne(db *sql.DB, uid string) (*model.ACL, error) {
	r, err := datastore.GetOne(db, uid)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func Check(db *sql.DB, userUID string, object string, permission string, action string) (bool, error) {
	var gr interface{}
	err := Transact(db, func(tx *sql.Tx) error {
		var err error
		gr, err = groupS.GetByUserUID(tx, userUID, 100, 0)
		return err
	})
	if err != nil {
		return false, err
	}

	allowed, err := CheckGroups(db, gr, object, permission, action)

	return allowed, err
}

func CheckGroup(db *sql.DB, userUID string, group string) (bool, error) {
	var gr interface{}
	err := Transact(db, func(tx *sql.Tx) error {
		var err error
		grs, err = groupS.GetByUserUID(tx, userUID, 100, 0)
		return err
	})
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

func CheckGroups(db *sql.DB, groups []groupModel.Group, object string, permission string, action string) (bool, error) {
	var e error
	for _, group := range groups {
		_, err := datastore.GetByGroupIDPermissionAndAction(db, group.Name, object, permission, action)

		if err != nil {
			e = err
			continue
		}

		return true, nil
	}

	return false, e
}

func CreateFromModel(db *sql.DB, m *model.ACL) error {
	m.UID = utils.NewUUID()
	m.CreatedAt = time.Now().UTC()

	err := datastore.Create(db, m)
	if err != nil {
		return err
	}

	return nil
}

func Create(db *sql.DB, userID string, objectID string, permission string, actions []string) error {

	var e error

	for _, action := range actions {
		acl := &model.ACL{
			GroupID:    userID,
			Object:     objectID,
			Permission: permission,
			Action:     action,
		}

		err := CreateFromModel(db, acl)
		if err != nil {
			e = err
			continue
		}
	}

	return e
}

func Update(db *sql.DB, m *model.ACL) error {
	err := datastore.Update(db, m)
	if err != nil {
		return err
	}

	return nil
}

func Delete(db *sql.DB, uid string) error {
	err := datastore.Delete(db, uid)
	if err != nil {
		return err
	}

	return nil
}
