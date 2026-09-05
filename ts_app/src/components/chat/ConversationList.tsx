import {
  Box,
  List,
  ListItemButton,
  ListItemText,
  Typography,
} from "@mui/material";

import type { Conversation } from "../../api/chat";

interface ConversationListProps {
  conversations: Conversation[];
  selectedConversationId: number | null;
  onSelect: (conversation: Conversation) => void;
  isLoading?: boolean;
  error?: string | null;
}

export default function ConversationList({
  conversations,
  selectedConversationId,
  onSelect,
  isLoading = false,
  error = null,
}: ConversationListProps) {
  if (isLoading) {
    return (
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          p: 3,
        }}
      >
        Loading...
      </Box>
    );
  }

  if (error) {
    return (
      <Box sx={{ p: 2 }}>
        <Typography color="error">{error}</Typography>
      </Box>
    );
  }

  if (conversations.length === 0) {
    return (
      <Box sx={{ p: 2 }}>
        <Typography color="text.secondary">No conversations yet</Typography>
      </Box>
    );
  }

  return (
    <List disablePadding>
      {conversations.map((conversation) => (
        <ListItemButton
          key={conversation.id}
          selected={conversation.id === selectedConversationId}
          onClick={() => onSelect(conversation)}
        >
          <ListItemText
            primary={conversation.other_user.email}
            secondary={`Conversation #${conversation.id}`}
          />
        </ListItemButton>
      ))}
    </List>
  );
}
