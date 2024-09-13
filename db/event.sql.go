// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: event.sql

package db

import (
	"context"
)

const deleteEvent = `-- name: DeleteEvent :exec
DELETE FROM
    event
WHERE
    id = ?
`

func (q *Queries) DeleteEvent(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteEvent, id)
	return err
}

const deleteEventContact = `-- name: DeleteEventContact :exec
DELETE FROM
    event_contacts
WHERE
    event_id = ?
`

func (q *Queries) DeleteEventContact(ctx context.Context, eventID int64) error {
	_, err := q.db.ExecContext(ctx, deleteEventContact, eventID)
	return err
}

const deleteEventSource = `-- name: DeleteEventSource :exec
DELETE FROM
    event_source
WHERE
    name = ?
`

func (q *Queries) DeleteEventSource(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, deleteEventSource, name)
	return err
}

const getEventSourceIdByName = `-- name: GetEventSourceIdByName :one
SELECT
    id
FROM
    event_source
WHERE
    name = ?
`

func (q *Queries) GetEventSourceIdByName(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getEventSourceIdByName, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const insertEvent = `-- name: InsertEvent :one
INSERT INTO
    event(
        source_id,
        job_application_id,
        title,
        date,
        notes
    )
VALUES
    (?, ?, ?, ?, ?) RETURNING id, source_id, job_application_id, title, date, notes
`

type InsertEventParams struct {
	SourceID         int64  `json:"source_id"`
	JobApplicationID int64  `json:"job_application_id"`
	Title            string `json:"title"`
	Date             int64  `json:"date"`
	Notes            string `json:"notes"`
}

func (q *Queries) InsertEvent(ctx context.Context, arg InsertEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, insertEvent,
		arg.SourceID,
		arg.JobApplicationID,
		arg.Title,
		arg.Date,
		arg.Notes,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.SourceID,
		&i.JobApplicationID,
		&i.Title,
		&i.Date,
		&i.Notes,
	)
	return i, err
}

const insertEventContact = `-- name: InsertEventContact :one
INSERT INTO
    event_contacts(event_id, contact_id)
VALUES
    (?, ?) RETURNING id, event_id, contact_id
`

type InsertEventContactParams struct {
	EventID   int64 `json:"event_id"`
	ContactID int64 `json:"contact_id"`
}

func (q *Queries) InsertEventContact(ctx context.Context, arg InsertEventContactParams) (EventContact, error) {
	row := q.db.QueryRowContext(ctx, insertEventContact, arg.EventID, arg.ContactID)
	var i EventContact
	err := row.Scan(&i.ID, &i.EventID, &i.ContactID)
	return i, err
}

const insertEventSource = `-- name: InsertEventSource :one
INSERT INTO
    event_source(name)
VALUES
    (?) RETURNING id, name
`

func (q *Queries) InsertEventSource(ctx context.Context, name string) (EventSource, error) {
	row := q.db.QueryRowContext(ctx, insertEventSource, name)
	var i EventSource
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listEvent = `-- name: ListEvent :many
SELECT
    event.id, event.source_id, event.job_application_id, event.title, event.date, event.notes,
    event_source.name AS source,
    company.name AS company_name,
    job_application.job_title AS job_title,
    contact.fist_name AS contact_fist_name,
    contact.last_name AS contact_last_name
FROM
    event
    INNER JOIN event_source ON event_source.id = event.source_id
    INNER JOIN job_application ON job_application.id = event.job_application_id
    INNER JOIN company ON company.id = job_application.company_id
    LEFT JOIN event_contacts ON event_contacts.event_id = event.id
    LEFT JOIN contact ON contact.id = event_contacts.contact_id
    AND contact.company_id = company.id
ORDER BY
    event.date DESC,
    company.name
`

type ListEventRow struct {
	ID               int64   `json:"id"`
	SourceID         int64   `json:"source_id"`
	JobApplicationID int64   `json:"job_application_id"`
	Title            string  `json:"title"`
	Date             int64   `json:"date"`
	Notes            string  `json:"notes"`
	Source           string  `json:"source"`
	CompanyName      string  `json:"company_name"`
	JobTitle         string  `json:"job_title"`
	ContactFistName  *string `json:"contact_fist_name"`
	ContactLastName  *string `json:"contact_last_name"`
}

func (q *Queries) ListEvent(ctx context.Context) ([]ListEventRow, error) {
	rows, err := q.db.QueryContext(ctx, listEvent)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListEventRow{}
	for rows.Next() {
		var i ListEventRow
		if err := rows.Scan(
			&i.ID,
			&i.SourceID,
			&i.JobApplicationID,
			&i.Title,
			&i.Date,
			&i.Notes,
			&i.Source,
			&i.CompanyName,
			&i.JobTitle,
			&i.ContactFistName,
			&i.ContactLastName,
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

const listEventSource = `-- name: ListEventSource :many
SELECT
    id, name,
    (
        SELECT
            count(*)
        FROM
            event
        WHERE
            event.source_id = event_source.id
    ) AS EVENTS
FROM
    event_source
ORDER BY
    event_source.name
`

type ListEventSourceRow struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Events int64  `json:"events"`
}

func (q *Queries) ListEventSource(ctx context.Context) ([]ListEventSourceRow, error) {
	rows, err := q.db.QueryContext(ctx, listEventSource)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListEventSourceRow{}
	for rows.Next() {
		var i ListEventSourceRow
		if err := rows.Scan(&i.ID, &i.Name, &i.Events); err != nil {
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

const updateEvent = `-- name: UpdateEvent :exec
UPDATE
    event
SET
    source_id = ?,
    job_application_id = ?,
    title = ?,
    date = ?,
    notes = ?
WHERE
    id = ?
`

type UpdateEventParams struct {
	SourceID         int64  `json:"source_id"`
	JobApplicationID int64  `json:"job_application_id"`
	Title            string `json:"title"`
	Date             int64  `json:"date"`
	Notes            string `json:"notes"`
	ID               int64  `json:"id"`
}

func (q *Queries) UpdateEvent(ctx context.Context, arg UpdateEventParams) error {
	_, err := q.db.ExecContext(ctx, updateEvent,
		arg.SourceID,
		arg.JobApplicationID,
		arg.Title,
		arg.Date,
		arg.Notes,
		arg.ID,
	)
	return err
}

const updateEventSource = `-- name: UpdateEventSource :exec
UPDATE
    event_source
SET
    name = ?
WHERE
    id = ?
`

type UpdateEventSourceParams struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

func (q *Queries) UpdateEventSource(ctx context.Context, arg UpdateEventSourceParams) error {
	_, err := q.db.ExecContext(ctx, updateEventSource, arg.Name, arg.ID)
	return err
}
