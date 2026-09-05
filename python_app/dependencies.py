import asyncpg
import jwt
from fastapi import HTTPException, Request

from auth.jwt import validate_token
from repositories.conversation import ConversationRepository
from repositories.message import MessageRepository
from services.chat import ChatService
from repositories.user import UserRepository

def get_db(request: Request) -> asyncpg.Pool:
    return request.app.state.db

def get_chat_service(request: Request) -> ChatService:
    db = get_db(request)

    conversation_repository = ConversationRepository(db)
    message_repository = MessageRepository(db)
    user_repository = UserRepository(db)

    return ChatService(
        conversation_repository=conversation_repository,
        message_repository=message_repository,
        user_repository=user_repository,
    )

def get_current_user_id(request: Request) -> int:
    authorization = request.headers.get("Authorization")

    if not authorization:
        raise HTTPException(
            status_code=401,
            detail="Authorization header is missing",
        )

    scheme, _, token = authorization.partition(" ")

    if scheme.lower() != "bearer" or not token:
        raise HTTPException(
            status_code=401,
            detail="Invalid authorization header",
        )

    try:
        return validate_token(token)
    except jwt.PyJWTError:
        raise HTTPException(
            status_code=401,
            detail="Invalid or expired token",
        )

def get_user_id_from_token(token: str) -> int:
    try:
        return validate_token(token)
    except jwt.PyJWTError:
        raise HTTPException(
            status_code=401,
            detail="Invalid or expired token",
        )