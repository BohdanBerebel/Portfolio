from pydantic import BaseModel


class CreateMessageRequest(BaseModel):
    content: str