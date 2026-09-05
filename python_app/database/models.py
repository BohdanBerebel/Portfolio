from dataclasses import dataclass
from datetime import datetime


@dataclass
class Conversation:
    id: int
    user1_id: int
    user2_id: int
    created_at: datetime


@dataclass
class Message:
    id: int
    conversation_id: int
    sender_id: int
    content: str
    created_at: datetime

@dataclass
class User:
    id: int
    email: str