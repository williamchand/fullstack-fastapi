import uuid

from sqlmodel import Field

from .base import BaseModel


class UserRole(BaseModel, table=True):
    __tablename__ = "user_role"
    user_id: uuid.UUID | None = Field(foreign_key="user.id", primary_key=True)
    role_id: int | None = Field(foreign_key="role.id", primary_key=True)

