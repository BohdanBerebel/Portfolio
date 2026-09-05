import type { Message } from "../api/chat";
import { getToken } from "../api/authStorage";

export type MessageHandler = (message: Message) => void;

export type ErrorHandler = () => void;

export class ChatSocket {
  private socket: WebSocket | null = null;

  connect(
    conversationId: number,
    onMessage: MessageHandler,
    onClose?: () => void,
    onOpen?: () => void,
    onError?: ErrorHandler,
  ): void {
    const token = getToken();

    if (!token) {
      onError?.();
      return;
    }

    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";

    const url =
      `${protocol}//${window.location.host}` +
      `/api/chat/ws/${conversationId}?token=${encodeURIComponent(token)}`;

    this.socket = new WebSocket(url);

    this.socket.onopen = () => {
      onOpen?.();
    };

    this.socket.onmessage = (event) => {
      const message: Message = JSON.parse(event.data);
      onMessage(message);
    };

    this.socket.onerror = () => {
      onError?.();
    };

    this.socket.onclose = () => {
      this.socket = null;
      onClose?.();
    };
  }

  send(content: string): void {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      throw new Error("WebSocket is not connected");
    }

    this.socket.send(content);
  }

  disconnect(): void {
    this.socket?.close();
    this.socket = null;
  }
}
