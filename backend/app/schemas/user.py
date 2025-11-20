# app/schemas/user.py
import uuid
from typing import TYPE_CHECKING

from pydantic import EmailStr
from sqlmodel import Field

from app.models.base import BaseModel
from app.models.role import RoleEnum
from app.models.user import UserBase
from app.schemas.base import FromModelMixin, PaginatedListResponseMixin

if TYPE_CHECKING:  # 👈 avoids circular import at runtime
    from app.models.role import Role

class UserResponseBase(FromModelMixin, UserBase):
    roles: list[RoleEnum] = Field(default_factory=list)

    def _transform_roles(rs: list["Role"] | None) -> list["RoleEnum"]:
        return [RoleEnum(r.name) for r in (rs or [])]

    __field_transformers__ = {
        "roles": _transform_roles
    }

# Properties to receive via API on creation
class UserCreate(UserResponseBase):
    password: str = Field(min_length=8, max_length=40)


class UserRegister(BaseModel):
    email: EmailStr = Field(max_length=255)
    password: str = Field(min_length=8, max_length=40)
    full_name: str | None = Field(default=None, max_length=255)


# Properties to receive via API on update, all are optional
class UserUpdate(UserResponseBase):
    roles: list["RoleEnum"] = Field(default_factory=list)
    email: EmailStr | None = Field(default=None, max_length=255)
    password: str | None = Field(default=None, min_length=8, max_length=40)


class UserUpdateMe(BaseModel):
    full_name: str | None = Field(default=None, max_length=255)
    email: EmailStr | None = Field(default=None, max_length=255)

# Properties to return via API, id is always required
class UserPublic(UserResponseBase):
    id: uuid.UUID

    # __field_transformers__ = {
    #     **UserResponseBase.__field_transformers__,
    # }


class UsersPublic(PaginatedListResponseMixin[UserPublic], BaseModel):
    pass
