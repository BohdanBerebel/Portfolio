import jwt
from fastapi import APIRouter, WebSocket, WebSocketDisconnect

from auth.jwt import validate_token
from repositories.conversation import ConversationRepository
from repositories.message import MessageRepository
from repositories.user import UserRepository
from services.chat import ChatService
from websocket.manager import ConnectionManager


router = APIRouter(
    prefix="/api/chat",
    tags=["websocket"],
)

manager = ConnectionManager()

@router.websocket("/ws/{conversation_id}")
async def websocket_endpoint(
    websocket: WebSocket,
    conversation_id: int,
):
    token = websocket.query_params.get("token")

    if not token:
        await websocket.close(code=1008)
        return

    try:
        user_id = validate_token(token)
    except jwt.PyJWTError:
        await websocket.close(code=1008)
        return

    db = websocket.app.state.db

    chat_service = ChatService(
        conversation_repository=ConversationRepository(db),
        message_repository=MessageRepository(db),
        user_repository=UserRepository(db),
    )

    conversation = await chat_service.get_conversation(
        conversation_id=conversation_id,
        user_id=user_id,
    )

    if conversation is None:
        await websocket.close(code=1008)
        return

    await manager.connect(user_id, websocket)

    try:
        while True:
            content = await websocket.receive_text()

            message = await chat_service.create_message(
                conversation_id=conversation_id,
                sender_id=user_id,
                content=content,
            )

            if message.sender_id == conversation.user1_id:
                recipient_id = conversation.user2_id
            else:
                recipient_id = conversation.user1_id

            message_data = {
                "id": message.id,
                "conversation_id": message.conversation_id,
                "sender_id": message.sender_id,
                "content": message.content,
                "created_at": message.created_at.isoformat(),
            }

            await manager.send_to_user(
                user_id,
                message_data,
            )

            await manager.send_to_user(
                recipient_id,
                message_data,
            )
            
    except WebSocketDisconnect:
        manager.disconnect(user_id, websocket)