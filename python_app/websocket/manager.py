from fastapi import WebSocket


class ConnectionManager:
    def __init__(self):
        self.connections: dict[int, set[WebSocket]] = {}

    async def connect(
        self,
        user_id: int,
        websocket: WebSocket,
    ):
        await websocket.accept()

        if user_id not in self.connections:
            self.connections[user_id] = set()

        self.connections[user_id].add(websocket)

    def disconnect(
        self,
        user_id: int,
        websocket: WebSocket,
    ):
        connections = self.connections.get(user_id)

        if connections is None:
            return

        connections.discard(websocket)

        if not connections:
            del self.connections[user_id]

    async def send_to_user(
        self,
        user_id: int,
        message: dict,
    ):
        connections = self.connections.get(user_id)

        if not connections:
            return

        disconnected = []

        for websocket in connections:
            try:
                await websocket.send_json(message)
            except Exception:
                disconnected.append(websocket)

        for websocket in disconnected:
            connections.discard(websocket)

        if not connections:
            del self.connections[user_id]