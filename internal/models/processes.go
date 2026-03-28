package models

import (
	"database/sql"
	"errors"
	"time"
)

// CREATE DATABASE sentinel CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
// USE sentinel;
// CREATE TABLE processes (id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,title VARCHAR(100) NOT NULL,startedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,active BOOLEAN NOT NULL DEFAULT FALSE);

// INSERT INTO processes (title, startedAt, active) VALUES ("google.com",UTC_TIMESTAMP(),TRUE);
// INSERT INTO processes (title, startedAt, active) VALUES ("youtube.com",DATE_ADD(UTC_TIMESTAMP(), INTERVAL 3 DAY),TRUE);
// INSERT INTO processes (title, startedAt, active) VALUES ("facebook.com",DATE_ADD(UTC_TIMESTAMP(), INTERVAL 1 DAY),FALSE);

// CREATE USER 'web'@'localhost';
// GRANT SELECT, INSERT, UPDATE, DELETE ON sentinel.* TO 'web'@'localhost';
// ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';

type Process struct {
	ID        int
	Title     string
	StartedAt time.Time
	Active    bool
}

type ProcessModel struct {
	DB *sql.DB
}

func (m *ProcessModel) Insert(title string) (int, error) {
	stmt := `INSERT INTO processes (title)VALUES(?)`
	res, err := m.DB.Exec(stmt, title)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *ProcessModel) Get(id int) (*Process, error) {
	stmt := `SELECT id, title, startedAt, active FROM processes WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)

	p := &Process{}

	err := row.Scan(&p.ID, &p.Title, &p.StartedAt, &p.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return p, nil
}
func (m *ProcessModel) Latest() ([]*Process, error) {
	stmt := `SELECT id, title, startedAt, active FROM processes ORDER BY id`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	process := []*Process{}
	for rows.Next() {
		p := &Process{}
		err = rows.Scan(&p.ID, &p.Title, &p.StartedAt, &p.Active)

		if err != nil {
			return nil, err
		} else {
			process = append(process, p)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return process, nil
}
