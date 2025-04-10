package models

import (
	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/bcrypt"
)

// Users -
type User struct {
	Id       int    `orm:"auto"`
	Username string `orm:"size(100);unique"`
	Password string `orm:"size(255)"` // Hashed password
	Role     string `orm:"size(50)"`  // e.g., "admin", "user"
}

func init() {
	orm.RegisterModel(new(User))
}

func AddUser(username, password, role string) error {
	o := orm.NewOrm()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &User{Username: username, Password: string(hashedPassword), Role: role}
	_, err = o.Insert(user)
	return err
}

func GetUserById(id int) (*User, error) {
	o := orm.NewOrm()
	user := &User{}
	err := o.QueryTable("user").Filter("id", id).One(user)
	return user, err
}

func GetUserByUsername(username string) (*User, error) {
	o := orm.NewOrm()
	user := &User{Username: username}
	err := o.Read(user, "Username")
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func GetAllUsers() ([]*User, error) {
	o := orm.NewOrm()
	var users []*User
	_, err := o.QueryTable("user").All(&users)
	return users, err
}

func GetUserCount() (int64, error) {
	o := orm.NewOrm()
	count, err := o.QueryTable("user").Count()
	return count, err
}

func UpdateUser(user *User) error {
	o := orm.NewOrm()
	_, err := o.Update(user)
	return err
}

func DeleteUser(userID int) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("user").Filter("id", userID).Delete()
	return err
}

func AuthenticateUser(username, password string) (*User, error) {
	user, err := GetUserByUsername(username)
	if err != nil || user == nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}
	return user, nil
}
