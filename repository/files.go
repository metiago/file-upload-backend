package repository

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/metiago/zbx1/common/dml"
	"github.com/metiago/zbx1/common/env"
)

func FileUpload(u User, f File) error {

	db := env.GetConnection()

	stmt, err := db.Prepare(dml.AddFile)
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := db.Begin()

	if err != nil {
		log.Println(err)
		return err
	}

	res, err := tx.Stmt(stmt).Exec(f.Name, f.Ext, f.Data, time.Now())

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err = db.Prepare(dml.AddUserFile)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = tx.Stmt(stmt).Exec(u.ID, id)
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
