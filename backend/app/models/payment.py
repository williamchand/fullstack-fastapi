import uuid
from typing import TYPE_CHECKING

from sqlmodel import JSON, Field, Relationship

from .base import BaseModel, BaseModelUUID

if TYPE_CHECKING:
    from .payment_method import PaymentMethod
    from .user import User


class PaymentBase(BaseModel):
    amount: float
    currency: str = Field(default="USD", max_length=10)
    status: str = Field(
        default="pending",
        sa_column_kwargs={"server_default": "pending"}  # matches Enum default
    )
    transaction_id: str = Field(max_length=255, unique=True)
    extra_metadata: dict = Field(default_factory=dict, sa_type=JSON)


class Payment(PaymentBase, BaseModelUUID, table=True):
    user_id: uuid.UUID = Field(foreign_key="user.id")
    payment_method_id: uuid.UUID | None = Field(default=None, foreign_key="payment_method.id")

    user: "User" = Relationship()
    payment_method: "PaymentMethod" = Relationship()
