import jwt

from config.settings import settings

def validate_token(token: str) -> int:
    payload = jwt.decode(
        token,
        settings.jwt_secret,
        algorithms=["HS256"],
    )

    user_id = payload.get("user_id")

    if user_id is None:
        raise ValueError("Token does not contain user_id")

    return int(user_id)