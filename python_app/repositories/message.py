import asyncpg

from database.models import Message


class MessageRepository:
    def __init__(self, db: asyncpg.Pool):
        self.db = db

    async def create(
        self,
        conversation_id: int,
        sender_id: int,
        content: str,
    ) -> Message:
        row = await self.db.fetchrow(
            """
            INSERT INTO messages (
                conversation_id,
                sender_id,
                content
            )
            VALUES ($1, $2, $3)
            RETURNING id, conversation_id, sender_id, content, created_at
            """,
            conversation_id,
            sender_id,
            content,
        )

        return Message(
            id=row["id"],
            conversation_id=row["conversation_id"],
            sender_id=row["sender_id"],
            content=row["content"],
            created_at=row["created_at"],
        )

    async def list_by_conversation(
        self,
        conversation_id: int,
    ) -> list[Message]:
        rows = await self.db.fetch(
            """
            SELECT
                id,
                conversation_id,
                sender_id,
                content,
                created_at
            FROM messages
            WHERE conversation_id = $1
            ORDER BY created_at ASC, id ASC
            """,
            conversation_id,
        )

        return [
            Message(
                id=row["id"],
                conversation_id=row["conversation_id"],
                sender_id=row["sender_id"],
                content=row["content"],
                created_at=row["created_at"],
            )
            for row in rows
        ]