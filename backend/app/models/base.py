# app/models/base.py
import uuid
from datetime import datetime

from sqlalchemy import text
from sqlalchemy.orm import with_loader_criteria
from sqlmodel import Field, SQLModel


class BaseModel(SQLModel):
    """
    Shared configuration for all entities
    (can be extended later for timestamps or common methods).
    """
    pass

class BaseModelUUID(SQLModel, table=False):
    id: uuid.UUID | None = Field(
        default=None,
        primary_key=True,
        nullable=False,
        sa_column_kwargs={"server_default": text("gen_random_uuid()")},
    )

class BaseModelSoftDelete(SQLModel):
    deleted_at: datetime | None = Field(default=None, nullable=True)

    @classmethod
    def _active_filter(cls):
        return cls.deleted_at.is_(None)

    @classmethod
    def __declare_last__(cls):
        from sqlalchemy import event
        from sqlalchemy.orm import Session

        # Apply global query criteria automatically
        @event.listens_for(Session, "do_orm_execute")
        def _add_filter(execute_state):
            if (
                execute_state.is_select
                and not execute_state.execution_options.get("include_deleted", False)
            ):
                execute_state.statement = execute_state.statement.options(
                    with_loader_criteria(cls, cls._active_filter())
                )

# Generic message
class Message(BaseModel):
    message: str
