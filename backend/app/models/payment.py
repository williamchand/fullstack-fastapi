from typing import TYPE_CHECKING

from sqlmodel import Field, Relationship

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
    metadata: dict = Field(default_factory=dict)


class Payment(PaymentBase, BaseModelUUID, table=True):
    user_id: str = Field(foreign_key="user.id")
    payment_method_id: str | None = Field(default=None, foreign_key="payment_method.id")

    user: "User" = Relationship()
    payment_method: "PaymentMethod" = Relationship()
