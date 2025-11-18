from sqlmodel import Field

from .base import BaseModel, BaseModelUUID


class PaymentMethod(BaseModel, BaseModelUUID, table=True):
    name: str = Field(max_length=100)
    provider: str = Field(max_length=50)
    config: dict = Field(default_factory=dict)
    is_active: bool = True
