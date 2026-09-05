from dataclasses import dataclass

import asyncpg


@dataclass
class User:
    id: int
    email: str


class UserRepository:
    def __init__(self, db: asyncpg.Pool):
        self.db = db

    async def get_by_email(
        self,
        email: str,
    ) -> User | None:
        row = await self.db.fetchrow(
            """
            SELECT id, email
            FROM users
            WHERE email = $1
            """,
            email,
        )

        if row is None:
            return None

        return User(
            id=row["id"],
            email=row["email"],
        )