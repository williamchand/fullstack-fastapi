# app/models/base.py
from sqlmodel import SQLModel


class BaseModel(SQLModel):
    """
    Shared configuration for all entities
    (can be extended later for timestamps or common methods).
    """
    pass

# Generic message
class Message(BaseModel):
    message: str
