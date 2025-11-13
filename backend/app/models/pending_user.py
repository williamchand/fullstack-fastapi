# app/models/pending_user.py
from __future__ import annotations

import uuid
from datetime import datetime

from sqlmodel import Field

from .base import BaseModel


class PendingUser(BaseModel, table=True):
    id: uuid.UUID | None = Field(default=None, primary_key=True)
    email: str | None = Field(default=None, index=True, nullable=True, max_length=255)
    phone_number: str | None = Field(default=None, unique=True, nullable=True, max_length=32)
    full_name: str | None = Field(default=None, max_length=255)
    hashed_password: str | None = None
    verification_code: str = Field(max_length=64)
    verification_type: str = Field(max_length=16)
    created_at: datetime = Field()
    expires_at: datetime = Field()
