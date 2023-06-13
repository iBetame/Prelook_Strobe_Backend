package Model

import "gorm.io/gorm"

type Admins struct {
	ID int `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Users struct {
	ID int
	NickName string `gorm:"column:NickName"`
	MediaName string `gorm:"column:MediaName"`
	UserName string `gorm:"column:UserName"`
	PassWord string `gorm:"column:PassWord"`
	UserType string `gorm:"column:UserType"`
	IP string
	Country string
}

// 定义个全局函数
func AdminTable() *gorm.DB {
	return db.Table("admins")
}

func CheckAuth(username, password string) bool {

	var auth Admins

	AdminTable().Select("id").Where("username =? AND password =?", username, password).First(&auth)

	if auth.ID > 0 {
		return true
	}

	return false
}

func VerifySudoUser(myuser Users) (result bool) {
	var tempPassWord = myuser.PassWord
	db.Where("username = ?", myuser.UserName).First(&myuser)
	if myuser.PassWord == tempPassWord {
		return true
	} else {
		return false
	}
}

func AddNewUser(myuser Users) (final Users) {
	var olduser Users
	// 先用用户名去找 id，存到 tempuser 里，如果 tempuser 没有 id，就 insert，否则 update
	db.Where("username = ?", myuser.UserName).First(&olduser)
	// olduser.ID == 0 证明是新的，需要以新的加入
	if olduser.ID == 0  {
		// 新的加入
		db.Create(&myuser)
	} else {
		// 用收到的覆盖掉查询出来的
		db.Where("username = ?", olduser.UserName).Updates(Users{NickName: myuser.NickName, MediaName: myuser.MediaName, Country: myuser.Country})
	}
	//返回刚添加这个用户
	db.Where("UserName = ?", myuser.UserName).Find(&final)
	return final
}

func GetPassWord() (toppassword string, err error) {
	var mysettings Settings
	// 获取全部设置选项
	db.Table("settings").Where("SettingsItem = ?", "TopPassWord").First(&mysettings)
	if err != nil {
		return mysettings.SettingsContent, err
	}
	return mysettings.SettingsContent, err
}

func GetUserName(UserID string)(UserName string)  {
	var myuser Users
	// 获取全部设置选项
	db.Table("users").Where("ID = ?", UserID).First(&myuser)
	return myuser.UserName
}

func GetUserByID(ID int) (user Users)  {
	var myuser Users
	// 获取全部设置选项
	db.Table("users").Where("ID = ?", ID).First(&myuser)
	return myuser
}