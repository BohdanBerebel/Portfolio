from contextlib import asynccontextmanager

from fastapi import FastAPI

from config.settings import settings
from database.connection import create_pool

from routes.conversations import router as conversations_router
from routes.websocket import router as websocket_router
from routes.users import router as users_router

@asynccontextmanager
async def lifespan(app: FastAPI):
    app.state.db = await create_pool()

    yield

    await app.state.db.close()

app = FastAPI(
    title="Chat Service",
    version="1.0.0",
    lifespan=lifespan,
)

app.include_router(conversations_router)
app.include_router(websocket_router)
app.include_router(users_router)

@app.get("/health")
async def health():
    return {"status": "ok"}

@app.get("/ready")
async def ready():
    await app.state.db.fetchval("SELECT 1")

    return {"status": "ready"}