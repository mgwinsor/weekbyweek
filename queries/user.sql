-- name: CreateUser :one
INSERT INTO user (id, username, email, password_hash, date_of_birth)
VALUES (
	?1,
	?2,
	?3,
	?4,
	?5
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM user
WHERE username = ?1;

-- name: GetUserByEmail :one
SELECT * FROM user
WHERE email = ?1;

-- name: UpdateUserEmail :exec
UPDATE user
SET email = ?1, updated_at = CURRENT_TIMESTAMP
WHERE id = ?2;

-- name: UpdateUserPassword :exec
UPDATE user
SET password_hash = ?1, updated_at = CURRENT_TIMESTAMP
WHERE id = ?2;

-- name: UpdateUserDateOfBirth :exec
UPDATE user
SET date_of_birth = ?1, updated_at = CURRENT_TIMESTAMP
WHERE id = ?2;

-- name: DeleteUser :exec
DELETE FROM user
WHERE id = ?1;
