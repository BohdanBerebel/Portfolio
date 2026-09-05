import { useEffect, useRef, useState } from "react";
import { Box, Paper, Typography } from "@mui/material";

import ConversationList from "../components/chat/ConversationList";
import MessageInput from "../components/chat/MessageInput";
import MessageItem from "../components/chat/MessageItem";
import { useChat } from "../hooks/useChat";
import { getUserId } from "../api/authStorage";
import { getConversations, type Conversation } from "../api/chat";
import NewConversationForm from "../components/chat/NewConversationForm";

export default function ChatPage() {
  const [selectedConversation, setSelectedConversation] =
    useState<Conversation | null>(null);

  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [conversationsLoading, setConversationsLoading] = useState(true);
  const [conversationsError, setConversationsError] = useState<string | null>(
    null,
  );

  const { messages, isLoading, connectionStatus, error, sendMessage } = useChat(
    selectedConversation?.id ?? null,
  );

  const currentUserId = getUserId();

  const messagesEndRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    async function loadConversations() {
      try {
        setConversationsLoading(true);
        setConversationsError(null);

        const data = await getConversations();

        setConversations(data);
      } catch {
        setConversationsError("Failed to load conversations");
      } finally {
        setConversationsLoading(false);
      }
    }

    void loadConversations();
  }, []);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({
      behavior: "smooth",
    });
  }, [messages]);

  function handleConversationCreated(conversation: Conversation) {
    setConversations((current) => {
      const exists = current.some((item) => item.id === conversation.id);

      if (exists) {
        return current;
      }

      return [conversation, ...current];
    });

    setSelectedConversation(conversation);
  }

  return (
    <Box
      sx={{
        height: "70vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        p: 2,
        bgcolor: "background.default",
      }}
    >
      <Paper
        elevation={3}
        sx={{
          width: "100%",
          maxWidth: 1200,
          height: "100%",
          display: "flex",
          overflow: "hidden",

          flexDirection: {
            xs: "column",
            md: "row",
          },
        }}
      >
        {/* Conversations */}
        <Box
          sx={{
            width: {
              xs: "100%",
              md: 320,
            },
            height: {
              xs: 180,
              md: "100%",
            },
            borderRight: {
              xs: 0,
              md: 1,
            },
            borderBottom: {
              xs: 1,
              md: 0,
            },
            borderColor: "divider",
            display: "flex",
            flexDirection: "column",
            flexShrink: 0,
          }}
        >
          <Box
            sx={{
              p: 2,
              borderBottom: 1,
              borderColor: "divider",
            }}
          >
            <Typography variant="h6">Conversations</Typography>
          </Box>

          <NewConversationForm onCreated={handleConversationCreated} />

          <Box sx={{ flex: 1, overflow: "auto" }}>
            <ConversationList
              conversations={conversations}
              selectedConversationId={selectedConversation?.id ?? null}
              onSelect={setSelectedConversation}
              isLoading={conversationsLoading}
              error={conversationsError}
            />
          </Box>
        </Box>

        {/* Chat */}
        <Box
          sx={{
            flex: 1,
            display: "flex",
            flexDirection: "column",
            minWidth: 0,
          }}
        >
          <Box
            sx={{
              p: 2,
              borderBottom: 1,
              borderColor: "divider",
            }}
          >
            <Typography variant="h6">
              {selectedConversation === null
                ? "Chat"
                : selectedConversation.other_user.email}
            </Typography>

            {selectedConversation !== null && (
              <Typography
                variant="body2"
                color={
                  connectionStatus === "connected"
                    ? "success.main"
                    : "text.secondary"
                }
              >
                {connectionStatus === "connected"
                  ? "Connected"
                  : connectionStatus === "connecting"
                    ? "Connecting..."
                    : "Disconnected"}
              </Typography>
            )}
          </Box>

          <Box
            sx={{
              flex: 1,
              overflow: "auto",
              p: 2,
            }}
          >
            {selectedConversation === null ? (
              <Box
                sx={{
                  height: "100%",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                }}
              >
                <Typography color="text.secondary">
                  Select a conversation
                </Typography>
              </Box>
            ) : isLoading ? (
              <Typography color="text.secondary">
                Loading messages...
              </Typography>
            ) : error ? (
              <Typography color="error">{error}</Typography>
            ) : messages.length === 0 ? (
              <Box
                sx={{
                  height: "100%",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                }}
              >
                <Typography color="text.secondary">No messages yet</Typography>
              </Box>
            ) : (
              <>
                {messages.map((message) => (
                  <MessageItem
                    key={message.id}
                    message={message}
                    isOwn={message.sender_id === currentUserId}
                  />
                ))}

                <div ref={messagesEndRef} />
              </>
            )}
          </Box>

          {/* Message input буде наступним кроком */}
          {selectedConversation !== null && (
            <Box
              sx={{
                p: 2,
                borderTop: 1,
                borderColor: "divider",
              }}
            >
              <MessageInput
                onSend={sendMessage}
                disabled={connectionStatus !== "connected"}
              />
            </Box>
          )}
        </Box>
      </Paper>
    </Box>
  );
}
