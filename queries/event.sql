-- name: ListEventSource :many
SELECT
    *,
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
    event_source.name;

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
    company.name;

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