package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/metiago/zbx1/common/dml"
	"github.com/metiago/zbx1/common/env"
	"github.com/metiago/zbx1/common/helper"
)

func AddUser(u *User) (*User, error) {

	u.Password, _ = helper.HashPassword(u.Password)

	db := env.GetConnection()

	tx, err := db.Begin()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var id int
	err = tx.QueryRow(dml.AddUser, u.Name, u.Email, u.Username, u.Password, time.Now()).Scan(&id)

	if err != nil {
		log.Printf("Error inserting user: %v", err)
		tx.Rollback()
		return nil, err
	}

	r, err := FindRoleByID(u.Role.ID)
	if r == nil {
		tx.Rollback()
		return nil, fmt.Errorf("Role with ID %d not found ", u.Role.ID)
	}

	_, err = tx.Exec(dml.AddUserRole, id, r.ID)
	if err != nil {
		log.Printf("Error inserting role for user: %v", err)
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return u, err
}

func UpdateUser(ID int, u *User) (*User, error) {

	u.Password, _ = helper.HashPassword(u.Password)

	db := env.GetConnection()

	stmt, err := db.Prepare(dml.UpdateUser)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	tx, err := db.Begin()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	_, err = tx.Stmt(stmt).Exec(u.Name, u.Email, u.Username, u.Password, time.Now(), ID)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, err
	}

	stmt, err = db.Prepare(dml.UpdateUserRole)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	r, err := FindRoleByID(u.Role.ID)
	if r == nil || r.ID == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("Role with ID %d not found ", u.Role.ID)
	}

	_, err = tx.Stmt(stmt).Exec(r.ID, ID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return u, err
}

func FindAllUsers() ([]User, error) {

	rows, err := env.GetConnection().Query(dml.FindAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u User
	var s = make([]User, 0)

	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.Username, &u.Created)
		if err != nil {
			log.Printf("Error reading row %s", err)
			return nil, err
		}

		u.Role, err = FindRoleByUserID(u.ID)
		if err != nil {
			return nil, err
		}

		s = append(s, u)
	}

	return s, err
}

func FindUserByID(ID int) (*User, error) {

	var u User

	stmt, err := env.GetConnection().Prepare(dml.FindUserByID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Created)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("User with ID %d not found", ID)
		return nil, err
	case err != nil:
		log.Printf("Error on find user %s", err)
		return nil, err
	default:
		return &u, nil
	}
}

func FindUserByUsername(username string) (*User, error) {

	var u User

	stmt, err := env.GetConnection().Prepare(dml.FindUserByUsername)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&u.ID, &u.Name, &u.Email, &u.Username, &u.Created)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("User with username %s not found", username)
		return nil, err
	case err != nil:
		log.Printf("Error on find user %s", err)
		return nil, err
	default:
		return &u, nil
	}
}

func DeleteUser(ID int) error {

	db := env.GetConnection()

	_, err := FindUserByID(ID)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(dml.DeleteUser)

	if err != nil {
		return err
	}

	defer stmt.Close()

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

// RepoAuthUser is to check if user's credentials are ok
func AuthUser(username string, password string) (bool, error) {

	log.Println("Authorizing user with username ", username)

	var u User

	stmt, err := env.GetConnection().Prepare(dml.AuthUser)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&u.Username, &u.Password)
	switch {
	case err == sql.ErrNoRows:
		log.Println(err)
		return false, nil
	case err != nil:
		log.Println(err)
		return false, err
	default:
		return helper.CheckPasswordHash(password, u.Password), nil
	}
}
