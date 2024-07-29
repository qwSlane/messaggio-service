package repository

const (
	createMessageQuery = `INSERT INTO messages 	(message_id, content, created_at, updated_at, status)
	VALUES ($1, $2, now(), now(), $3) RETURNING message_id, content, created_at, updated_at`
)
