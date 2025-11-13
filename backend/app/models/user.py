# app/models/user.py
from __future__ import annotations

import uuid
from typing import TYPE_CHECKING, List

from sqlalchemy.orm import Mapped, relationship
from sqlmodel import Field, Relationship

from .base import BaseModel
from .user_role import UserRole

if TYPE_CHECKING:  # 👈 avoids circular import at runtime
    from .oauth_account import OAuthAccount
    from .role import Role


# Shared properties
class UserBase(BaseModel):
    email: str | None = Field(default=None, index=True, nullable=True, unique=True)
    phone_number: str | None = Field(default=None, index=True, nullable=True, unique=True)
    is_active: bool = True
    is_superuser: bool = False
    is_email_verified: bool = Field(default=False)
    is_phone_verified: bool = Field(default=False)
    full_name: str | None = Field(default=None, max_length=255)


# Database model, database table inferred from class name
class User(UserBase, table=True):
    id: uuid.UUID | None = Field(default=None, primary_key=True)
    hashed_password: str
    roles: Mapped[list[Role]] = Relationship(
        sa_relationship=relationship(
            "Role",
            secondary=UserRole.__table__,
            back_populates="users"
        )
    )
    oauth_accounts: Mapped[list[OAuthAccount]] = Relationship(sa_relationship=relationship(back_populates="user"))
