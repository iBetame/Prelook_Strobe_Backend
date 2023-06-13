package Model

import (
	"github.com/imroc/req/v3"
	"github.com/sunbelife/Prelook_Strobe_Backend/Config"
	"log"
	"strconv"
	"time"
)

type TokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	ErrCode int `json:"errcode"`
	ErrMsg int `json:"errmsg"`
}

type UserLogin struct {
	Code string `json:"code"`
	UserInfo WxUser `json:"userInfo"`
}

type WxUser struct {
	ID int `json:"id" gorm:"column:id"`
	NickName  string `json:"nickName" gorm:"column:nickName"`
	Gender    int    `json:"gender" gorm:"column:gender"`
	Language  string `json:"language" gorm:"column:language"`
	City      string `json:"city" gorm:"column:city"`
	Province  string `json:"province" gorm:"column:province"`
	Country   string `json:"country" gorm:"column:country"`
	AvatarUrl string `json:"avatarUrl" gorm:"column:avatarUrl"`
	SessionKey string `json:"session_key" gorm:"column:session_key"`
	OpenID string `json:"openid" gorm:"column:openid"`
	LeftSubScribeTimes int `json:"LeftSubScribeTimes" gorm:"column:left_subscribe_times"`
	SettingsData string `json:"SettingsData" gorm:"column:settings_data"`
}

type UserSettings struct {
	OpenID string `json:"OpenID" gorm:"openid"`
	SettingsData string `json:"SettingsData" gorm:"settings_data"`
}

func GetAllSubscribeUsers() (ResultUser []string, err error) {
	err = db.Table("mp_user").Select("openid").Where("left_subscribe_times > ?", 0).Find(&ResultUser).Error
	return ResultUser, err
}

func InsertUser(myWxUser WxUser) (ResultUser WxUser, err error) {
	var OldUser WxUser

	db.Table("mp_user").Where("openid =?", myWxUser.OpenID).Find(&OldUser)

	// 新用户
	if OldUser.ID == 0 {
		err = db.Table("mp_user").Create(&myWxUser).Error
	} else {
		// 老用户
		db.Table("mp_user").Select("session_key", "nickName", "gender", "city", "province", "country", "avatarUrl").Where("openid = ?", myWxUser.OpenID).Updates(WxUser{Country: myWxUser.Country, Province: myWxUser.Province, City: myWxUser.City, Language: myWxUser.Language, Gender: myWxUser.Gender, NickName: myWxUser.NickName, AvatarUrl: myWxUser.AvatarUrl, SessionKey: myWxUser.SessionKey, OpenID: myWxUser.OpenID})
	}

	db.Table("mp_user").Where("openid =?", myWxUser.OpenID).Find(&ResultUser)

	return ResultUser, err
}

func DeleteSubscribeState(myUserSettings UserSettings) (finalUserSettings WxUser, err error) {
	err = db.Raw("UPDATE `mp_user` SET `left_subscribe_times` = IF(`left_subscribe_times`<1, 0, `left_subscribe_times`-1) WHERE `openid` = '" + myUserSettings.OpenID + "'").Scan(&finalUserSettings).Error
	err = db.Table("mp_user").Where("openid =?", myUserSettings.OpenID).Find(&finalUserSettings).Error
	return finalUserSettings, err
}

func AddSubscribeTimes(myUserSettings UserSettings) (finalUserSettings WxUser, err error) {
	err = db.Raw("UPDATE `mp_user` SET `left_subscribe_times` = `left_subscribe_times`+1 WHERE `openid` = '" + myUserSettings.OpenID + "'").Scan(&finalUserSettings).Error
	err = db.Table("mp_user").Where("openid =?", myUserSettings.OpenID).Find(&finalUserSettings).Error
	return finalUserSettings, err
}

func CheckSubscribeState(myUserSettings UserSettings) (finalUserSettings WxUser, err error) {
	err = db.Table("mp_user").Select("left_subscribe_times", "settings_data").Where("openid =?", myUserSettings.OpenID).Find(&finalUserSettings).Error
	return finalUserSettings, err
}

func ResetSubscribeTimes(myUserSettings UserSettings) (finalUserSettings WxUser, err error) {
	err = db.Raw("UPDATE `mp_user` SET `left_subscribe_times` = 0 WHERE `openid` = '" + myUserSettings.OpenID + "'").Scan(&finalUserSettings).Error
	return SetUserSettings(myUserSettings)
}

func GetUserSettings(myUserSettings UserSettings) (finalUserSettings WxUser, err error) {
	err = db.Table("mp_user").Where("openid =?", myUserSettings.OpenID).Find(&finalUserSettings).Error
	return finalUserSettings, err
}

func SetUserSettings(myUserSettings UserSettings) (finalUserSettings WxUser, err error) {
	err = db.Table("mp_user").Where("openid =?", myUserSettings.OpenID).Updates(WxUser{OpenID: myUserSettings.OpenID, SettingsData: myUserSettings.SettingsData}).Error
	err = db.Table("mp_user").Where("openid =?", myUserSettings.OpenID).Find(&finalUserSettings).Error
	return finalUserSettings, err
}

func GetAccessToken() TokenResult {
	var myTokenResult TokenResult

	// 导入配置
	AppConfig, err := Config.GetAppConfig()
	if err != nil {
		log.Println("app config err")
	}

	access_token, _ := GetSettings("access_token")
	access_token_expired, _ := GetSettings("access_token_expired")
	access_token_start, _ := GetSettings("access_token_start")
	Int_access_token_expired, _ := strconv.Atoi(access_token_expired.SettingsContent)
	Int_access_token_start, _ := strconv.Atoi(access_token_start.SettingsContent)

	strNowTime := strconv.FormatInt(time.Now().Unix(), 10)
	IntnowTime, _ := strconv.Atoi(strNowTime)

	// 没过期
	if ((Int_access_token_start + Int_access_token_expired) > IntnowTime) {
		myTokenResult.AccessToken = access_token.SettingsContent
		myTokenResult.ExpiresIn, _ = strconv.Atoi(access_token_expired.SettingsContent)
		return myTokenResult
		// 过期了
	} else {
		TokenURL := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + AppConfig.App.APPID + "&secret=" + AppConfig.App.APPSecret

		client := req.C()
		request, err := client.R().SetResult(&myTokenResult).Get(TokenURL)

		if err != nil {}
		if request.IsSuccess() && myTokenResult.ErrCode == 0 {
			SetSettings(Settings{SettingsItem: "access_token", SettingsContent: myTokenResult.AccessToken})
			SetSettings(Settings{SettingsItem: "access_token_expired", SettingsContent: strconv.Itoa(myTokenResult.ExpiresIn)})
			SetSettings(Settings{SettingsItem: "access_token_start", SettingsContent: strNowTime})
			return myTokenResult
		} else {
			return myTokenResult
		}
	}
}