import { api } from "./client";

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  id: number;
  email: string;
  token: string;
}

export async function register(data: RegisterRequest): Promise<AuthResponse> {
  const response = await api.post<{ data: AuthResponse }>(
    "/auth/register",
    data,
  );

  return response.data.data;
}

export async function login(data: LoginRequest): Promise<AuthResponse> {
  const response = await api.post<{ data: AuthResponse }>("/auth/login", data);

  return response.data.data;
}
