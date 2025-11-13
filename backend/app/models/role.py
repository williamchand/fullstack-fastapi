# app/models/role.py
from __future__ import annotations

from enum import Enum
from typing import TYPE_CHECKING, List

from sqlalchemy.orm import Mapped, relationship
from sqlmodel import Field, Relationship

from .base import BaseModel
from .user_role import UserRole

if TYPE_CHECKING:  # 👈 avoids circular import at runtime
    from .user import User


# Shared properties
class RoleBase(BaseModel):
    name: str = Field(index=True, max_length=50)
    description: str| None = Field(default=None, max_length=255)

class Role(RoleBase, table=True):
    id: int | None = Field(default=None, primary_key=True)
    users: Mapped[list[User]] = Relationship(
        sa_relationship=relationship(
            "User",
            secondary=UserRole.__table__,
            back_populates="roles"
        )
    )


class RoleEnum(str, Enum):
    CUSTOMER = "customer"
    SALON_OWNER = "salon_owner"
    EMPLOYEE = "employee"
