import { api } from "./client";

export interface UserPreview {
  id: number;
  email: string;
}

export interface Conversation {
  id: number;
  user1_id: number;
  user2_id: number;
  other_user: UserPreview;
  created_at: string;
}

export interface Message {
  id: number;
  conversation_id: number;
  sender_id: number;
  content: string;
  created_at: string;
}

export interface CreateConversationRequest {
  user2_id: number;
}

export interface CreateMessageRequest {
  content: string;
}

export async function getConversations(): Promise<Conversation[]> {
  const response = await api.get<Conversation[]>("/chat/conversations/");

  return response.data;
}

export async function createConversation(
  data: CreateConversationRequest,
): Promise<Conversation> {
  const response = await api.post<Conversation>("/chat/conversations/", data);

  return response.data;
}

export async function searchUserByEmail(email: string): Promise<UserPreview> {
  const response = await api.get<UserPreview>("/chat/users/search", {
    params: { email },
  });

  return response.data;
}

export async function getConversation(
  conversationId: number,
): Promise<Conversation> {
  const response = await api.get<Conversation>(
    `/chat/conversations/${conversationId}`,
  );

  return response.data;
}

export async function getMessages(conversationId: number): Promise<Message[]> {
  const response = await api.get<Message[]>(
    `/chat/conversations/${conversationId}/messages`,
  );

  return response.data;
}

export async function createMessage(
  conversationId: number,
  data: CreateMessageRequest,
): Promise<Message> {
  const response = await api.post<Message>(
    `/chat/conversations/${conversationId}/messages`,
    data,
  );

  return response.data;
}
