package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func checkAdmin(role string) error {
	if role != "admin" {
		return errors.New("forbidden")
	}
	return nil
}

func CreateUserService(email, name, password, role, currentRole string) (*User, error) {

	if err := checkAdmin(currentRole); err != nil {
		return nil, err
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return CreateUser(email, name, string(hash), role)
}

func GetUsersService(currentRole string) ([]User, error) {

	if err := checkAdmin(currentRole); err != nil {
		return nil, err
	}

	return GetAllUsers()
}

func UpdateUserService(id int64, email, name, role, currentRole string) error {

	if err := checkAdmin(currentRole); err != nil {
		return err
	}

	return UpdateUser(id, email, name, role)
}

func DeleteUserService(id int64, currentRole string) error {

	if err := checkAdmin(currentRole); err != nil {
		return err
	}

	return DeleteUser(id)
}
