from dataclasses import dataclass
from datetime import datetime

import asyncpg

from database.models import Conversation


@dataclass
class ConversationWithUser:
    id: int
    user1_id: int
    user2_id: int
    other_user_id: int
    other_user_email: str
    created_at: datetime


class ConversationRepository:
    def __init__(self, db: asyncpg.Pool):
        self.db = db

    async def create(
        self,
        user1_id: int,
        user2_id: int,
    ) -> Conversation:
        row = await self.db.fetchrow(
            """
            INSERT INTO conversations (user1_id, user2_id)
            VALUES ($1, $2)
            ON CONFLICT DO NOTHING
            RETURNING id, user1_id, user2_id, created_at
            """,
            user1_id,
            user2_id,
        )

        if row is None:
            row = await self.db.fetchrow(
                """
                SELECT id, user1_id, user2_id, created_at
                FROM conversations
                WHERE LEAST(user1_id, user2_id) = LEAST($1::integer, $2::integer)
                AND GREATEST(user1_id, user2_id) = GREATEST($1::integer, $2::integer)
                """,
                user1_id,
                user2_id,
            )

        return Conversation(
            id=row["id"],
            user1_id=row["user1_id"],
            user2_id=row["user2_id"],
            created_at=row["created_at"],
        )

    async def get_by_id(
        self,
        conversation_id: int,
    ) -> Conversation | None:
        row = await self.db.fetchrow(
            """
            SELECT id, user1_id, user2_id, created_at
            FROM conversations
            WHERE id = $1
            """,
            conversation_id,
        )

        if row is None:
            return None

        return Conversation(
            id=row["id"],
            user1_id=row["user1_id"],
            user2_id=row["user2_id"],
            created_at=row["created_at"],
        )

    async def get_for_user(
        self,
        conversation_id: int,
        user_id: int,
    ) -> Conversation | None:
        row = await self.db.fetchrow(
            """
            SELECT id, user1_id, user2_id, created_at
            FROM conversations
            WHERE id = $1
              AND (user1_id = $2 OR user2_id = $2)
            """,
            conversation_id,
            user_id,
        )

        if row is None:
            return None

        return Conversation(
            id=row["id"],
            user1_id=row["user1_id"],
            user2_id=row["user2_id"],
            created_at=row["created_at"],
        )

    async def list_for_user(
        self,
        user_id: int,
    ) -> list[ConversationWithUser]:
        rows = await self.db.fetch(
            """
            SELECT
                c.id,
                c.user1_id,
                c.user2_id,
                c.created_at,
                u.id AS other_user_id,
                u.email AS other_user_email
            FROM conversations c
            JOIN users u
              ON u.id = CASE
                  WHEN c.user1_id = $1 THEN c.user2_id
                  ELSE c.user1_id
              END
            WHERE c.user1_id = $1
               OR c.user2_id = $1
            ORDER BY c.created_at DESC
            """,
            user_id,
        )

        return [
            ConversationWithUser(
                id=row["id"],
                user1_id=row["user1_id"],
                user2_id=row["user2_id"],
                other_user_id=row["other_user_id"],
                other_user_email=row["other_user_email"],
                created_at=row["created_at"],
            )
            for row in rows
        ]