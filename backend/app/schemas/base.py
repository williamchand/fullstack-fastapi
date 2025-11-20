from collections.abc import Callable
from typing import Any, Generic, Protocol, TypeVar, get_args

from pydantic import BaseModel


# ----------------------------------------------------------------------
# ORM Model Protocol
# ----------------------------------------------------------------------
class ORMModel(Protocol):
    """
    Minimal protocol for ORM-like objects.
    Allows SQLAlchemy/SQLModel instances to pass type checking.
    """
    ...


M = TypeVar("M", bound=ORMModel)
P = TypeVar("P", bound="FromModelMixin")
T = TypeVar("T", bound="FromModelMixin")


# ----------------------------------------------------------------------
# Mixin: Convert ORM → Pydantic models
# ----------------------------------------------------------------------
class FromModelMixin(BaseModel):
    """
    Adds:
        - .from_model(orm_obj)
        - .from_list(list_of_orm)

    Supports:
        - custom per-field transformers
        - SQLAlchemy and SQLModel instances
    """

    __field_transformers__: dict[str, Callable[[Any], Any]] = {}

    @classmethod
    def from_model(cls: type[P], model_obj: M) -> P:
        # Step 1: create empty Pydantic model without validation
        obj: P = cls.model_construct()

        # Step 2: Fill fields from the model object (bypass validation)
        for field_name, field in cls.model_fields.items():
            value = getattr(model_obj, field_name, None)
            setattr(obj, field_name, value)

        # Step 3: Run custom field transformers (manual conversion)
        for field_name, transformer in cls.__field_transformers__.items():
            raw_value = getattr(model_obj, field_name, None)
            transformed = transformer(raw_value)
            setattr(obj, field_name, transformed)

        # Step 4: Validate the final object
        obj = cls.model_validate(obj)

        return obj

    @classmethod
    def from_list(cls: type[P], items: list[M]) -> list[P]:
        return [cls.from_model(item) for item in items]


# ----------------------------------------------------------------------
# Pagination Metadata
# ----------------------------------------------------------------------
class ListMeta(BaseModel):
    count: int
    page: int | None = None
    page_size: int | None = None
    total: int | None = None
    has_next: bool | None = None


# ----------------------------------------------------------------------
# Paginated Response Mixin
# ----------------------------------------------------------------------
class PaginatedListResponseMixin(BaseModel, Generic[T]):
    """
    data: list[T]  (T must be FromModelMixin)
    meta: pagination metadata
    """

    data: list[T]
    meta: ListMeta

    @classmethod
    def from_model(
        cls,
        models: list[M],
        *,
        skip: int = 0,
        limit: int = 100,
        total: int | None = None,
    ):
        # Extract T from data: list[T]
        data_annotation = cls.model_fields["data"].annotation
        (item_type,) = get_args(data_annotation)

        page_size = limit
        page = (skip // limit) + 1
        meta = ListMeta(
            count=len(models),
            page=page,
            page_size=page_size,
            total=total,
            has_next=(page_size is not None and len(models) == page_size),
        )

        return cls(
            data=item_type.from_list(models),
            meta=meta,
        )
