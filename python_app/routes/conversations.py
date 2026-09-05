from fastapi import APIRouter, Depends, HTTPException

from dependencies import get_chat_service, get_current_user_id
from schemas.conversation import (
    ConversationResponse,
    CreateConversationRequest,
    UserPreview,
)
from services.chat import ChatService
from schemas.message import CreateMessageRequest


router = APIRouter(
    prefix="/api/chat/conversations",
    tags=["conversations"],
)

@router.get("/", response_model=list[ConversationResponse])
async def list_conversations(
    user_id: int = Depends(get_current_user_id),
    chat_service: ChatService = Depends(get_chat_service),
):
    conversations = await chat_service.list_conversations(user_id)

    return [
        ConversationResponse(
            id=conversation.id,
            user1_id=conversation.user1_id,
            user2_id=conversation.user2_id,
            other_user=UserPreview(
                id=conversation.other_user_id,
                email=conversation.other_user_email,
            ),
            created_at=conversation.created_at,
        )
        for conversation in conversations
    ]

@router.get("/{conversation_id}")
async def get_conversation(
    conversation_id: int,
    user_id: int = Depends(get_current_user_id),
    chat_service: ChatService = Depends(get_chat_service),
):
    conversation = await chat_service.get_conversation(
        conversation_id=conversation_id,
        user_id=user_id,
    )

    if conversation is None:
        raise HTTPException(
            status_code=404,
            detail="Conversation not found",
        )

    return {
        "id": conversation.id,
        "user1_id": conversation.user1_id,
        "user2_id": conversation.user2_id,
        "created_at": conversation.created_at,
    }

@router.get("/{conversation_id}/messages")
async def list_messages(
    conversation_id: int,
    user_id: int = Depends(get_current_user_id),
    chat_service: ChatService = Depends(get_chat_service),
):
    try:
        messages = await chat_service.list_messages(
            conversation_id=conversation_id,
            user_id=user_id,
        )
    except ValueError as error:
        raise HTTPException(
            status_code=404,
            detail=str(error),
        )

    return [
        {
            "id": message.id,
            "conversation_id": message.conversation_id,
            "sender_id": message.sender_id,
            "content": message.content,
            "created_at": message.created_at,
        }
        for message in messages
    ]    

@router.post("/", status_code=201)
async def create_conversation(
    data: CreateConversationRequest,
    user_id: int = Depends(get_current_user_id),
    chat_service: ChatService = Depends(get_chat_service),
):
    try:
        conversation = await chat_service.create_conversation(
            user1_id=user_id,
            user2_id=data.user2_id,
        )

        conversations = await chat_service.list_conversations(user_id)

        created_conversation = next(
            item
            for item in conversations
            if item.id == conversation.id
        )

    except ValueError as error:
        raise HTTPException(
            status_code=400,
            detail=str(error),
        )

    return {
        "id": created_conversation.id,
        "user1_id": created_conversation.user1_id,
        "user2_id": created_conversation.user2_id,
        "other_user": {
            "id": created_conversation.other_user_id,
            "email": created_conversation.other_user_email,
        },
        "created_at": created_conversation.created_at,
    }

@router.post("/{conversation_id}/messages", status_code=201)
async def create_message(
    conversation_id: int,
    data: CreateMessageRequest,
    user_id: int = Depends(get_current_user_id),
    chat_service: ChatService = Depends(get_chat_service),
):
    try:
        message = await chat_service.create_message(
            conversation_id=conversation_id,
            sender_id=user_id,
            content=data.content,
        )
    except ValueError as error:
        raise HTTPException(
            status_code=404,
            detail=str(error),
        )

    return {
        "id": message.id,
        "conversation_id": message.conversation_id,
        "sender_id": message.sender_id,
        "content": message.content,
        "created_at": message.created_at,
    }



@router.get("/test")
async def test_chat_service(
    chat_service: ChatService = Depends(get_chat_service),
    user_id: int = Depends(get_current_user_id),
):
    return {
        "status": "chat service connected",
        "user_id": user_id,
    }