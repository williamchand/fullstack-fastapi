# app/models/pending_user.py
import uuid
from datetime import datetime
from sqlmodel import Field
from .role import UserRole
from .base import BaseModel

class PendingUser(BaseModel):
    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    email: str | None = Field(default=None, index=True, nullable=True)
    phone_number: str | None = Field(default=None, unique=True, nullable=True)
    hashed_password: str | None = None
    verification_code: str = Field()
    role: UserRole = Field(default=UserRole.CUSTOMER)
    created_at: datetime = Field(default_factory=datetime.utcnow)
    expires_at: datetime = Field(default_factory=lambda: datetime.utcnow())
