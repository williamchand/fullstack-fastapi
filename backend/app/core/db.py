from sqlmodel import Session, SQLModel, create_engine, select

from app.core.config import settings
from app.crud.user import user_crud
from app.models.role import Role, RoleEnum
from app.models.user import User
from app.models.user_role import UserRole
from app.schemas.user import UserCreate

from .base import import_all_models

import_all_models()
engine = create_engine(str(settings.SQLALCHEMY_DATABASE_URI))


# make sure all SQLModel models are imported (app.models) before initializing DB
# otherwise, SQLModel might fail to initialize relationships properly
# for more details: https://github.com/fastapi/full-stack-fastapi-template/issues/28

def init_db(session: Session) -> None:
    # Tables should be created with Alembic migrations
    # But if you don't want to use migrations, create
    # the tables un-commenting the next lines
    # from sqlmodel import SQLModel

    _ = SQLModel.metadata  # load metadata to ensure model mappers are registered
    # This works because the models are already imported and registered from app.models
    # SQLModel.metadata.create_all(engine)

    role = session.exec(
        select(Role).where(Role.name == RoleEnum.SUPERUSER)
    ).first()

    # 2. Ensure superuser user exists
    user = session.exec(
        select(User).where(User.email == settings.FIRST_SUPERUSER)
    ).first()

    if not user:
        user_in = UserCreate(
            email=settings.FIRST_SUPERUSER,
            password=settings.FIRST_SUPERUSER_PASSWORD,
            is_email_verified=True,
            is_phone_verified=True,
        )
        user = user_crud.create_user(session=session, user_create=user_in)

    # 3. Attach superuser role to the user (if not already attached)
    has_super_role = any(r.name == RoleEnum.SUPERUSER for r in user.roles)
    if not has_super_role:
        user_role = UserRole(user_id=user.id, role_id=role.id)
        session.add(user_role)
        session.commit()

    return user

    user = session.exec(
        select(User).where(User.email == settings.FIRST_SUPERUSER)
    ).first()
    if not user:
        user_in = UserCreate(
            email=settings.FIRST_SUPERUSER,
            password=settings.FIRST_SUPERUSER_PASSWORD,
            is_superuser=True,
            is_email_verified=True,
        )
        user = user_crud.create_user(session=session, user_create=user_in)
