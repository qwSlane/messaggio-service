package repository

const (
	createMessageQuery = `INSERT INTO messages 	(message_id, content, created_at, updated_at, status)
	VALUES ($1, $2, now(), now(), $3) RETURNING message_id, content, created_at, updated_at`

	updateMessageQuery = `UPDATE messages SET status = true WHERE message_id = $1 RETURNING message_id, content, created_at, updated_at`

	getStatisticsQuery = `SELECT
    COUNT(*)::bigint AS all_messages_count,
    COUNT(CASE WHEN status = true THEN 1 END)::bigint AS processed_messages_count,
    (SELECT message_id FROM messages WHERE status = true ORDER BY updated_at DESC LIMIT 1) AS last_processed_id
FROM messages;
`
)
