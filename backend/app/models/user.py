# app/models/user.py
from typing import TYPE_CHECKING

from sqlmodel import Field, Relationship

from .base import BaseModel, BaseModelUUID
from .role import RoleEnum
from .user_role import UserRole

if TYPE_CHECKING:  # 👈 avoids circular import at runtime
    from .oauth_account import OAuthAccount
    from .role import Role


# Shared properties
class UserBase(BaseModel):
    email: str | None = Field(default=None, index=True, nullable=True, unique=True)
    phone_number: str | None = Field(default=None, index=True, nullable=True, unique=True)
    is_active: bool = True
    is_superuser: bool = True
    is_email_verified: bool = Field(default=False)
    is_phone_verified: bool = Field(default=False)
    is_totp_enabled: bool = Field(default=False)
    totp_secret: str | None = Field(default=None)
    full_name: str | None = Field(default=None, max_length=255)

# Database model, database table inferred from class name
class User(UserBase, BaseModelUUID, table=True):
    hashed_password: str
    oauth_accounts: list["OAuthAccount"]  = Relationship(back_populates="user")
    roles: list["Role"] = Relationship(
        back_populates="users",
        link_model=UserRole,
    )

    @property
    def is_superuser(self) -> bool:
        """
        Automatically computed:
        If user has any role named admin / superadmin → Superuser
        """
        if not self.roles:
            return False

        privileged_roles = {RoleEnum.SUPERUSER}
        return any(role.name.lower() in privileged_roles for role in self.roles)

