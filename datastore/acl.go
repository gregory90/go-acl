package datastore

import (
	"database/sql"

	"bitbucket.org/pqstudio/go-acl/model"

	. "bitbucket.org/pqstudio/go-acl/db"
	. "bitbucket.org/pqstudio/go-webutils/logger"
)

const (
	table string = "acl"

	selectQuery string = `
        SELECT 
            lower(hex(uid)), 
            groupID, 
            lower(hex(objectID)), 
            permission, 
            action, 
            createdAt 
        FROM ` + table + " "

	insertQuery string = `
        INSERT  ` + table + ` SET 
            uid=unhex(?),
            groupID=?,
            objectID=unhex(?), 
            permission=?, 
            action=?, 
            createdAt=?
             `
	updateQuery string = `
        UPDATE  ` + table + ` SET 
            groupID=?,
            objectID=unhex(?), 
            permission=?, 
            action=? 
             `
	deleteQuery string = `
        DELETE FROM  ` + table + ` `
)

func scanSelectSingle(m *model.ACL, row *sql.Row) error {
	err := row.Scan(
		&m.UID,
		&m.GroupID,
		&m.Object,
		&m.Permission,
		&m.Action,
		&m.CreatedAt,
	)
	return err
}

func execUpdate(m *model.ACL, stmt *sql.Stmt) error {
	_, err := stmt.Exec(
		m.GroupID,
		m.Object,
		m.Permission,
		m.Action,
		m.UID,
	)

	return err
}

func execInsert(m *model.ACL, stmt *sql.Stmt) error {
	_, err := stmt.Exec(
		m.UID,
		m.GroupID,
		m.Object,
		m.Permission,
		m.Action,
		m.CreatedAt,
	)

	return err
}

func GetOne(uid string) (*model.ACL, error) {
	r := &model.ACL{}
	row := DB.QueryRow(selectQuery+"WHERE uid = unhex(?)", uid)

	err := scanSelectSingle(r, row)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func GetByGroupIDPermissionAndAction(groupID string, object string, permission string, action string) (*model.ACL, error) {
	r := &model.ACL{}
	Log.Debug(groupID + " " + object + " " + permission + " " + action)
	row := DB.QueryRow(selectQuery+"WHERE groupID = ? AND objectID = unhex(?) AND permission = ? AND action = ?", groupID, object, permission, action)

	err := scanSelectSingle(r, row)
	if err != nil {
		Log.Debug("%+v", err)
		return nil, err
	}
	Log.Debug("%+v", r)

	return r, nil
}

func Create(m *model.ACL) error {
	stmt, err := DB.Prepare(insertQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = execInsert(m, stmt)

	return err
}

func Update(m *model.ACL) error {
	stmt, err := DB.Prepare(updateQuery + "WHERE uid=unhex(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = execUpdate(m, stmt)
	return err
}

func Delete(uid string) error {
	stmt, err := DB.Prepare(deleteQuery + "WHERE uid=unhex(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(uid)
	if err != nil {
		return err
	}

	return nil
}
