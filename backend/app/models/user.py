# app/models/user.py
from typing import TYPE_CHECKING

from sqlmodel import Field, Relationship, Session, select

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

    def validate_role(self, allowed: set[RoleEnum]) -> bool:
        """Check if the user has any allowed roles."""
        if not self.roles:
            return False

        allowed_lower = {role.value.lower() for role in allowed}
        return any(r.name.lower() in allowed_lower for r in self.roles)

    def set_roles(self, roles: list[RoleEnum], session: Session) -> None:
        """Replace all user roles with the provided list of RoleEnum values."""
        from .role import Role
        role_names = [r.value for r in roles]

        db_roles = session.exec(
            select(Role).where(Role.name.in_(role_names))
        ).all()

        if len(db_roles) != len(roles):
            missing = set(role_names) - {r.name for r in db_roles}
            raise ValueError(f"Roles not found: {missing}")

        self.roles = db_roles
        session.add(self)
        session.commit()
        session.refresh(self)
