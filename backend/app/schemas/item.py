
import uuid

from sqlmodel import Field

from app.models.base import BaseModel
from app.models.item import ItemBase


# Properties to receive on item creation
class ItemCreate(ItemBase):
    pass


# Properties to receive on item update
class ItemUpdate(ItemBase):
    title: str | None = Field(default=None, min_length=1, max_length=255)  # type: ignore


# Properties to return via API, id is always required
class ItemPublic(ItemBase):
    id: uuid.UUID
    owner_id: uuid.UUID


class ItemsPublic(BaseModel):
    data: list[ItemPublic]
    count: int

