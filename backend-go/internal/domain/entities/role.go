package entities

type RoleEnum string

const (
	RoleCustomer   RoleEnum = "customer"
	RoleSalonOwner RoleEnum = "salon_owner"
	RoleEmployee   RoleEnum = "employee"
	RoleSuperuser  RoleEnum = "superuser"
)
