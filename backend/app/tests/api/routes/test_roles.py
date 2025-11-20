import uuid
from sqlmodel import Session, select
from fastapi.testclient import TestClient

from app.models.role import Role, RoleEnum
from app.models.user import User
from app.schemas.user import UserCreate
from app.crud.user import user_crud
from app.core.config import settings
from app.tests.utils.utils import random_lower_string, random_email


# ---------------------------------------------------------
# ROLE LISTING — SUPERUSER ONLY
# ---------------------------------------------------------
def test_list_roles_superuser(client: TestClient, superuser_token_headers: dict[str, str], db: Session):
    r = client.get(f"{settings.API_V1_STR}/roles/", headers=superuser_token_headers)
    assert r.status_code == 200
    data = r.json()

    # default roles that must exist
    expected = set([role.value for role in RoleEnum])

    returned = set([item["name"] for item in data])
    assert expected == returned


def test_list_roles_normal_user_forbidden(client: TestClient, normal_user_token_headers: dict[str, str]):
    r = client.get(f"{settings.API_V1_STR}/roles/", headers=normal_user_token_headers)
    assert r.status_code == 403
    assert r.json() == {"detail": "The user doesn't have enough privileges"}


# ---------------------------------------------------------
# ASSIGN ROLE TO USER — SUPERUSER ONLY
# ---------------------------------------------------------
def test_assign_role_to_user(client: TestClient, superuser_token_headers: dict[str, str], db: Session):
    email = random_email()
    pwd = random_lower_string()

    # Create new user with CUSTOMER role
    user_in = UserCreate(email=email, password=pwd, roles=[RoleEnum.CUSTOMER])
    user = user_crud.create_user(session=db, user_create=user_in)

    # Promote to STAFF
    r = client.patch(
        f"{settings.API_V1_STR}/roles/{user.id}",
        headers=superuser_token_headers,
        json={"roles": [RoleEnum.STAFF]},
    )

    assert r.status_code == 200
    updated = r.json()
    assert updated["roles"] == [RoleEnum.STAFF]

    db.refresh(user)
    assert [RoleEnum(r.name) for r in user.roles] == [RoleEnum.STAFF]


def test_assign_role_to_user_forbidden(client: TestClient, normal_user_token_headers: dict[str, str], db: Session):
    email = random_email()
    pwd = random_lower_string()

    user_in = UserCreate(email=email, password=pwd, roles=[RoleEnum.CUSTOMER])
    user = user_crud.create_user(session=db, user_create=user_in)

    r = client.patch(
        f"{settings.API_V1_STR}/roles/{user.id}",
        headers=normal_user_token_headers,
        json={"roles": [RoleEnum.ADMIN]},
    )

    assert r.status_code == 403
    assert r.json() == {"detail": "The user doesn't have enough privileges"}


# ---------------------------------------------------------
# GET SINGLE USER ROLE
# ---------------------------------------------------------
def test_get_user_roles(client: TestClient, superuser_token_headers: dict[str, str], db: Session):
    email = random_email()
    pwd = random_lower_string()

    user_in = UserCreate(email=email, password=pwd, roles=[RoleEnum.CUSTOMER])
    user = user_crud.create_user(session=db, user_create=user_in)

    r = client.get(
        f"{settings.API_V1_STR}/roles/{user.id}",
        headers=superuser_token_headers,
    )

    assert r.status_code == 200
    roles = r.json()
    assert roles == [RoleEnum.CUSTOMER]


def test_get_user_roles_not_found(client: TestClient, superuser_token_headers: dict[str, str]):
    r = client.get(
        f"{settings.API_V1_STR}/roles/{uuid.uuid4()}",
        headers=superuser_token_headers,
    )
    assert r.status_code == 404
    assert r.json()["detail"] == "User not found"


# ---------------------------------------------------------
# REMOVE ROLE FROM USER
# ---------------------------------------------------------
def test_remove_role_from_user(client: TestClient, superuser_token_headers: dict[str, str], db: Session):
    email = random_email()
    pwd = random_lower_string()

    user_in = UserCreate(email=email, password=pwd, roles=[RoleEnum.ADMIN])
    user = user_crud.create_user(session=db, user_create=user_in)

    r = client.delete(
        f"{settings.API_V1_STR}/roles/{user.id}/{RoleEnum.ADMIN}",
        headers=superuser_token_headers,
    )

    assert r.status_code == 200
    result = r.json()
    assert result["message"] == "Role removed successfully"

    db.refresh(user)
    assert user.roles == []


def test_remove_role_forbidden(client: TestClient, normal_user_token_headers: dict[str, str], db: Session):
    email = random_email()
    pwd = random_lower_string()

    user_in = UserCreate(email=email, password=pwd, roles=[RoleEnum.ADMIN])
    user = user_crud.create_user(session=db, user_create=user_in)

    r = client.delete(
        f"{settings.API_V1_STR}/roles/{user.id}/{RoleEnum.ADMIN}",
        headers=normal_user_token_headers,
    )

    assert r.status_code == 403


# ---------------------------------------------------------
# PREVENT SUPERUSER ROLE REMOVAL FROM SELF
# ---------------------------------------------------------
def test_remove_own_superuser_role_forbidden(client: TestClient, superuser_token_headers: dict[str, str], db: Session):
    superuser = db.exec(select(User).where(User.email == settings.FIRST_SUPERUSER)).first()
    assert superuser is not None

    r = client.delete(
        f"{settings.API_V1_STR}/roles/{superuser.id}/{RoleEnum.SUPERUSER}",
        headers=superuser_token_headers,
    )

    assert r.status_code == 403
    assert r.json()["detail"] == "Super users are not allowed to remove their own superuser role"
