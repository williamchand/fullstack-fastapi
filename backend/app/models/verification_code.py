import uuid
from datetime import datetime

from sqlmodel import JSON, Column, DateTime, Field, text

from .base import BaseModel, BaseModelUUID


class VerificationCode(BaseModel, BaseModelUUID, table=True):
    __tablename__ = "verification_code"

    user_id: uuid.UUID = Field(foreign_key="user.id", index=True)

    verification_code: str = Field(index=True)
    verification_type: str = Field(index=True)  # "email_verification", "phone_verification", "totp_recovery", "password_reset"

    created_at: datetime | None = Field(
        default=None,
        sa_column=Column(
            DateTime(timezone=True),
            server_default=text('NOW()'),
        ),
    )
    expires_at: datetime
    extra_metadata: dict = Field(default_factory=dict, sa_type=JSON)
