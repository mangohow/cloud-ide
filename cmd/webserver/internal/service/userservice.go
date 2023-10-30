package service

import (
	"errors"
	"time"

	"github.com/mangohow/cloud-ide/cmd/webserver/internal/code"
	"github.com/mangohow/cloud-ide/cmd/webserver/internal/dao"
	"github.com/mangohow/cloud-ide/cmd/webserver/internal/model"
	"github.com/mangohow/cloud-ide/pkg/logger"
	"github.com/mangohow/cloud-ide/pkg/utils/encrypt"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type UserService struct {
	logger       *logrus.Logger
	dao          *dao.UserDao
	emailService EmailService
}

func NewUserService(service EmailService) *UserService {
	return &UserService{
		logger:       logger.Logger(),
		dao:          dao.NewUserDao(),
		emailService: service,
	}
}

var ErrUserDeleted = errors.New("user deleted")

func (u *UserService) Login(username, password string) (*model.User, error) {
	// 1、从数据库中查询
	user, err := u.dao.FindByUsernameAndPassword(username, password)
	if err != nil {
		return nil, err
	}

	// 2、检查用户状态是否正常
	if code.UserStatus(user.Status) == code.StatusDeleted {
		return nil, ErrUserDeleted
	}

	// 3、生成token
	token, err := encrypt.CreateToken(user.Id, user.Username, user.Uid)
	if err != nil {
		return nil, err
	}
	user.Token = token

	return user, nil
}

func (u *UserService) CheckUsernameAvailable(username string) bool {
	err := u.dao.FindByUsername(username)
	// 如果能查询到记录， err == nil
	if err != nil {
		return false
	}

	return true
}

var (
	ErrEmailCodeIncorrect = errors.New("email code incorrect")
	ErrEmailAlreadyInUse  = errors.New("this email had been registered")
)

func (u *UserService) UserRegister(info *model.RegisterInfo) error {
	// 1.验证邮箱验证码
	err := u.emailService.VerifyEmailValidateCode(info.Email, info.EmailCode)
	if err != nil {
		u.logger.Infof("verify email code failed err:%v", err)
		return ErrEmailCodeIncorrect
	}

	// 2.验证是否已经存在账号与该邮箱关联，一个邮箱只能创建一个账号
	if !u.emailService.IsEmailAvailable(info.Email) { // 如果err == nil说明查找到了记录
		return ErrEmailAlreadyInUse
	}

	// 3.生成新的用户
	now := time.Now()
	user := &model.User{
		Uid:        bson.NewObjectId().Hex(),
		Username:   info.Username,
		Password:   info.Password,
		Nickname:   info.Nickname,
		Email:      info.Email,
		CreateTime: now,
		DeleteTime: now,
	}

	err = u.dao.AddUser(user)
	if err != nil {
		return err
	}

	return nil
}
