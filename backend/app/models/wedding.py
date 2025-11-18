# app/models/wedding.py
import uuid
from typing import TYPE_CHECKING

from sqlmodel import JSON, Field, Relationship

from .base import BaseModel, BaseModelSoftDelete, BaseModelUUID

if TYPE_CHECKING:
    from .guest import Guest
    from .payment import Payment
    from .template import Template
    from .user import User


class WeddingBase(BaseModel):
    status: str = Field(default="draft")
    custom_domain: str | None = Field(default=None, max_length=255)
    slug: str | None = Field(default=None, max_length=150)
    config_data: dict = Field(default_factory=dict, sa_type=JSON)


class Wedding(WeddingBase, BaseModelUUID, BaseModelSoftDelete, table=True):
    user_id: uuid.UUID = Field(foreign_key="user.id")
    template_id: uuid.UUID | None = Field(default=None, foreign_key="template.id")
    payment_id: uuid.UUID | None = Field(default=None, foreign_key="payment.id")

    user: "User" = Relationship()
    template: "Template" = Relationship()
    payment: "Payment" = Relationship()
    guests: list["Guest"] = Relationship(back_populates="wedding")
