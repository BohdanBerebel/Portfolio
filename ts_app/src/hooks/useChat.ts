import { useCallback, useEffect, useRef, useState } from "react";

import { getMessages, type Message } from "../api/chat";

import { ChatSocket } from "../websocket/chatSocket";

type MessageState = {
  conversationId: number | null;
  messages: Message[];
};

type ConnectionStatus = "disconnected" | "connecting" | "connected";

type ConnectionState = {
  conversationId: number | null;
  status: ConnectionStatus;
};

type LoadingState = {
  conversationId: number | null;
  loading: boolean;
};

type ErrorState = {
  conversationId: number | null;
  error: string | null;
};

export function useChat(conversationId: number | null) {
  const [messageState, setMessageState] = useState<MessageState>({
    conversationId: null,
    messages: [],
  });

  const [connectionState, setConnectionState] = useState<ConnectionState>({
    conversationId: null,
    status: "disconnected",
  });

  const [loadingState, setLoadingState] = useState<LoadingState>({
    conversationId: null,
    loading: false,
  });

  const [errorState, setErrorState] = useState<ErrorState>({
    conversationId: null,
    error: null,
  });

  const socketRef = useRef<ChatSocket | null>(null);

  useEffect(() => {
    if (conversationId === null) {
      socketRef.current?.disconnect();
      socketRef.current = null;
      return;
    }

    const currentConversationId = conversationId;
    let cancelled = false;

    const socket = new ChatSocket();

    socketRef.current = socket;

    async function loadMessages() {
      setLoadingState({
        conversationId: currentConversationId,
        loading: true,
      });

      setErrorState({
        conversationId: currentConversationId,
        error: null,
      });

      try {
        const data = await getMessages(currentConversationId);

        if (cancelled) {
          return;
        }

        setMessageState({
          conversationId: currentConversationId,
          messages: data,
        });
      } catch {
        if (cancelled) {
          return;
        }

        setErrorState({
          conversationId: currentConversationId,
          error: "Failed to load messages",
        });
      } finally {
        if (!cancelled) {
          setLoadingState({
            conversationId: currentConversationId,
            loading: false,
          });
        }
      }
    }

    void loadMessages();

    setConnectionState({
      conversationId: currentConversationId,
      status: "connecting",
    });

    socket.connect(
      currentConversationId,
      (message) => {
        if (cancelled) {
          return;
        }

        setMessageState((current) => ({
          conversationId: currentConversationId,
          messages:
            current.conversationId === currentConversationId
              ? [...current.messages, message]
              : [message],
        }));
      },
      () => {
        if (cancelled) {
          return;
        }

        setConnectionState({
          conversationId: currentConversationId,
          status: "disconnected",
        });
      },
      () => {
        if (cancelled) {
          return;
        }

        setConnectionState({
          conversationId: currentConversationId,
          status: "connected",
        });
      },
      () => {
        if (cancelled) {
          return;
        }

        setErrorState({
          conversationId: currentConversationId,
          error: "Failed to connect to chat",
        });
      },
    );

    return () => {
      cancelled = true;
      socket.disconnect();
    };
  }, [conversationId]);

  const messages =
    messageState.conversationId === conversationId ? messageState.messages : [];

  const isLoading =
    loadingState.conversationId === conversationId && loadingState.loading;

  const connectionStatus =
    connectionState.conversationId === conversationId
      ? connectionState.status
      : "disconnected";

  const error =
    errorState.conversationId === conversationId ? errorState.error : null;

  const sendMessage = useCallback((content: string) => {
    const trimmedContent = content.trim();

    if (!trimmedContent) {
      return;
    }

    socketRef.current?.send(trimmedContent);
  }, []);

  return {
    messages,
    isLoading,
    connectionStatus,
    error,
    sendMessage,
  };
}
