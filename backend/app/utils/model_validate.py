from typing import Any, TypeVar

from pydantic import BaseModel

OrmModelT = TypeVar("OrmModelT")
PydanticModelT = TypeVar("PydanticModelT", bound=BaseModel)

def validate_to_orm(
    orm_model: type[OrmModelT],
    pydantic_model: type[PydanticModelT],
    data: Any,
    **updates: Any,
) -> OrmModelT:
    validated = pydantic_model.model_validate(data, update=updates)
    payload = validated.model_dump(exclude_unset=True)
    return orm_model(**payload)