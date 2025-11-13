# app/models/oauth_account.py
from __future__ import annotations

import uuid
from typing import TYPE_CHECKING

from sqlalchemy.orm import Mapped, relationship
from sqlmodel import Field, Relationship

from .base import BaseModel

if TYPE_CHECKING:  # 👈 avoids circular import at runtime
    from .user import User


class OAuthAccount(BaseModel, table=True):
    __tablename__ = "oauth_account"
    id: uuid.UUID | None = Field(default=None, primary_key=True)
    user_id: uuid.UUID = Field(foreign_key="user.id")
    provider: str = Field(index=True, max_length=50)  # e.g., "google"
    provider_user_id: str = Field(max_length=255, index=True)  # Google sub ID
    user: Mapped[User] = Relationship(sa_relationship=relationship(back_populates="oauth_accounts"))
    access_token: str | None = None
