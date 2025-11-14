# app/models/base.py
import uuid

from sqlalchemy import text
from sqlmodel import Field, SQLModel


class BaseModel(SQLModel):
    """
    Shared configuration for all entities
    (can be extended later for timestamps or common methods).
    """
    pass

class BaseModelUUID(SQLModel, table=False):
    id: uuid.UUID | None = Field(
        default=None,
        primary_key=True,
        nullable=False,
        sa_column_kwargs={"server_default": text("gen_random_uuid()")},
    )

# Generic message
class Message(BaseModel):
    message: str
