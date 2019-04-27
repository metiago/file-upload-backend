package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/metiago/zbx1/common/dml"
	"github.com/metiago/zbx1/common/env"
)

var ErrFileExists = errors.New("File already exists")

func FileUpload(u User, f *File) error {

	err := fileExists(u, f)
	if err != nil {
		log.Println(err)
		return err
	}

	db := env.GetConnection()

	tx, err := db.Begin()

	if err != nil {
		log.Println(err)
		return err
	}

	var id int
	err = tx.QueryRow(dml.AddFile, f.Name, f.Ext, f.Data, time.Now()).Scan(&id)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(dml.AddUserFile, u.ID, id)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func FileDownload(ID int) (*File, error) {

	var f File

	stmt, err := env.GetConnection().Prepare(dml.FindFileByID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(ID).Scan(&f.ID, &f.Name, &f.Ext, &f.Data)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("File with ID %d not found", ID)
		return nil, err
	case err != nil:
		log.Printf("Error on find user %s", err)
		return nil, err
	default:
		return &f, nil
	}
}

func FindaAllFilesByUsername(username string, offset int) ([]File, error) {
	rows, err := env.GetConnection().Query(dml.FindAllFilesByUsername, username, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f File
	var fl = make([]File, 0)

	for rows.Next() {
		err = rows.Scan(&f.ID, &f.Name, &f.Ext, &f.Created, &f.Data, &f.Username)
		if err != nil {
			log.Printf("Error reading row %s", err)
			return nil, err
		}

		fl = append(fl, f)
	}

	return fl, err
}

func fileExists(u User, f *File) error {

	rows, err := env.GetConnection().Query(dml.FileExist, u.Username, f.Name)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		return ErrFileExists
	}

	return nil
}

func FindFileByID(ID int) (*File, error) {

	var f File

	stmt, err := env.GetConnection().Prepare(dml.FindFileByID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(&f.ID, &f.Name, &f.Ext, &f.Data)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("User with ID %d not found", ID)
		return nil, err
	case err != nil:
		log.Printf("Error on find file %s", err)
		return nil, err
	default:
		return &f, nil
	}
}

func DeleteFile(ID int) error {

	db := env.GetConnection()

	_, err := FindFileByID(ID)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(dml.DeleteFile)

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
