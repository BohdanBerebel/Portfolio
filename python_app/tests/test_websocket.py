from datetime import datetime
from unittest.mock import AsyncMock, patch

import jwt
import pytest
from fastapi.testclient import TestClient
from starlette.websockets import WebSocketDisconnect

from database.models import Conversation, Message
from main import app
from routes.websocket import manager

def test_websocket_connects_with_valid_token():
    conversation = Conversation(
    id=10,
    user1_id=1,
    user2_id=2,
    created_at=datetime.now(),
    )

    mock_db = AsyncMock()

    mock_chat_service = AsyncMock()
    mock_chat_service.get_conversation.return_value = conversation

    with TestClient(app) as client:
        app.state.db = mock_db

        with patch(
            "routes.websocket.validate_token",
            return_value=1,
        ), patch(
            "routes.websocket.ChatService",
            return_value=mock_chat_service,
        ):
            with client.websocket_connect(
                "/api/chat/ws/10?token=valid-token"
            ) as websocket:
                assert websocket is not None

    mock_chat_service.get_conversation.assert_awaited_once_with(
        conversation_id=10,
        user_id=1,
    )

def test_websocket_rejects_missing_token():
    with TestClient(app) as client:
        app.state.db = AsyncMock()

        with pytest.raises(WebSocketDisconnect) as exc_info:
            with client.websocket_connect(
                "/api/chat/ws/10"
            ):
                pass

    assert exc_info.value.code == 1008

def test_websocket_rejects_invalid_token():
    with TestClient(app) as client:
        app.state.db = AsyncMock()

        with patch(
            "routes.websocket.validate_token",
            side_effect=jwt.InvalidTokenError,
        ):
            with pytest.raises(WebSocketDisconnect) as exc_info:
                with client.websocket_connect(
                    "/api/chat/ws/10?token=invalid-token"
                ):
                    pass

    assert exc_info.value.code == 1008

def test_websocket_rejects_non_member():
    mock_db = AsyncMock()

    mock_chat_service = AsyncMock()
    mock_chat_service.get_conversation.return_value = None

    with TestClient(app) as client:
        app.state.db = mock_db

        with patch(
            "routes.websocket.validate_token",
            return_value=3,
        ), patch(
            "routes.websocket.ChatService",
            return_value=mock_chat_service,
        ):
            with pytest.raises(WebSocketDisconnect) as exc_info:
                with client.websocket_connect(
                    "/api/chat/ws/10?token=valid-token"
                ):
                    pass

    assert exc_info.value.code == 1008

    mock_chat_service.get_conversation.assert_awaited_once_with(
        conversation_id=10,
        user_id=3,
    )

def test_websocket_sends_message_to_recipient():
    conversation = Conversation(
        id=10,
        user1_id=1,
        user2_id=2,
        created_at=datetime.now(),
    )

    message = Message(
        id=100,
        conversation_id=10,
        sender_id=1,
        content="Hello!",
        created_at=datetime.now(),
    )

    mock_db = AsyncMock()

    mock_chat_service = AsyncMock()
    mock_chat_service.get_conversation.return_value = conversation
    mock_chat_service.create_message.return_value = message

    with TestClient(app) as client:
        app.state.db = mock_db

        with patch(
            "routes.websocket.validate_token",
            return_value=1,
        ), patch(
            "routes.websocket.ChatService",
            return_value=mock_chat_service,
        ), patch(
            "routes.websocket.manager.send_to_user",
            new_callable=AsyncMock,
        ) as mock_send:
            with client.websocket_connect(
                "/api/chat/ws/10?token=valid-token"
            ) as websocket:
                websocket.send_text("Hello!")

    mock_chat_service.create_message.assert_awaited_once_with(
        conversation_id=10,
        sender_id=1,
        content="Hello!",
    )

    message_data = {
        "id": 100,
        "conversation_id": 10,
        "sender_id": 1,
        "content": "Hello!",
        "created_at": message.created_at.isoformat(),
    }

    assert mock_send.await_count == 2

    mock_send.assert_any_await(1, message_data)
    mock_send.assert_any_await(2, message_data)

def test_websocket_disconnect_removes_connection():
    conversation = Conversation(
        id=10,
        user1_id=1,
        user2_id=2,
        created_at=datetime.now(),
    )

    mock_db = AsyncMock()

    mock_chat_service = AsyncMock()
    mock_chat_service.get_conversation.return_value = conversation

    with TestClient(app) as client:
        app.state.db = mock_db

        with patch(
            "routes.websocket.validate_token",
            return_value=1,
        ), patch(
            "routes.websocket.ChatService",
            return_value=mock_chat_service,
        ):
            with client.websocket_connect(
                "/api/chat/ws/10?token=valid-token"
            ):
                assert 1 in manager.connections
                assert len(manager.connections[1]) == 1

            assert 1 not in manager.connections