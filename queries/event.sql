-- name: ListEventSource :many
SELECT
    *
FROM
    event_source;

-- name: GetEventSourceIdByName :one
SELECT
    id
FROM
    event_source
WHERE
    name = ?;

-- name: DeleteEventSource :exec
DELETE FROM
    event_source
WHERE
    name = ?;

-- name: UpdateEventSource :exec
UPDATE
    event_source
SET
    name = ?
WHERE
    id = ?;

-- name: InsertEventSource :one
INSERT INTO
    event_source(name)
VALUES
    (?) RETURNING *;

-- name: DeleteEventContact :exec
DELETE FROM
    event_contacts
WHERE
    event_id = ?;

-- name: InsertEventContact :one
INSERT INTO
    event_contacts(event_id, contact_id)
VALUES
    (?, ?) RETURNING *;

-- name: ListEvent :many
SELECT
    event.*,
    event_source.name AS source,
    company.name AS company_name,
    job_application.job_title AS job_title,
    sqlc.embed(contact)
FROM
    event
    INNER JOIN event_source ON event_source.id = event.source_id
    INNER JOIN job_application ON job_application.id = event.job_application_id
    INNER JOIN company ON company.id = job_application.company_id
    LEFT JOIN event_contacts ON event_contacts.event_id = event.id
    LEFT JOIN contact ON contact.id = event_contacts.contact_id;

-- name: DeleteEvent :exec
DELETE FROM
    event
WHERE
    id = ?;

-- name: UpdateEvent :exec
UPDATE
    event
SET
    source_id = ?,
    job_application_id = ?,
    title = ?,
    date = ?,
    notes = ?
WHERE
    id = ?;

-- name: InsertEvent :one
INSERT INTO
    event(
        source_id,
        job_application_id,
        title,
        date,
        notes
    )
VALUES
    (?, ?, ?, ?, ?) RETURNING *;