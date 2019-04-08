package repository

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/metiago/zbx1/common/dml"
	"github.com/metiago/zbx1/common/env"
	
)

func FindAllRoles() ([]Role, error) {

	rows, err := env.GetConnection().Query(dml.FindAllRoles)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var r Role
	var s = make([]Role, 0)

	for rows.Next() {
		err = rows.Scan(&r.ID, &r.Name, &r.Created)
		if err != nil {
			log.Printf("Error reading row %s", err)
			return nil, err
		}
		s = append(s, r)
	}

	return s, err
}

func AddRole(r *Role) (*Role, error) {

	db := env.GetConnection()

	stmt, err := db.Prepare(dml.AddRole)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	tx, err := db.Begin()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = tx.Stmt(stmt).Exec(r.Name, time.Now())

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return r, err
}

func UpdateRole(ID int, r *Role) (*Role, error) {

	db := env.GetConnection()

	stmt, err := db.Prepare(dml.UpdateRole)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	tx, err := db.Begin()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = tx.Stmt(stmt).Exec(r.Name, time.Now(), ID)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return r, err
}

func DeleteRole(ID int) error {

	db := env.GetConnection()

	_, err := FindRoleByID(ID)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(dml.DeleteRole)

	defer stmt.Close()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	defer tx.Commit()

	if err != nil {
		log.Println(err)
		return err
	}

	_, err = tx.Stmt(stmt).Exec(ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func FindRoleByID(ID int) (*Role, error) {

	var o Role

	stmt, err := env.GetConnection().Prepare(dml.FindRoleByID)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(&o.ID, &o.Name, &o.Created)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("Role with ID %d not found", ID)
		return &o, nil
	case err != nil:
		log.Printf("Error finding by roles %s", err)
		return nil, err
	default:
		return &o, nil
	}
}

func FindRoleByUserID(ID int) (*Role, error) {

	var o Role

	stmt, err := env.GetConnection().Prepare(dml.FindRoleByUserID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(&o.ID, &o.Name, &o.Created)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("Role with ID %d not found", ID)
		return &o, nil
	case err != nil:
		log.Printf("Error finding by roles by user %s", err)
		return nil, err
	default:
		return &o, nil
	}

}
