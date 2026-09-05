import { useEffect, useState } from "react";
import { deleteNote, getNotes, updateNote } from "../api/notes";
import type { Note } from "../api/notes";

interface NotesListProps {
  refreshKey: number;
}

function NotesList({ refreshKey }: NotesListProps) {
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [editingId, setEditingId] = useState<number | null>(null);
  const [editTitle, setEditTitle] = useState("");
  const [editContent, setEditContent] = useState("");
  const [actionLoading, setActionLoading] = useState(false);

  useEffect(() => {
    getNotes()
      .then((result) => {
        setNotes(result);
        setError("");
      })
      .catch(() => {
        setError("Failed to load notes.");
      })
      .finally(() => {
        setLoading(false);
      });
  }, [refreshKey]);

  function startEditing(note: Note) {
    setEditingId(note.id);
    setEditTitle(note.title);
    setEditContent(note.content);
    setError("");
  }

  function cancelEditing() {
    setEditingId(null);
    setEditTitle("");
    setEditContent("");
  }

  async function handleUpdate(id: number) {
    setActionLoading(true);
    setError("");

    try {
      const updatedNote = await updateNote(id, {
        title: editTitle,
        content: editContent,
      });

      setNotes((currentNotes) =>
        currentNotes.map((note) => (note.id === id ? updatedNote : note)),
      );

      cancelEditing();
    } catch {
      setError("Failed to update note.");
    } finally {
      setActionLoading(false);
    }
  }

  async function handleDelete(id: number) {
    const confirmed = window.confirm(
      "Are you sure you want to delete this note?",
    );

    if (!confirmed) {
      return;
    }

    setActionLoading(true);
    setError("");

    try {
      await deleteNote(id);

      setNotes((currentNotes) => currentNotes.filter((note) => note.id !== id));
    } catch {
      setError("Failed to delete note.");
    } finally {
      setActionLoading(false);
    }
  }

  if (loading) {
    return <p>Loading notes...</p>;
  }

  if (error && notes.length === 0) {
    return <p className="error-message">{error}</p>;
  }

  return (
    <div>
      <div className="notes-header">
        <h2>My Notes</h2>
      </div>

      {error && <p className="error-message">{error}</p>}

      {notes.length === 0 ? (
        <p>No notes yet.</p>
      ) : (
        <div className="notes-grid">
          {notes.map((note) => (
            <article className="note-card" key={note.id}>
              {editingId === note.id ? (
                <div className="edit-form">
                  <div className="form-group">
                    <label htmlFor={`edit-title-${note.id}`}>Title</label>

                    <input
                      id={`edit-title-${note.id}`}
                      type="text"
                      value={editTitle}
                      onChange={(event) => setEditTitle(event.target.value)}
                    />
                  </div>

                  <div className="form-group">
                    <label htmlFor={`edit-content-${note.id}`}>Content</label>

                    <textarea
                      id={`edit-content-${note.id}`}
                      value={editContent}
                      onChange={(event) => setEditContent(event.target.value)}
                    />
                  </div>

                  <div className="note-actions">
                    <button
                      type="button"
                      onClick={() => handleUpdate(note.id)}
                      disabled={actionLoading}
                    >
                      {actionLoading ? "Saving..." : "Save"}
                    </button>

                    <button
                      className="secondary-button"
                      type="button"
                      onClick={cancelEditing}
                      disabled={actionLoading}
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              ) : (
                <>
                  <h3>{note.title}</h3>

                  <p>{note.content}</p>

                  <div className="note-actions">
                    <button
                      type="button"
                      onClick={() => startEditing(note)}
                      disabled={actionLoading}
                    >
                      Edit
                    </button>

                    <button
                      className="danger-button"
                      type="button"
                      onClick={() => handleDelete(note.id)}
                      disabled={actionLoading}
                    >
                      Delete
                    </button>
                  </div>
                </>
              )}
            </article>
          ))}
        </div>
      )}
    </div>
  );
}

export default NotesList;
