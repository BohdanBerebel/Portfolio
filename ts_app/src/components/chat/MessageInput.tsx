import { useState } from "react";
import { Box, Button, TextField } from "@mui/material";

interface MessageInputProps {
  onSend: (content: string) => void;
  disabled?: boolean;
}

export default function MessageInput({
  onSend,
  disabled = false,
}: MessageInputProps) {
  const [content, setContent] = useState("");

  function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const trimmedContent = content.trim();

    if (!trimmedContent || disabled) {
      return;
    }

    onSend(trimmedContent);
    setContent("");
  }

  return (
    <Box
      component="form"
      onSubmit={handleSubmit}
      sx={{
        display: "flex",
        gap: 1,
      }}
    >
      <TextField
        fullWidth
        size="small"
        placeholder="Type a message..."
        value={content}
        onChange={(event) => setContent(event.target.value)}
        disabled={disabled}
      />

      <Button
        type="submit"
        variant="contained"
        disabled={disabled || !content.trim()}
      >
        Send
      </Button>
    </Box>
  );
}
