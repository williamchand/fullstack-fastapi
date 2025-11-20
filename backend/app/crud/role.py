from sqlmodel import Session, select

from app.models.role import Role
from app.models.user_role import UserRole

from .base import CRUDBase


class CRUDRole(CRUDBase[Role]):

    def create_role(self, session: Session, name: str) -> Role:
        """
        Create a new role if not exists.
        """
        existing = session.exec(
            select(Role).where(Role.name == name)
        ).first()
        if existing:
            return existing

        role = Role(name=name)
        session.add(role)
        session.commit()
        session.refresh(role)
        return role

    def update_role(self, session: Session, role: Role, *, name: str) -> Role:
        """
        Update role name.
        """
        role.name = name
        session.add(role)
        session.commit()
        session.refresh(role)
        return role

    def delete_role(self, session: Session, role: Role) -> None:
        """
        Delete role and all UserRole links.
        """
        session.exec(
            select(UserRole).where(UserRole.role_id == role.id)
        ).delete()

        session.delete(role)
        session.commit()

    def get_by_name(self, session: Session, name: str) -> Role | None:
        return session.exec(
            select(Role).where(Role.name == name)
        ).first()

    def assign_to_user(self, session: Session, user_id, role_name: str):
        """
        Add a role to an existing user.
        """
        role = self.create_role(session, role_name)

        link = UserRole(user_id=user_id, role_id=role.id)
        session.add(link)
        session.commit()

    def remove_from_user(self, session: Session, user_id, role_name: str):
        """
        Remove a role from the user.
        """
        role = self.get_by_name(session, role_name)
        if not role:
            return

        session.exec(
            select(UserRole)
            .where(UserRole.user_id == user_id)
            .where(UserRole.role_id == role.id)
        ).delete()

        session.commit()

# 👇 Common pattern: create only once, import everywhere
role_crud = CRUDRole(Role)
