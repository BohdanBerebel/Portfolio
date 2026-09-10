from contextlib import asynccontextmanager
import logging

from fastapi import FastAPI

from config.settings import settings
from database.connection import create_pool
from logging_config import setup_logging

from routes.conversations import router as conversations_router
from routes.websocket import router as websocket_router
from routes.users import router as users_router


setup_logging()
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info("starting chat service")

    app.state.db = await create_pool()
    logger.info("database connection pool created")

    yield

    logger.info("shutting down chat service")

    await app.state.db.close()
    logger.info("database connection pool closed")


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