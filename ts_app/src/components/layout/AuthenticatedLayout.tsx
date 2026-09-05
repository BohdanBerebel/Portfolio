import { useState } from "react";
import {
  Box,
  Divider,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Paper,
  Typography,
} from "@mui/material";

import {
  Chat as ChatIcon,
  Home as HomeIcon,
  Logout as LogoutIcon,
  Notes as NotesIcon,
} from "@mui/icons-material";

import ChatPage from "../../pages/ChatPage";
import CreateNoteForm from "../CreateNoteForm";
import NotesList from "../NotesList";
import { removeToken } from "../../api/authStorage";

type Page = "home" | "notes" | "chat";

interface AuthenticatedLayoutProps {
  onLogout: () => void;
}

export default function AuthenticatedLayout({
  onLogout,
}: AuthenticatedLayoutProps) {
  const [page, setPage] = useState<Page>("home");
  const [notesRefreshKey, setNotesRefreshKey] = useState(0);

  function handleLogout() {
    removeToken();
    onLogout();
  }

  function handleNoteCreated() {
    setNotesRefreshKey((key) => key + 1);
  }

  function renderPage() {
    switch (page) {
      case "notes":
        return (
          <Box sx={{ p: 3 }}>
            <Paper sx={{ p: 3, mb: 3 }}>
              <CreateNoteForm onCreated={handleNoteCreated} />
            </Paper>

            <Paper sx={{ p: 3 }}>
              <NotesList refreshKey={notesRefreshKey} />
            </Paper>
          </Box>
        );

      case "chat":
        return <ChatPage />;

      case "home":
      default:
        return (
          <Box
            sx={{
              height: "100%",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              p: 3,
            }}
          >
            <Box sx={{ textAlign: "center" }}>
              <Typography variant="h3" gutterBottom>
                Welcome back!
              </Typography>

              <Typography variant="h6" color="text.secondary">
                What would you like to do?
              </Typography>
            </Box>
          </Box>
        );
    }
  }

  return (
    <Box
      sx={{
        minHeight: "100vh",
        display: "flex",
        bgcolor: "background.default",
      }}
    >
      <Paper
        square
        elevation={0}
        sx={{
          width: 240,
          flexShrink: 0,
          display: "flex",
          flexDirection: "column",
          borderRight: 1,
          borderColor: "divider",
        }}
      >
        <Box sx={{ px: 3, py: 2.5 }}>
          <Typography
            variant="h6"
            sx={{
              fontWeight: 700,
            }}
          >
            Notes App
          </Typography>

          <Typography variant="body2" color="text.secondary">
            Personal workspace
          </Typography>
        </Box>

        <Divider />

        <List sx={{ px: 1, py: 1 }}>
          <ListItemButton
            selected={page === "home"}
            onClick={() => setPage("home")}
            sx={{
              borderRadius: 1.5,
              mb: 0.5,
            }}
          >
            <ListItemIcon>
              <HomeIcon />
            </ListItemIcon>
            <ListItemText primary="Home" />
          </ListItemButton>

          <ListItemButton
            selected={page === "notes"}
            onClick={() => setPage("notes")}
            sx={{
              borderRadius: 1.5,
              mb: 0.5,
            }}
          >
            <ListItemIcon>
              <NotesIcon />
            </ListItemIcon>
            <ListItemText primary="Notes" />
          </ListItemButton>

          <ListItemButton
            selected={page === "chat"}
            onClick={() => setPage("chat")}
            sx={{
              borderRadius: 1.5,
              mb: 0.5,
            }}
          >
            <ListItemIcon>
              <ChatIcon />
            </ListItemIcon>
            <ListItemText primary="Chat" />
          </ListItemButton>
        </List>

        <Box sx={{ mt: "auto", p: 2 }}>
          <ListItemButton
            onClick={handleLogout}
            sx={{
              borderRadius: 1.5,
            }}
          >
            <ListItemIcon>
              <LogoutIcon />
            </ListItemIcon>
            <ListItemText primary="Logout" />
          </ListItemButton>
        </Box>
      </Paper>

      <Box
        component="main"
        sx={{
          flex: 1,
          minWidth: 0,
        }}
      >
        {renderPage()}
      </Box>
    </Box>
  );
}
