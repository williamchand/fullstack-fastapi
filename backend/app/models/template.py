from sqlmodel import Field

from .base import BaseModel, BaseModelUUID


class Template(BaseModel, BaseModelUUID, table=True):
    name: str = Field(max_length=100)
    theme_config: dict = Field(default_factory=dict)
    config_schema: dict = Field(default_factory=dict)
    preview_url: str | None = Field(default=None, max_length=512)
    price: float = 0.0
