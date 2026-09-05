import { api } from "./client";

export interface Note {
  id: number;
  user_id: number;
  title: string;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface CreateNoteRequest {
  title: string;
  content: string;
}

export interface UpdateNoteRequest {
  title: string;
  content: string;
}

export async function getNotes(): Promise<Note[]> {
  const response = await api.get<{ data: Note[] }>("/notes");

  return response.data.data;
}

export async function getNote(id: number): Promise<Note> {
  const response = await api.get<{ data: Note }>(`/notes/${id}`);

  return response.data.data;
}

export async function createNote(data: CreateNoteRequest): Promise<Note> {
  const response = await api.post<{ data: Note }>("/notes", data);

  return response.data.data;
}

export async function updateNote(
  id: number,
  data: UpdateNoteRequest,
): Promise<Note> {
  const response = await api.put<{ data: Note }>(`/notes/${id}`, data);

  return response.data.data;
}

export async function deleteNote(id: number): Promise<void> {
  await api.delete(`/notes/${id}`);
}
