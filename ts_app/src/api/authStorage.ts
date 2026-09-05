const TOKEN_KEY = "auth_token";
const USER_ID_KEY = "user_id";

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

export function getUserId(): number | null {
  const value = localStorage.getItem(USER_ID_KEY);

  if (value === null) {
    return null;
  }

  const userId = Number(value);

  return Number.isNaN(userId) ? null : userId;
}

export function setUserId(userId: number): void {
  localStorage.setItem(USER_ID_KEY, String(userId));
}

export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(USER_ID_KEY);
}
