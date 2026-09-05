from repositories.conversation import ConversationRepository
from repositories.message import MessageRepository
from repositories.user import UserRepository

class ChatService:
    def __init__(
        self,
        conversation_repository: ConversationRepository,
        message_repository: MessageRepository,
        user_repository: UserRepository,
    ):
        self.conversation_repository = conversation_repository
        self.message_repository = message_repository
        self.user_repository = user_repository

    async def get_user_by_email(
        self,
        email: str,
    ):
        return await self.user_repository.get_by_email(email)

    async def create_conversation(
        self,
        user1_id: int,
        user2_id: int,
    ):
        if user1_id == user2_id:
            raise ValueError("Users must be different")

        user1_id, user2_id = sorted((user1_id, user2_id))

        return await self.conversation_repository.create(
            user1_id,
            user2_id,
        )

    async def get_conversation(
        self,
        conversation_id: int,
        user_id: int,
    ):
        return await self.conversation_repository.get_for_user(
            conversation_id,
            user_id,
        )

    async def list_conversations(
        self,
        user_id: int,
    ):
        return await self.conversation_repository.list_for_user(
            user_id,
        )

    async def create_message(
        self,
        conversation_id: int,
        sender_id: int,
        content: str,
    ):
        conversation = await self.conversation_repository.get_for_user(
            conversation_id,
            sender_id,
        )

        if conversation is None:
            raise ValueError("Conversation not found")

        content = content.strip()

        if not content:
            raise ValueError("Message content cannot be empty")

        return await self.message_repository.create(
            conversation_id,
            sender_id,
            content,
        )

    async def list_messages(
        self,
        conversation_id: int,
        user_id: int,
    ):
        conversation = await self.conversation_repository.get_for_user(
            conversation_id,
            user_id,
        )

        if conversation is None:
            raise ValueError("Conversation not found")

        return await self.message_repository.list_by_conversation(
            conversation_id,
        )