from typing import Any

from sqlmodel import Session, select

from app.core.security import get_password_hash, verify_password
from app.models.user import User
from app.schemas.user import UserCreate, UserUpdate

from .base import CRUDBase


class CRUDUser(CRUDBase[User]):
    def create_user(self, session: Session, user_create: UserCreate) -> User:
        user_data = user_create.model_dump(exclude={"roles"})
        db_obj = User.model_validate(
            user_data, update={"hashed_password": get_password_hash(user_create.password)}
        )
        session.add(db_obj)
        session.commit()
        session.refresh(db_obj)
        if user_create.roles:
            db_obj.set_roles(roles=user_create.roles, session=session)

        return db_obj


    def update_user(self, session: Session, db_user: User, user_in: UserUpdate) -> Any:
        user_data = user_in.model_dump(exclude_unset=True, exclude={"roles"})
        extra_data = {}
        if "password" in user_data:
            password = user_data["password"]
            hashed_password = get_password_hash(password)
            extra_data["hashed_password"] = hashed_password
        db_user.sqlmodel_update(user_data, update=extra_data)
        session.add(db_user)
        session.commit()
        session.refresh(db_user)
        if user_in.roles:
            db_user.set_roles(roles=user_in.roles, session=session)
        return db_user


    def get_user_by_email(self, session: Session, email: str | None) -> User | None:
        statement = select(User).where(User.email == email)
        session_user = session.exec(statement).first()
        return session_user


    def authenticate(self, session: Session, email: str, password: str) -> User | None:
        db_user = self.get_user_by_email(session=session, email=email)
        if not db_user:
            return None
        if not verify_password(password, db_user.hashed_password):
            return None
        return db_user


# 👇 Common pattern: create only once, import everywhere
user_crud = CRUDUser(User)
