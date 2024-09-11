-- name: ListJobApplicationStatus :many
SELECT
    *
FROM
    job_application_status;

-- name: GetJobApplicationStatusIdByName :one
SELECT
    id
FROM
    job_application_status
WHERE
    name = ?;

-- name: DeleteJobApplicationStatus :exec
DELETE FROM
    job_application_status
WHERE
    name = ?;

-- name: UpdateJobApplicationStatus :exec
UPDATE
    job_application_status
SET
    name = ?
WHERE
    id = ?;

-- name: InsertJobApplicationStatus :one
INSERT INTO
    job_application_status(name)
VALUES
    (?) RETURNING *;

-- name: ListJobApplication :many
SELECT
    job_application.*,
    job_application_status.name AS status_name,
    company.name AS company_name,
    (
        SELECT
            count(*)
        FROM
            event
        WHERE
            event.job_application_id = job_application.id
    ) AS event_cnt,
    (
        SELECT
            cast(max(event.date) AS integer)
        FROM
            event
        WHERE
            event.job_application_id = job_application.id
            AND event.date <= unixepoch()
    ) AS last_event,
    (
        SELECT
            cast(min(event.date) AS integer)
        FROM
            event
        WHERE
            event.job_application_id = job_application.id
            AND event.date >= unixepoch()
    ) AS next_event
FROM
    job_application
    INNER JOIN job_application_status ON job_application_status.id = job_application.status_id
    INNER JOIN company ON company.id = job_application.company_id;

-- name: GetJobApplicationIdByName :one
SELECT
    job_application.id
FROM
    job_application
    INNER JOIN company ON company.id = job_application.company_id
WHERE
    job_application.job_title = ?
    AND company.name = ?
LIMIT
    1;

-- name: DeleteJobApplication :exec
DELETE FROM
    job_application
WHERE
    id = ?;

-- name: UpdateJobApplication :exec
UPDATE
    job_application
SET
    company_id = ?,
    status_id = ?,
    job_title = ?,
    notes = ?
WHERE
    id = ?;

-- name: InsertJobApplication :one
INSERT INTO
    job_application(company_id, status_id, job_title, notes)
VALUES
    (?, ?, ?, ?) RETURNING *;