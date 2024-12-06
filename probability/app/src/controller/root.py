from fastapi import APIRouter

router = APIRouter()

@router.get("/")
def Hello():
    return {"Hello": "World"}