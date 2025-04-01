package services

import (
	sql "github.com/fegig/goLang_class/database"
)

type UserData struct {
	ID        int    `json:"id"`
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func GetUsers(limit, offset int) ([]UserData, error) {
	orderBy := "id"
	users, err := sql.SelectData("users", []string{"id", "userId", "firstName", "lastName", "email"}, []sql.FieldValue{}, &orderBy, &limit, &offset)
	if err != nil {
		return nil, err
	}
	defer users.Close()

	var usersData []UserData
	for users.Next() {
		var user UserData
		if err := users.Scan(&user.ID, &user.UserId, &user.FirstName, &user.LastName, &user.Email); err != nil {
			return nil, err
		}
		usersData = append(usersData, user)
	}
	return usersData, nil
}

func GetUser(userID string) (UserData, error) {
	user, err := sql.SelectData("users", []string{"id", "userId", "firstName", "lastName", "email"}, []sql.FieldValue{
		{Field: "userId", Value: userID},
	}, nil, nil, nil)
	if err != nil {
		return UserData{}, err
	}
	defer user.Close()

	var userData UserData
	if user.Next() {
		if err := user.Scan(&userData.ID, &userData.UserId, &userData.FirstName, &userData.LastName, &userData.Email); err != nil {
			return UserData{}, err
		}
	}
	return userData, nil
}
