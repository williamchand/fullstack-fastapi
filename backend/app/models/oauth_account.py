# app/models/oauth_account.py

import uuid
from typing import TYPE_CHECKING

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
    user: "User" = Relationship(back_populates="oauth_accounts")
    access_token: str | None = None
