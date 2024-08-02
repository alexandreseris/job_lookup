// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: job_application.sql

package db

import (
	"context"
)

const deleteJobApplication = `-- name: DeleteJobApplication :exec
DELETE FROM
    job_application
WHERE
    id = ?
`

func (q *Queries) DeleteJobApplication(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteJobApplication, id)
	return err
}

const deleteJobApplicationStatus = `-- name: DeleteJobApplicationStatus :exec
DELETE FROM
    job_application_status
WHERE
    name = ?
`

func (q *Queries) DeleteJobApplicationStatus(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, deleteJobApplicationStatus, name)
	return err
}

const getJobApplicationIdByName = `-- name: GetJobApplicationIdByName :one
SELECT
    job_application.id
FROM
    job_application
    INNER JOIN company ON company.id = job_application.company_id
WHERE
    job_application.job_title = ?
    AND company.name = ?
LIMIT
    1
`

type GetJobApplicationIdByNameParams struct {
	JobTitle string `json:"job_title"`
	Name     string `json:"name"`
}

func (q *Queries) GetJobApplicationIdByName(ctx context.Context, arg GetJobApplicationIdByNameParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getJobApplicationIdByName, arg.JobTitle, arg.Name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getJobApplicationStatusIdByName = `-- name: GetJobApplicationStatusIdByName :one
SELECT
    id
FROM
    job_application_status
WHERE
    name = ?
`

func (q *Queries) GetJobApplicationStatusIdByName(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getJobApplicationStatusIdByName, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const insertJobApplication = `-- name: InsertJobApplication :one
INSERT INTO
    job_application(company_id, status_id, job_title, notes)
VALUES
    (?, ?, ?, ?) RETURNING id, company_id, status_id, job_title, notes
`

type InsertJobApplicationParams struct {
	CompanyID int64  `json:"company_id"`
	StatusID  int64  `json:"status_id"`
	JobTitle  string `json:"job_title"`
	Notes     string `json:"notes"`
}

func (q *Queries) InsertJobApplication(ctx context.Context, arg InsertJobApplicationParams) (JobApplication, error) {
	row := q.db.QueryRowContext(ctx, insertJobApplication,
		arg.CompanyID,
		arg.StatusID,
		arg.JobTitle,
		arg.Notes,
	)
	var i JobApplication
	err := row.Scan(
		&i.ID,
		&i.CompanyID,
		&i.StatusID,
		&i.JobTitle,
		&i.Notes,
	)
	return i, err
}

const insertJobApplicationStatus = `-- name: InsertJobApplicationStatus :one
INSERT INTO
    job_application_status(name)
VALUES
    (?) RETURNING id, name
`

func (q *Queries) InsertJobApplicationStatus(ctx context.Context, name string) (JobApplicationStatus, error) {
	row := q.db.QueryRowContext(ctx, insertJobApplicationStatus, name)
	var i JobApplicationStatus
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listJobApplication = `-- name: ListJobApplication :many
SELECT
    job_application.id, job_application.company_id, job_application.status_id, job_application.job_title, job_application.notes,
    job_application_status.name AS status_name,
    company.name AS company_name
FROM
    job_application
    INNER JOIN job_application_status ON job_application_status.id = job_application.status_id
    INNER JOIN company ON company.id = job_application.company_id
`

type ListJobApplicationRow struct {
	ID          int64  `json:"id"`
	CompanyID   int64  `json:"company_id"`
	StatusID    int64  `json:"status_id"`
	JobTitle    string `json:"job_title"`
	Notes       string `json:"notes"`
	StatusName  string `json:"status_name"`
	CompanyName string `json:"company_name"`
}

func (q *Queries) ListJobApplication(ctx context.Context) ([]ListJobApplicationRow, error) {
	rows, err := q.db.QueryContext(ctx, listJobApplication)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListJobApplicationRow{}
	for rows.Next() {
		var i ListJobApplicationRow
		if err := rows.Scan(
			&i.ID,
			&i.CompanyID,
			&i.StatusID,
			&i.JobTitle,
			&i.Notes,
			&i.StatusName,
			&i.CompanyName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listJobApplicationStatus = `-- name: ListJobApplicationStatus :many
SELECT
    id, name
FROM
    job_application_status
`

func (q *Queries) ListJobApplicationStatus(ctx context.Context) ([]JobApplicationStatus, error) {
	rows, err := q.db.QueryContext(ctx, listJobApplicationStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []JobApplicationStatus{}
	for rows.Next() {
		var i JobApplicationStatus
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateJobApplication = `-- name: UpdateJobApplication :exec
UPDATE
    job_application
SET
    company_id = ?,
    status_id = ?,
    job_title = ?,
    notes = ?
WHERE
    id = ?
`

type UpdateJobApplicationParams struct {
	CompanyID int64  `json:"company_id"`
	StatusID  int64  `json:"status_id"`
	JobTitle  string `json:"job_title"`
	Notes     string `json:"notes"`
	ID        int64  `json:"id"`
}

func (q *Queries) UpdateJobApplication(ctx context.Context, arg UpdateJobApplicationParams) error {
	_, err := q.db.ExecContext(ctx, updateJobApplication,
		arg.CompanyID,
		arg.StatusID,
		arg.JobTitle,
		arg.Notes,
		arg.ID,
	)
	return err
}

const updateJobApplicationStatus = `-- name: UpdateJobApplicationStatus :exec
UPDATE
    job_application_status
SET
    name = ?
WHERE
    id = ?
`

type UpdateJobApplicationStatusParams struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

func (q *Queries) UpdateJobApplicationStatus(ctx context.Context, arg UpdateJobApplicationStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateJobApplicationStatus, arg.Name, arg.ID)
	return err
}
