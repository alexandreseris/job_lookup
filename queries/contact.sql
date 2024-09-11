-- name: ListContact :many
SELECT
    contact.*,
    company.name AS company_name,
    (
        SELECT
            cast(
                CASE
                    WHEN max(event.date) IS NULL THEN 0
                    ELSE max(event.date)
                END AS integer
            )
        FROM
            event
            INNER JOIN event_contacts ON event_contacts.event_id = event.id
        WHERE
            event_contacts.contact_id = contact.id
            AND event.date <= unixepoch()
    ) AS last_event,
    (
        SELECT
            cast(
                CASE
                    WHEN min(event.date) IS NULL THEN 0
                    ELSE min(event.date)
                END AS integer
            )
        FROM
            event
            INNER JOIN event_contacts ON event_contacts.event_id = event.id
        WHERE
            event_contacts.contact_id = contact.id
            AND event.date >= unixepoch()
    ) AS next_event
FROM
    contact
    INNER JOIN company ON company.id = contact.company_id;

-- name: GetContactIdByNames :one
SELECT
    contact.id
FROM
    contact
    INNER JOIN company ON company.id = contact.company_id
WHERE
    contact.fist_name = ?
    AND contact.last_name = ?
    AND company.name = ?
LIMIT
    1;

-- name: DeleteContact :exec
DELETE FROM
    contact
WHERE
    id = ?;

-- name: UpdateContact :exec
UPDATE
    contact
SET
    company_id = ?,
    job_position = ?,
    fist_name = ?,
    last_name = ?,
    email = ?,
    phone_number = ?,
    notes = ?
WHERE
    id = ?;

-- name: InsertContact :one
INSERT INTO
    contact (
        company_id,
        job_position,
        fist_name,
        last_name,
        email,
        phone_number,
        notes
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?) RETURNING *;