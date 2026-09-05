import { Box, Paper, Typography } from "@mui/material";
import type { Message } from "../../api/chat";

interface MessageItemProps {
  message: Message;
  isOwn: boolean;
}

export default function MessageItem({ message, isOwn }: MessageItemProps) {
  const time = new Date(message.created_at).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: isOwn ? "flex-end" : "flex-start",
        mb: 1.5,
      }}
    >
      <Paper
        elevation={1}
        sx={{
          maxWidth: "70%",
          px: 2,
          py: 1,
          borderRadius: 2,
          bgcolor: isOwn ? "primary.main" : "background.paper",
          color: isOwn ? "primary.contrastText" : "text.primary",
        }}
      >
        {!isOwn && (
          <Typography
            variant="caption"
            color="text.secondary"
            sx={{
              display: "block",
              mb: 0.25,
            }}
          >
            User {message.sender_id}
          </Typography>
        )}

        <Typography
          sx={{
            overflowWrap: "anywhere",
            whiteSpace: "pre-wrap",
          }}
        >
          {message.content}
        </Typography>

        <Typography
          variant="caption"
          sx={{
            display: "block",
            textAlign: "right",
            mt: 0.5,
            color: isOwn ? "inherit" : "text.secondary",
            opacity: isOwn ? 0.8 : 1,
          }}
        >
          {time}
        </Typography>
      </Paper>
    </Box>
  );
}
