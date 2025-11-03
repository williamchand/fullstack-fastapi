# app/models/oauth_account.py
import uuid

from sqlmodel import Field, Relationship

from .base import BaseModel
from .user import User


class OAuthAccount(BaseModel):
    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    provider: str = Field(index=True)  # e.g., "google"
    provider_account_id: str = Field(index=True)  # Google sub ID
    email: str = Field(index=True)
    user_id: uuid.UUID = Field(foreign_key="user.id")

    user: User = Relationship(back_populates="oauth_accounts")
