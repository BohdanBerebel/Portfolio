from fastapi import APIRouter, Depends, HTTPException, Query

from dependencies import get_chat_service, get_current_user_id
from services.chat import ChatService


router = APIRouter(
    prefix="/api/chat/users",
    tags=["users"],
)


@router.get("/search")
async def search_user(
    email: str = Query(...),
    user_id: int = Depends(get_current_user_id),
    chat_service: ChatService = Depends(get_chat_service),
):
    user = await chat_service.get_user_by_email(email)

    if user is None:
        raise HTTPException(
            status_code=404,
            detail="User not found",
        )

    if user.id == user_id:
        raise HTTPException(
            status_code=400,
            detail="Cannot create conversation with yourself",
        )

    return {
        "id": user.id,
        "email": user.email,
    }