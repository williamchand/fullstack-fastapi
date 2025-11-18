# app/models/wedding.py
from typing import TYPE_CHECKING

from sqlmodel import Field, Relationship

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
    config_data: dict = Field(default_factory=dict)
    deleted_at: str | None = None  # datetime, nullable


class Wedding(WeddingBase, BaseModelUUID, BaseModelSoftDelete, table=True):
    user_id: str = Field(foreign_key="user.id")
    template_id: str | None = Field(default=None, foreign_key="template.id")
    payment_id: str | None = Field(default=None, foreign_key="payment.id")

    user: "User" = Relationship()
    template: "Template" = Relationship()
    payment: "Payment" = Relationship()
    guests: list["Guest"] = Relationship(back_populates="wedding")
