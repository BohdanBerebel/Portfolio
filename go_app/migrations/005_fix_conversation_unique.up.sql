ALTER TABLE conversations
DROP CONSTRAINT conversations_user1_id_user2_id_key;

CREATE UNIQUE INDEX idx_conversations_users_unique
ON conversations (
    LEAST(user1_id, user2_id),
    GREATEST(user1_id, user2_id)
);
