import { useState } from "react";
import type { SubmitEvent } from "react";
import { createNote } from "../api/notes";

interface CreateNoteFormProps {
  onCreated: () => void;
}

function CreateNoteForm({ onCreated }: CreateNoteFormProps) {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    setError("");
    setLoading(true);

    try {
      await createNote({
        title,
        content,
      });

      setTitle("");
      setContent("");

      onCreated();
    } catch {
      setError("Failed to create note.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <h2>Create Note</h2>

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="note-title">Title</label>

          <input
            id="note-title"
            type="text"
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="note-content">Content</label>

          <textarea
            id="note-content"
            value={content}
            onChange={(event) => setContent(event.target.value)}
            required
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "Creating..." : "Create Note"}
        </button>
      </form>

      {error && <p className="error-message">{error}</p>}
    </div>
  );
}

export default CreateNoteForm;
