# app/models/user.py
import uuid

from sqlmodel import Field

from .base import BaseModel


# Shared properties
class UserBase(BaseModel):
    email: str | None = Field(default=None, index=True, nullable=True, unique=True)
    phone_number: str | None = Field(default=None, index=True, nullable=True, unique=True)
    is_active: bool = True
    is_superuser: bool = False
    is_email_verified: bool = Field(default=False)
    is_phone_verified: bool = Field(default=False)
    full_name: str | None = Field(default=None, max_length=255)
    # role: UserRole = Field(sa_column=Column(Enum(UserRole)), default=UserRole.CUSTOMER)
    # oauth_accounts: list["OAuthAccount"] = Relationship(back_populates="user")


# Database model, database table inferred from class name
class User(UserBase, table=True):
    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    hashed_password: str

