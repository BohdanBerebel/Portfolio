from datetime import datetime

from pydantic import BaseModel


class CreateConversationRequest(BaseModel):
    user2_id: int


class UserPreview(BaseModel):
    id: int
    email: str


class ConversationResponse(BaseModel):
    id: int
    user1_id: int
    user2_id: int
    other_user: UserPreview
    created_at: datetime