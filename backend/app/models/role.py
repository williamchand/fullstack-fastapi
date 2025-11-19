# app/models/role.py

from enum import Enum
from typing import TYPE_CHECKING

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
    users: list["User"] = Relationship(
        back_populates="roles",
        link_model=UserRole,
    )


class RoleEnum(str, Enum):
    CUSTOMER = "customer"
    SALON_OWNER = "salon_owner"
    EMPLOYEE = "employee"
    SUPERUSER = "superuser"
