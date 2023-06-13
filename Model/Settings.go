package Model

type Settings struct {
	ID int `json:"ID" gorm:"column:id"`
	SettingsItem string `json:"SettingsItem" binding:"required" gorm:"column:settings_item"`
	SettingsContent string `json:"SettingsContent" binding:"required" gorm:"column:settings_content"`
}

type SettingsData struct {
	Data []Settings `binding:"required,dive"`
}

func GetAllSettings() (Settings []Settings, err error)  {
	err = db.Table("mp_settings").Find(&Settings).Error
	return Settings, err
}

func SetSettings(mySettings Settings) (Result Settings, err error) {

	err = db.Debug().Table("mp_settings").Where("settings_item = ?", mySettings.SettingsItem).Update("settings_content", mySettings.SettingsContent).Error
	db.Table("mp_settings").Where("settings_item = ?", mySettings.SettingsItem).Find(&Result)

	if err != nil {
		return Result, err
	}
	return Result, err
}

// 单独获取设置选项的接口
func GetSettings(SettingsItem string) (ResultSettings Settings, err error) {
	err = db.Table("mp_settings").Where("settings_item = ?", SettingsItem).First(&ResultSettings).Error
	return ResultSettings, err
}