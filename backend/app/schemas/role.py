import uuid

from sqlmodel import Field

from app.models.base import BaseModel
from app.models.role import RoleEnum
from app.schemas.base import FromModelMixin, PaginatedListResponseMixin


class RoleBase(BaseModel):
    name: RoleEnum = Field(description="System role name")

class RoleCreate(RoleBase):
    """Payload used when creating a new role."""
    pass

class RoleUpdate(BaseModel):
    name: RoleEnum

class RolePublic(FromModelMixin, RoleBase):
    id: int

    @staticmethod
    def _transform_name(name: str) -> RoleEnum:
        return RoleEnum(name)

    __field_transformers__ = {
        "name": _transform_name,
    }

class RolesPublic(PaginatedListResponseMixin[RolePublic], BaseModel):
    pass
