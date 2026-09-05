import { useState } from "react";
import { Box, Button, TextField } from "@mui/material";

import {
  createConversation,
  searchUserByEmail,
  type Conversation,
} from "../../api/chat";

interface NewConversationFormProps {
  onCreated: (conversation: Conversation) => void;
}

export default function NewConversationForm({
  onCreated,
}: NewConversationFormProps) {
  const [email, setEmail] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const trimmedEmail = email.trim();

    if (!trimmedEmail) {
      setError("Enter a user email");
      return;
    }

    try {
      setIsLoading(true);
      setError(null);

      const user = await searchUserByEmail(trimmedEmail);

      const conversation = await createConversation({
        user2_id: user.id,
      });

      onCreated(conversation);
      setEmail("");
    } catch {
      setError("Failed to find user or create conversation");
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <Box
      component="form"
      onSubmit={handleSubmit}
      sx={{
        display: "flex",
        gap: 1,
        p: 2,
      }}
    >
      <TextField
        fullWidth
        size="small"
        label="User email"
        type="email"
        value={email}
        onChange={(event) => setEmail(event.target.value)}
        disabled={isLoading}
      />

      <Button type="submit" variant="contained" disabled={isLoading}>
        {isLoading ? "Creating..." : "New conversation"}
      </Button>

      {error && (
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            color: "error.main",
          }}
        >
          {error}
        </Box>
      )}
    </Box>
  );
}
