from fastapi import FastAPI
from controller.root import router as root_router
from controller.probability import router as probability_router
import uvicorn

app = FastAPI()
app.include_router(root_router)
app.include_router(probability_router)

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8090)