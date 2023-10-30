package code

import "github.com/mangohow/cloud-ide/pkg/serialize"

const (
	QuerySuccess = iota
	QueryFailed

	LoginSuccess
	LoginFailed
	LoginUserDeleted

	SpaceCreateSuccess
	SpaceCreateFailed
	SpaceCreateNameDuplicate
	SpaceCreateReachMaxCount

	SpaceStartSuccess
	SpaceStartFailed

	SpaceDeleteSuccess
	SpaceDeleteFailed
	SpaceDeleteIsRunning

	SpaceStopSuccess
	SpaceStopFailed
	SpaceStopIsNotRunning

	UserNameAvailable
	UserNameUnavailable
	UserSendValidateCodeSuccess
	UserSendValidateCodeFailed
	UserEmailCodeInvalid
	UserEmailInvalid
	UserUsernameExist
	UserRegisterSuccess
	UserRegisterFailed
	UserEmailCodeIncorrect
	UserEmailAlreadyInUse

	SpaceStartNotExist
	SpaceOtherSpaceIsRunning

	SpaceNameModifySuccess
	SpaceNameModifyFailed
	SpaceAlreadyExist
	SpaceNotFound
	ResourceExhausted
)

type UserStatus uint32

const (
	StatusNormal UserStatus = iota
	StatusDeleted
)

var messageForCode = map[int]string{
	QuerySuccess:                "查询成功",
	QueryFailed:                 "查询失败",
	LoginSuccess:                "登录成功",
	LoginFailed:                 "登录失败",
	LoginUserDeleted:            "用户已注销",
	SpaceCreateSuccess:          "创建成功",
	SpaceCreateFailed:           "创建失败",
	SpaceCreateNameDuplicate:    "不能和已有工作空间名称重复",
	SpaceCreateReachMaxCount:    "达到最大工作空间创建上限,请删除其它工作空间后重试",
	SpaceStartSuccess:           "工作空间启动成功",
	SpaceStartFailed:            "工作空间启动失败,请重试",
	SpaceDeleteSuccess:          "删除工作空间成功",
	SpaceDeleteFailed:           "删除工作空间失败",
	SpaceDeleteIsRunning:        "无法删除正在运行的工作空间,请先停止运行",
	SpaceStopSuccess:            "停止工作空间成功",
	SpaceStopFailed:             "停止工作空间失败",
	SpaceStopIsNotRunning:       "工作空间未运行",
	UserNameAvailable:           "用户名可用",
	UserNameUnavailable:         "用户名不可用",
	UserSendValidateCodeFailed:  "验证码发送失败,请重试",
	UserSendValidateCodeSuccess: "验证码发送成功,五分钟内有效",
	UserEmailCodeInvalid:        "邮箱验证码不合法",
	UserEmailInvalid:            "邮箱格式不正确",
	UserUsernameExist:           "用户名已存在",
	UserRegisterSuccess:         "注册成功",
	UserRegisterFailed:          "注册失败",
	UserEmailCodeIncorrect:      "邮箱验证码不正确",
	UserEmailAlreadyInUse:       "该邮箱已经被注册",
	SpaceStartNotExist:          "工作空间不存在",
	SpaceOtherSpaceIsRunning:    "检测到有其它工作空间正在运行,请先停止正在运行的工作空间",
	SpaceNameModifySuccess:      "名称修改成功",
	SpaceNameModifyFailed:       "名称修改失败",
	SpaceAlreadyExist:           "工作空间已存在",
	SpaceNotFound:               "未找到该工作空间",
	ResourceExhausted:           "资源不足,无法启动工作空间",
}

func GetMessage(code int) string {
	return messageForCode[code]
}

func init() {
	serialize.SetCodeMessager(GetMessage)
}
