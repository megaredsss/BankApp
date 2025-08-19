-- name: GetUserIDByEmail :one
SELECT
    id
FROM
    users
WHERE
    email = $1
LIMIT
    1;

DELETE FROM
    users
WHERE
    id = $1;

---
-- name: GetUserNames :one
SELECT
    first_name,
    last_name
FROM
    users_profiles
WHERE
    users_id = $1
LIMIT
    1;

---
-- name: GetUserBalance :one
SELECT
    balance
FROM
    users_profiles
WHERE
    users_id = $1
LIMIT
    1;

---
-- name: CreateUser :one
INSERT INTO
    users(email, password)
VALUES
    ($1, $2) RETURNING id;

---
-- name: DeleteUser :exec
DELETE FROM
    users
WHERE id = $1;

---
-- name: CreateUserProfile :one
INSERT INTO
    users_profiles(users_id, first_name, last_name)
VALUES
    ($1, $2, $3) RETURNING *;

---
-- name: MoneyTransactionIncreaseByUserId :exec
UPDATE users_profiles SET balance = balance + $2
WHERE id = $1;

---
-- name: MoneyTransactionDecreaseByUserId :exec
UPDATE users_profiles SET balance = balance - $2
WHERE id = $1;
