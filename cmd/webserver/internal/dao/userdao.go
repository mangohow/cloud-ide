package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/mangohow/cloud-ide/cmd/webserver/internal/dao/db"
	"github.com/mangohow/cloud-ide/cmd/webserver/internal/model"
)

type UserDao struct {
	db *sqlx.DB
}

func NewUserDao() *UserDao {
	return &UserDao{
		db: db.DB(),
	}
}

func (u *UserDao) FindByUsernameAndPassword(username, password string) (user *model.User, _ error) {
	sql := `SELECT id, uid, username, nickname, email, avatar, status FROM t_user WHERE username = ? AND password = ?`
	user = &model.User{}
	err := u.db.Get(user, sql, username, password)
	return user, err
}

func (u *UserDao) FindByUsername(username string) error {
	sql := "SELECT 1 FROM t_user WHERE username = ?"
	var n int
	return u.db.Get(&n, sql, username)
}

func (u *UserDao) FindByEmail(email string) error {
	sql := `SELECT 1 FROM t_user WHERE email = ?`
	var n int
	return u.db.Get(&n, sql, email)
}

func (u *UserDao) AddUser(user *model.User) error {
	sql := `Insert into t_user (uid, username, password, nickname, email, create_time, delete_time, status) values (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := u.db.Exec(sql, user.Uid, user.Username, user.Password, user.Nickname, user.Email, user.CreateTime, user.DeleteTime, user.Status)
	return err
}
