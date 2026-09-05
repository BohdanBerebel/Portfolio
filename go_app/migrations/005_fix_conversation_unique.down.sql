DROP INDEX IF EXISTS idx_conversations_users_unique;

ALTER TABLE conversations
ADD CONSTRAINT conversations_user1_id_user2_id_key
UNIQUE (user1_id, user2_id);