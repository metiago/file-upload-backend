package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/metiago/zbx1/common/dml"
	"github.com/metiago/zbx1/common/env"
	"github.com/metiago/zbx1/common/helper"
)

var ErrUsernameExists = errors.New("Username already exists")
var ErrMatchPassword = errors.New("Old password not match")

func AddUser(u *User) (*User, error) {

	existentUser, err := userExists(u)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existentUser.ID == 0 {

		u.Password, _ = helper.EncryptPassword(u.Password)

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

		tx.Commit()

		return u, nil
	}

	return nil, ErrUsernameExists
}

func UpdateUser(u *User) (*User, error) {

	existentUser, err := userExists(u)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existentUser.ID != 0 && existentUser.ID != u.ID {
		return nil, ErrUsernameExists
	}

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

	_, err = tx.Stmt(stmt).Exec(u.Name, u.Email, u.Username, time.Now(), u.ID)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, err
	}

	defer stmt.Close()

	tx.Commit()

	return u, nil

}

func UpdateUserPassword(u *User) (*User, error) {

	existentUser, err := FindUserByID(u.ID)
	if err != nil {
		return nil, err
	}

	match := helper.CheckPasswordHash(u.Password, existentUser.Password)

	if match {

		db := env.GetConnection()

		stmt, err := db.Prepare(dml.UpdateUserPassword)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		tx, err := db.Begin()

		if err != nil {
			log.Println(err)
			return nil, err
		}

		p, _ := helper.EncryptPassword(u.UpdatedPassword)
		_, err = tx.Stmt(stmt).Exec(p, time.Now(), u.ID)

		if err != nil {
			log.Println(err)
			tx.Rollback()
			return nil, err
		}

		defer stmt.Close()

		tx.Commit()

		return u, err
	}

	return nil, ErrMatchPassword
}

func FindAllUsers() ([]User, error) {

	rows, err := env.GetConnection().Query(dml.FindAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u User
	var users = make([]User, 0)

	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.Username, &u.Created)
		if err != nil {
			log.Printf("Error reading row %s", err)
			return nil, err
		}

		users = append(users, u)
	}

	return users, err
}

func FindUserByID(ID int) (*User, error) {

	var u User

	stmt, err := env.GetConnection().Prepare(dml.FindUserByID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(&u.ID, &u.Name, &u.Email, &u.Username, &u.Password, &u.Created)

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
	err = stmt.QueryRow(username).Scan(&u.ID, &u.Name, &u.Email, &u.Username, &u.Password, &u.Created)
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

func AuthUser(username string, password string) (bool, error) {

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

func userExists(user *User) (*User, error) {

	stmt, err := env.GetConnection().Prepare(dml.FindUserByUsername)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var u User
	err = stmt.QueryRow(user.Username).Scan(&u.ID, &u.Name, &u.Email, &u.Username, &u.Password, &u.Created)
	return &u, err
}
