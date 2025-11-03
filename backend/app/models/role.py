# app/models/role.py
from enum import Enum


class UserRole(str, Enum):
    CUSTOMER = "customer"
    SALON_OWNER = "salon_owner"
    EMPLOYEE = "employee"
