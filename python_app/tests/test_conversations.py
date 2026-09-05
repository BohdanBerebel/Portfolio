from datetime import datetime
from unittest.mock import AsyncMock

import pytest
from fastapi.testclient import TestClient

from database.models import Conversation, Message
from repositories.conversation import ConversationWithUser
from dependencies import get_chat_service, get_current_user_id
from main import app


client = TestClient(app)

@pytest.fixture
def mock_chat_service():
    return AsyncMock()

@pytest.fixture
def authenticated_client(mock_chat_service):
    app.dependency_overrides[get_current_user_id] = lambda: 1
    app.dependency_overrides[get_chat_service] = lambda: mock_chat_service

    yield mock_chat_service

    app.dependency_overrides.clear()

def test_list_conversations_requires_authentication():
    response = client.get("/api/chat/conversations/")

    assert response.status_code == 401
    assert response.json() == {
        "detail": "Authorization header is missing"
    }

def test_list_conversations_rejects_invalid_authentication_scheme():
    response = client.get(
        "/api/chat/conversations/",
        headers={"Authorization": "Basic abc123"},
    )

    assert response.status_code == 401
    assert response.json() == {
        "detail": "Invalid authorization header"
    }

def test_list_conversations_rejects_empty_bearer_token():
    response = client.get(
        "/api/chat/conversations/",
        headers={"Authorization": "Bearer"},
    )

    assert response.status_code == 401
    assert response.json() == {
        "detail": "Invalid authorization header"
    }

def test_list_conversations(authenticated_client):
    mock_service = authenticated_client

    conversations = [
        ConversationWithUser(
            id=1,
            user1_id=1,
            user2_id=2,
            other_user_id=2,
            other_user_email="admin@tubely.com",
            created_at=datetime.now(),
        ),
        ConversationWithUser(
            id=2,
            user1_id=1,
            user2_id=3,
            other_user_id=3,
            other_user_email="user3@example.com",
            created_at=datetime.now(),
        ),
    ]

    mock_service.list_conversations.return_value = conversations

    response = client.get(
        "/api/chat/conversations/",
        headers={"Authorization": "Bearer fake-token"},
    )

    assert response.status_code == 200

    data = response.json()

    assert data[0]["id"] == 1
    assert data[0]["user1_id"] == 1
    assert data[0]["user2_id"] == 2
    assert data[0]["other_user"] == {
        "id": 2,
        "email": "admin@tubely.com",
    }

    assert data[1]["id"] == 2
    assert data[1]["user1_id"] == 1
    assert data[1]["user2_id"] == 3
    assert data[1]["other_user"] == {
        "id": 3,
        "email": "user3@example.com",
    }
    mock_service.list_conversations.assert_awaited_once_with(1)

def test_create_conversation(authenticated_client):
    mock_service = authenticated_client

    conversation = Conversation(
        id=10,
        user1_id=1,
        user2_id=2,
        created_at=datetime.now(),
    )

    conversation_with_user = ConversationWithUser(
        id=10,
        user1_id=1,
        user2_id=2,
        other_user_id=2,
        other_user_email="admin@tubely.com",
        created_at=conversation.created_at,
    )

    mock_service.create_conversation.return_value = conversation
    mock_service.list_conversations.return_value = [
        conversation_with_user,
    ]

    response = client.post(
        "/api/chat/conversations/",
        json={"user2_id": 2},
        headers={"Authorization": "Bearer fake-token"},
    )

    assert response.status_code == 201

    data = response.json()

    assert data["id"] == 10
    assert data["user1_id"] == 1
    assert data["user2_id"] == 2
    assert data["other_user"]["id"] == 2
    assert data["other_user"]["email"] == "admin@tubely.com"

    mock_service.create_conversation.assert_awaited_once_with(
        user1_id=1,
        user2_id=2,
    )

    mock_service.list_conversations.assert_awaited_once_with(1)

def test_create_conversation_rejects_same_user(authenticated_client):
    mock_service = authenticated_client

    mock_service.create_conversation.side_effect = ValueError(
        "Users must be different"
    )

    response = client.post(
        "/api/chat/conversations/",
        json={"user2_id": 1},
        headers={"Authorization": "Bearer fake-token"},
    )

    assert response.status_code == 400

    assert response.json() == {
        "detail": "Users must be different"
    }

    mock_service.create_conversation.assert_awaited_once_with(
        user1_id=1,
        user2_id=1,
    )

def test_get_conversation(authenticated_client):
    mock_service = authenticated_client

    conversation = Conversation(
        id=10,
        user1_id=1,
        user2_id=2,
        created_at=datetime.now(),
    )

    mock_service.get_conversation.return_value = conversation

    response = client.get(
        "/api/chat/conversations/10",
        headers={"Authorization": "Bearer fake-token"},
    )

    assert response.status_code == 200

    data = response.json()

    assert data["id"] == 10
    assert data["user1_id"] == 1
    assert data["user2_id"] == 2

    mock_service.get_conversation.assert_awaited_once_with(
        conversation_id=10,
        user_id=1,
    )

def test_get_conversation_forbidden_for_non_member(authenticated_client):
    mock_service = authenticated_client

    mock_service.get_conversation.return_value = None

    response = client.get(
        "/api/chat/conversations/10",
        headers={"Authorization": "Bearer fake-token"},
    )

    assert response.status_code == 404

    assert response.json() == {
        "detail": "Conversation not found"
    }

    mock_service.get_conversation.assert_awaited_once_with(
        conversation_id=10,
        user_id=1,
    )

def test_create_message(authenticated_client):
    mock_service = authenticated_client

    message = Message(
        id=20,
        conversation_id=10,
        sender_id=1,
        content="Hello!",
        created_at=datetime.now(),
    )

    mock_service.create_message.return_value = message

    response = client.post(
        "/api/chat/conversations/10/messages",
        json={"content": "Hello!"},
        headers={"Authorization": "Bearer fake-token"},
    )

    assert response.status_code == 201

    data = response.json()

    assert data["id"] == 20
    assert data["conversation_id"] == 10
    assert data["sender_id"] == 1
    assert data["content"] == "Hello!"

    mock_service.create_message.assert_awaited_once_with(
        conversation_id=10,
        sender_id=1,
        content="Hello!",
    )

def test_create_message_rejects_unknown_conversation(authenticated_client):
    mock_service = authenticated_client

    mock_service.create_message.side_effect = ValueError(
        "Conversation not found"
    )

    response = client.post(
        "/api/chat/conversations/999/messages",
        json={"content": "Hello!"},
        headers={"Authorization": "Bearer fake-token"},
    )

    assert response.status_code == 404

    assert response.json() == {
        "detail": "Conversation not found"
    }

def test_list_messages(authenticated_client):
    mock_service = authenticated_client

    messages = [
        Message(
            id=1,
            conversation_id=10,
            sender_id=1,
            content="Hello!",
            created_at=datetime.now(),
        ),
        Message(
            id=2,
            conversation_id=10,
            sender_id=2,
            content="Hi!",
            created_at=datetime.now(),
        ),
    ]

    mock_service.list_messages.return_value = messages

    response = client.get(
        "/api/chat/conversations/10/messages",
        headers={"Authorization": "Bearer fake-token"},
    )

    assert response.status_code == 200

    data = response.json()

    assert len(data) == 2

    assert data[0]["id"] == 1
    assert data[0]["conversation_id"] == 10
    assert data[0]["sender_id"] == 1
    assert data[0]["content"] == "Hello!"

    assert data[1]["id"] == 2
    assert data[1]["conversation_id"] == 10
    assert data[1]["sender_id"] == 2
    assert data[1]["content"] == "Hi!"

    mock_service.list_messages.assert_awaited_once_with(
        conversation_id=10,
        user_id=1,
    )

def test_list_messages_for_unknown_conversation(authenticated_client):
    mock_service = authenticated_client

    mock_service.list_messages.side_effect = ValueError(
        "Conversation not found"
    )

    response = client.get(
        "/api/chat/conversations/999/messages",
        headers={"Authorization": "Bearer fake-token"},
    )

    assert response.status_code == 404

    assert response.json() == {
        "detail": "Conversation not found"
    }

    mock_service.list_messages.assert_awaited_once_with(
        conversation_id=999,
        user_id=1,
    )
