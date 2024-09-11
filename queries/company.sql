-- name: GetCompany :one
SELECT
    company.*,
    sqlc.embed(company_type)
FROM
    company
    INNER JOIN company_type_rel ON company_type_rel.company_id = company.id
    INNER JOIN company_type ON company_type.id = company_type_rel.company_type_id
WHERE
    company.id = ?
LIMIT
    1;

-- name: GetCompanyIdByName :one
SELECT
    company.id
FROM
    company
WHERE
    company.name = ?
LIMIT
    1;

-- name: ListCompany :many
SELECT
    company.*,
    sqlc.embed(company_type),
    (
        SELECT
            count(*)
        FROM
            job_application
        WHERE
            job_application.company_id = company.id
    ) AS application_cnt,
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
            INNER JOIN job_application ON job_application.id = event.job_application_id
        WHERE
            job_application.company_id = company.id
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
            INNER JOIN job_application ON job_application.id = event.job_application_id
        WHERE
            job_application.company_id = company.id
            AND event.date >= unixepoch()
    ) AS next_event
FROM
    company
    INNER JOIN company_type_rel ON company_type_rel.company_id = company.id
    INNER JOIN company_type ON company_type.id = company_type_rel.company_type_id;

-- name: InsertCompany :one
INSERT INTO
    company(name, notes)
VALUES
    (?, ?) RETURNING *;

-- name: UpdateCompany :exec
UPDATE
    company
SET
    name = ?,
    notes = ?
WHERE
    id = ?;

-- name: DeleteCompany :exec
DELETE FROM
    company
WHERE
    name = ?;

-- name: GetCompanyType :one
SELECT
    *
FROM
    company_type
WHERE
    name = ?
LIMIT
    1;

-- name: ListCompanyType :many
SELECT
    *,
    (
        SELECT
            count(*)
        FROM
            company
            INNER JOIN company_type_rel ON company_type_rel.company_id = company.id
        WHERE
            company_type_rel.company_type_id = company_type.id
    ) AS companies
FROM
    company_type;

-- name: InsertCompanyType :one
INSERT INTO
    company_type(name)
VALUES
    (?) RETURNING *;

-- name: UpdateCompanyType :exec
UPDATE
    company_type
SET
    name = ?
WHERE
    id = ?;

-- name: DeleteCompanyType :exec
DELETE FROM
    company_type
WHERE
    name = ?;

-- name: InsertCompanyTypeRel :one
INSERT INTO
    company_type_rel(company_id, company_type_id)
VALUES
    (?, ?) RETURNING *;

-- name: DeleteCompanyTypeRel :exec
DELETE FROM
    company_type_rel
WHERE
    company_id = ?;