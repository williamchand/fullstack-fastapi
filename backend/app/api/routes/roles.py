# app/api/routes/role.py

import uuid

from fastapi import APIRouter, Depends, HTTPException
from sqlmodel import func, select

from app.api.deps import SessionDep, get_current_active_superuser
from app.crud.role import role_crud
from app.models.base import Message
from app.models.role import Role
from app.schemas.role import (
    RoleCreate,
    RolePublic,
    RolesPublic,
    RoleUpdate,
)

router = APIRouter(prefix="/roles", tags=["roles"])


# ─────────────────────────────────────────────
# LIST ROLES
# ─────────────────────────────────────────────
@router.get(
    "/",
    dependencies=[Depends(get_current_active_superuser)],
    response_model=RolesPublic,
)
def list_roles(session: SessionDep, skip: int = 0, limit: int = 100):
    count = session.exec(select(func.count()).select_from(Role)).one()

    roles = session.exec(
        select(Role).offset(skip).limit(limit)
    ).all()

    return RolesPublic.from_model(
        roles,
        total=count,
        skip=skip,
        limit=limit,
    )


# ─────────────────────────────────────────────
# CREATE ROLE
# ─────────────────────────────────────────────
@router.post(
    "/",
    dependencies=[Depends(get_current_active_superuser)],
    response_model=RolePublic,
)
def create_role(session: SessionDep, role_in: RoleCreate):
    existing = role_crud.get_by_name(session, role_in.name)
    if existing:
        raise HTTPException(409, detail="Role with this name already exists")

    role = role_crud.create(session, role_in)
    return RolePublic.from_model(role)


# ─────────────────────────────────────────────
# GET ROLE BY ID
# ─────────────────────────────────────────────
@router.get(
    "/{role_id}",
    dependencies=[Depends(get_current_active_superuser)],
    response_model=RolePublic,
)
def get_role(role_id: uuid.UUID, session: SessionDep):
    role = session.get(Role, role_id)
    if not role:
        raise HTTPException(404, detail="Role not found")
    return RolePublic.from_model(role)


# ─────────────────────────────────────────────
# UPDATE ROLE
# ─────────────────────────────────────────────
@router.patch(
    "/{role_id}",
    dependencies=[Depends(get_current_active_superuser)],
    response_model=RolePublic,
)
def update_role(role_id: uuid.UUID, session: SessionDep, role_in: RoleUpdate):
    role = session.get(Role, role_id)
    if not role:
        raise HTTPException(404, detail="Role not found")

    if role_in.name:
        existing = role_crud.get_by_name(session, role_in.name)
        if existing and existing.id != role_id:
            raise HTTPException(409, detail="Role with this name already exists")

    role = role_crud.update(session, role, role_in)
    return RolePublic.from_model(role)


# ─────────────────────────────────────────────
# DELETE ROLE
# ─────────────────────────────────────────────
@router.delete(
    "/{role_id}",
    dependencies=[Depends(get_current_active_superuser)],
    response_model=Message,
)
def delete_role(role_id: uuid.UUID, session: SessionDep):
    role = session.get(Role, role_id)
    if not role:
        raise HTTPException(404, detail="Role not found")

    session.delete(role)
    session.commit()
    return Message(message="Role deleted successfully")