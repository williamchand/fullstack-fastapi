# app/models/user.py
import uuid

from sqlmodel import Field, Relationship

from .base import BaseModel


# Shared properties
class UserBase(BaseModel):
    email: str | None = Field(default=None, index=True, nullable=True, unique=True)
    # phone_number: Optional[str] = Field(default=None, index=True, nullable=True, unique=True)
    is_active: bool = True
    is_superuser: bool = False
    # is_email_verified: bool = Field(default=False)  # True if email OR phone is verified
    # is_phone_verified: bool = Field(default=False)  # True if email OR phone is verified
    full_name: str | None = Field(default=None, max_length=255)
    # role: UserRole = Field(sa_column=Column(Enum(UserRole)), default=UserRole.CUSTOMER)
    # oauth_accounts: list["OAuthAccount"] = Relationship(back_populates="user")


# Database model, database table inferred from class name
class User(UserBase, table=True):
    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    hashed_password: str

