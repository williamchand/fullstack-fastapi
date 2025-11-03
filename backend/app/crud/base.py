# app/crud/crud_user.py
from typing import Generic, TypeVar

from sqlmodel import Session, SQLModel, select

ModelType = TypeVar("ModelType", bound=SQLModel)

class CRUDBase(Generic[ModelType]):
    def __init__(self, model: type[ModelType]):
        """
        CRUD object with default methods to Create, Read, Update, Delete (CRUD).
        **model**: a SQLModel class
        """
        self.model = model

    def get(self, session: Session, id: str) -> ModelType | None:
        return session.get(self.model, id)

    def get_multi(self, session: Session, skip: int = 0, limit: int = 100):
        statement = select(self.model).offset(skip).limit(limit)
        return session.exec(statement).all()

    def create(self, session: Session, obj_in: SQLModel) -> ModelType:
        db_obj = self.model(**obj_in.model_dump())
        session.add(db_obj)
        session.commit()
        session.refresh(db_obj)
        return db_obj

    def update(self, session: Session, db_obj: ModelType, obj_in: SQLModel) -> ModelType:
        obj_data = obj_in.model_dump(exclude_unset=True)
        for field, value in obj_data.items():
            setattr(db_obj, field, value)
        session.add(db_obj)
        session.commit()
        session.refresh(db_obj)
        return db_obj

    def remove(self, session: Session, id: str) -> ModelType:
        obj = session.get(self.model, id)
        if obj:
            session.delete(obj)
            session.commit()
        return obj

