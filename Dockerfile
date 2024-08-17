


from fastapi import Request, HTTPException
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
import httpx

security = HTTPBearer()

async def verify_token(token: str) -> dict:
    async with httpx.AsyncClient() as client:
        response = await client.post(
            "http://auth-service/verify-token",
            json={"token": token}
        )
        if response.status_code == 200:
            return response.json()
        else:
            raise HTTPException(status_code=401, detail="Invalid token")

async def get_current_user(credentials: HTTPAuthorizationCredentials = Depends(security)):
    token = credentials.credentials
    user_data = await verify_token(token)
    return user_data

# Use this as a dependency in your routes
# @app.get("/protected-route")
# async def protected_route(current_user: dict = Depends(get_current_user)):
#     return {"message": "This is a protected route", "user": current_user}
