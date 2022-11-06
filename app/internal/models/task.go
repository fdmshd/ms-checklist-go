package models

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Task struct {
	Id          int    `json:"id"`
	UserID      int    `json:"user_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"descr"`
	IsDone      bool   `json:"is_done"`
	CreatedAt   string `json:"created_at"`
}

type TaskModel struct {
	DB *sql.DB
}

func (m *TaskModel) Get(id int) (*Task, error) {
	stmt := `SELECT id,user_id,name,descr,is_done,created_at FROM Tasks WHERE id =?`
	row := m.DB.QueryRow(stmt, id)
	task := new(Task)
	err := row.Scan(&task.Id, &task.UserID, &task.Name, &task.Description, &task.IsDone, &task.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return task, nil
}

func (m *TaskModel) GetByUser(userID int) ([]*Task, error) {
	stmt := `SELECT id,user_id,name,descr,is_done,created_at FROM Tasks WHERE user_id = ?`
	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}

	tasks := make([]*Task, 0)
	n := 0
	for rows.Next() {
		tasks = append(tasks, &Task{})
		err := rows.Scan(
			&tasks[n].Id, &tasks[n].UserID,
			&tasks[n].Name, &tasks[n].Description,
			&tasks[n].IsDone, &tasks[n].CreatedAt,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNoRecord
			} else {
				return nil, err
			}
		}
		n++
	}
	return tasks, nil
}

func (m *TaskModel) Add(t *Task) (int, error) {
	stmt := `INSERT INTO Tasks (user_id, name, descr,created_at)
	VALUES(?,?,?,?)`

	res, err := m.DB.Exec(stmt, t.UserID, t.Name, t.Description, t.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *TaskModel) Update(t *Task) error {
	stmt := `UPDATE Tasks SET name = ?, descr = ? WHERE id = ?`
	res, err := m.DB.Exec(stmt, t.Name, t.Description, t.Id)
	if err != nil {
		return err
	}
	affectedNum, _ := res.RowsAffected()
	if affectedNum == 0 {
		return ErrNoRecord
	}
	return nil
}

func (m *TaskModel) UpdateFields(fields map[string]interface{}, id int) error {
	stmt := "UPDATE Tasks SET "
	for name, val := range fields {
		stmt += fmt.Sprintf("%s=%v, ", name, val)
	}
	stmt = stmt[:len(stmt)-2] + " "
	stmt += "WHERE id = ?"
	res, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	affectedNum, _ := res.RowsAffected()
	if affectedNum == 0 {
		return ErrNoRecord
	}
	return nil
}

func (m *TaskModel) Delete(id int) error {
	stmt := `DELETE FROM Tasks WHERE id = ?`
	res, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	affectedNum, _ := res.RowsAffected()
	if affectedNum == 0 {
		return ErrNoRecord
	}
	return nil
}

func (m *TaskModel) DeleteByUser(userID int) error {
	stmt := `DELETE FROM Tasks WHERE user_id = ?`
	res, err := m.DB.Exec(stmt, userID)
	if err != nil {
		return err
	}
	affectedNum, _ := res.RowsAffected()
	if affectedNum == 0 {
		return ErrNoRecord
	}
	return nil
}
