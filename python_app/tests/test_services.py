from datetime import datetime
from unittest.mock import AsyncMock

import pytest

from database.models import Conversation, Message
from repositories.user import User
from services.chat import ChatService


@pytest.fixture
def repositories():
    conversation_repository = AsyncMock()
    message_repository = AsyncMock()
    user_repository = AsyncMock()

    return (
        conversation_repository,
        message_repository,
        user_repository,
    )

@pytest.fixture
def service(repositories):
    (
        conversation_repository,
        message_repository,
        user_repository,
    ) = repositories

    return ChatService(
        conversation_repository=conversation_repository,
        message_repository=message_repository,
        user_repository=user_repository,
    )


@pytest.mark.asyncio
async def test_create_conversation_sorts_user_ids(
    service,
    repositories,
):
    conversation_repository, _, _ = repositories

    conversation = Conversation(
        id=1,
        user1_id=1,
        user2_id=3,
        created_at=datetime.now(),
    )

    conversation_repository.create.return_value = conversation

    result = await service.create_conversation(
        user1_id=3,
        user2_id=1,
    )

    assert result == conversation

    conversation_repository.create.assert_awaited_once_with(1, 3)


@pytest.mark.asyncio
async def test_create_conversation_rejects_same_user(
    service,
    repositories,
):
    conversation_repository, _, _ = repositories

    with pytest.raises(ValueError, match="Users must be different"):
        await service.create_conversation(
            user1_id=1,
            user2_id=1,
        )

    conversation_repository.create.assert_not_awaited()


@pytest.mark.asyncio
async def test_create_message(
    service,
    repositories,
):
    conversation_repository, message_repository, _ = repositories

    conversation = Conversation(
        id=1,
        user1_id=1,
        user2_id=2,
        created_at=datetime.now(),
    )

    message = Message(
        id=1,
        conversation_id=1,
        sender_id=1,
        content="Hello!",
        created_at=datetime.now(),
    )

    conversation_repository.get_for_user.return_value = conversation
    message_repository.create.return_value = message

    result = await service.create_message(
        conversation_id=1,
        sender_id=1,
        content="  Hello!  ",
    )

    assert result == message

    message_repository.create.assert_awaited_once_with(
        1,
        1,
        "Hello!",
    )


@pytest.mark.asyncio
async def test_create_message_rejects_unknown_conversation(
    service,
    repositories,
):
    conversation_repository, message_repository, _ = repositories

    conversation_repository.get_for_user.return_value = None

    with pytest.raises(ValueError, match="Conversation not found"):
        await service.create_message(
            conversation_id=999,
            sender_id=1,
            content="Hello!",
        )

    message_repository.create.assert_not_awaited()


@pytest.mark.asyncio
async def test_create_message_rejects_empty_content(
    service,
    repositories,
):
    conversation_repository, message_repository, _ = repositories

    conversation_repository.get_for_user.return_value = Conversation(
        id=1,
        user1_id=1,
        user2_id=2,
        created_at=datetime.now(),
    )

    with pytest.raises(
        ValueError,
        match="Message content cannot be empty",
    ):
        await service.create_message(
            conversation_id=1,
            sender_id=1,
            content="   ",
        )

    message_repository.create.assert_not_awaited()

@pytest.mark.asyncio
async def test_get_user_by_email(
    service,
    repositories,
):
    _, _, user_repository = repositories

    user = User(
        id=2,
        email="admin@tubely.com",
    )

    user_repository.get_by_email.return_value = user

    result = await service.get_user_by_email(
        "admin@tubely.com",
    )

    assert result == user

    user_repository.get_by_email.assert_awaited_once_with(
        "admin@tubely.com",
    )