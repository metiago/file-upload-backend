package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/metiago/zbx1/common/dml"
	"github.com/metiago/zbx1/common/env"
)

func FileUpload(u User, f File) error {

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

// FindaAllFilesByUsername get all files by a give username
func FindaAllFilesByUsername(username string) ([]File, error) {
	rows, err := env.GetConnection().Query(dml.FindAllFilesByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f File
	var fl = make([]File, 0)

	for rows.Next() {
		err = rows.Scan(&f.ID, &f.Name, &f.Ext, &f.Created, &f.Data)
		if err != nil {
			log.Printf("Error reading row %s", err)
			return nil, err
		}

		fl = append(fl, f)
	}

	return fl, err
}
