from sqlmodel import JSON, Field

from .base import BaseModel, BaseModelUUID


class PaymentMethod(BaseModel, BaseModelUUID, table=True):
    __tablename__ = "payment_method"
    name: str = Field(max_length=100)
    provider: str = Field(max_length=50)
    config: dict = Field(default_factory=dict, sa_type=JSON)
    is_active: bool = True
