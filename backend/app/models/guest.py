# app/models/guest.py
from typing import TYPE_CHECKING

from sqlmodel import Field, Relationship

from .base import BaseModel, BaseModelSoftDelete, BaseModelUUID

if TYPE_CHECKING:
    from .wedding import Wedding


class GuestBase(BaseModel):
    name: str = Field(max_length=255)
    contact: str = Field(max_length=255)
    rsvp_status: str = Field(default="maybe")
    message: str | None = None


class Guest(GuestBase, BaseModelUUID, BaseModelSoftDelete, table=True):
    wedding_id: str = Field(foreign_key="wedding.id")

    wedding: "Wedding" = Relationship(back_populates="guests")
