package Model

import (
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

type News struct {
	ID int `gorm:"column:id"`
	UnionID string `gorm:"column:union_id"`
	NewsTitle string `gorm:"column:news_title"`
	NewsSummary string `gorm:"column:news_summary"`
	NewsCoverURL string `gorm:"column:news_cover_url"`
	NewsContent string `gorm:"column:news_content"`
	NewsType string `gorm:"column:news_type"`
	NewsAddTime time.Time `gorm:"column:news_add_time"`
	NewsEditTime time.Time `gorm:"column:news_edit_time"`
	NewsSortTime time.Time `gorm:"column:news_sort_time"`
	// 0 值时 gorm 不给更新，所以要 string
	NewsIsShow string `gorm:"column:news_is_show"`
	NewsIsDelete string `gorm:"column:news_is_delete"`
}

type NewsData struct {
	Data []News `binding:"required,dive"`
}

//func GetNewsList(Start string, Limit string, Type string) (Results []News, DataTotal int, AllTotal int64, IsShowTotal int64, IsHideTotal int64, DeleteTotal int64, err error) {
//
//	IntStart, _ := strconv.Atoi(Start)
//	IntLimit, _ := strconv.Atoi(Limit)
//
//	switch Type {
//		case "All":
//			err = db.Table("mp_news").Where("news_is_delete = ?", FALSE).Limit(IntLimit).Offset(IntStart).Order("news_sort_time desc").Order("news_add_time desc").Find(&Results).Error
//		case "Show":
//			err = db.Table("mp_news").Where("news_is_delete = ? AND news_is_show = ?", FALSE, TRUE).Limit(IntLimit).Offset(IntStart).Order("news_sort_time desc").Order("news_add_time desc").Find(&Results).Error
//		case "Hide":
//			err = db.Table("mp_news").Where("news_is_delete = ? AND news_is_show = ?", FALSE, FALSE).Limit(IntLimit).Offset(IntStart).Order("news_sort_time desc").Order("news_add_time desc").Find(&Results).Error
//		case "Delete":
//			err = db.Table("mp_news").Where("news_is_delete = ?", TRUE).Limit(IntLimit).Offset(IntStart).Order("news_sort_time desc").Order("news_add_time desc").Find(&Results).Error
//	}
//
//	// AllTotal
//	err = db.Table("mp_news").Where("news_is_delete = ?", FALSE).Count(&AllTotal).Error
//	// ISShowTotal
//	err = db.Table("mp_news").Where("news_is_delete = ? AND news_is_show = ?", FALSE, TRUE).Count(&IsShowTotal).Error
//	// ISHideTotal
//	err = db.Table("mp_news").Where("news_is_delete = ? AND news_is_show = ?", FALSE, FALSE).Count(&IsHideTotal).Error
//	// ISDeleteTotal
//	err = db.Table("mp_news").Where("news_is_delete = ?", TRUE).Count(&DeleteTotal).Error
//
//	DataTotal = len(Results)
//
//	return Results, DataTotal, AllTotal, IsShowTotal, IsHideTotal, DeleteTotal, err
//}

func SearchNewsList(Start string, Limit string, Type string, KeyWord string) (Results []News, DataTotal int64, AllTotal int64, IsShowTotal int64, IsHideTotal int64, DeleteTotal int64, err error) {
	intStart, _ := strconv.Atoi(Start)
	intLimit, _ := strconv.Atoi(Limit)
	tableName := "mp_news"

	concatColumns := "CONCAT_WS (news_title, news_summary, news_content, news_type)"

	mpTable := db.Table(tableName)
	// 没有 KeyWords 的情况
	if KeyWord != "" {
		mpTable = mpTable.Where(concatColumns + "like ?", "%" + KeyWord + "%")
	}

	switch Type {
		case "All":
			mpTable = db.Table(tableName).Where("news_is_delete = ?", FALSE)
			break
		case "Show":
			mpTable = db.Table(tableName).Where("news_is_delete = ? AND news_is_show = ?", FALSE, TRUE)
			break
		case "Hide":
			mpTable = db.Table(tableName).Where("news_is_delete = ? AND news_is_show = ?", FALSE, FALSE)
			break
		case "Delete":
			mpTable = db.Table(tableName).Where("news_is_delete = ?", TRUE)
			break
	}

	err = mpTable.Limit(intLimit).Offset(intStart).Order("news_sort_time desc").Order("news_add_time desc").Find(&Results).Error

	// AllTotal
	err = db.Table(tableName).Where("news_is_delete = ?", FALSE).Count(&AllTotal).Error
	// ISShowTotal
	err = db.Table(tableName).Where("news_is_delete = ? AND news_is_show = ?", FALSE, TRUE).Count(&IsShowTotal).Error
	// ISHideTotal
	err = db.Table(tableName).Where("news_is_delete = ? AND news_is_show = ?", FALSE, FALSE).Count(&IsHideTotal).Error
	// ISDeleteTotal
	err = db.Table(tableName).Where("news_is_delete = ?", TRUE).Count(&DeleteTotal).Error
	// DataTotal
	DataTotal = int64(len(Results))

	return Results, DataTotal, AllTotal, IsShowTotal, IsHideTotal, DeleteTotal, err
}

// 删除商品，数组
func EmptyNews() (err error) {
	var myNews News

	err = db.Table("mp_news").Delete(&myNews, "news_is_delete = ?", TRUE).Error

	return err
}

func SetNewsIsDelete(IDs []int, IsDelete string) (finalNews []News, err error) {

	if IsDelete == TRUE {
		// NewsIsDelete 字段设置为 true, IsShow 设置为 False
		err = db.Table("mp_news").Where("id IN ?", IDs).Updates(News{NewsIsDelete: TRUE, NewsIsShow: FALSE}).Error
		db.Table("mp_news").Find(&finalNews, IDs)
	} else if IsDelete == FALSE {
		// NewsIsDelete 字段设置为 false, IsShow 依然设置为 False（总不能刚恢复就显示吧）
		err = db.Table("mp_news").Where("id IN ?", IDs).Updates(News{NewsIsDelete: FALSE, NewsIsShow: FALSE}).Error
		db.Table("mp_news").Find(&finalNews, IDs)
	}

	return finalNews, err
}

// 添加商品
func UpdateNews(myNews News) (finalNews News, err error)  {
	// -1 为新增，0 为异常，其他则需要修改
	if myNews.ID == -1 {
		myNews.UnionID = uuid.NewV4().String()
		// 新的加入
		db.Table("mp_news").Select("union_id", "news_title", "news_summary", "news_cover_url", "news_content", "news_type").Create(&myNews)

		// 返回刚添加的这个商品，此处没有 union id 就返回不了了
		err = db.Table("mp_news").Where("union_id = ?", myNews.UnionID).First(&finalNews).Error
	} else if myNews.ID == 0 {
		return finalNews, err
	} else {
		// 用收到的覆盖掉查询出来的
		db.Table("mp_news").Select("news_title", "news_summary", "news_cover_url", "news_content", "news_type").Where("id = ?", myNews.ID).Updates(&myNews)

		// 返回覆盖的，传过来的时候没有 union_id，所以得用 id 查
		err = db.Table("mp_news").Where("id = ?", myNews.ID).First(&finalNews).Error
	}

	return finalNews, err
}

// 用 ID 获取商品
func GetNewsByUnionID(UnionID string) (finalNews News, err error)  {
	db.Table("mp_news").Where("union_id = ?", UnionID).First(&finalNews)
	return finalNews, err
}

func SetNewsIsTop(ID string, IsTop string) (finalNews News, err error) {
	// 置顶则加入 NewsSortTime 时间，反之则清空
	if IsTop == FALSE {
		// 如果不置顶，就把时间清空
		db.Table("mp_news").Where("id = ?", ID).Update("news_sort_time", nil)
	} else if IsTop == TRUE {
		// 如果置顶，就把隐藏给关掉
		NewsSortTime := time.Now().Format(TIME_LAYOUT)
		db.Debug().Table("mp_news").Where("id = ?", ID).Updates(map[string]interface{}{"news_sort_time": NewsSortTime, "news_is_show": TRUE})
	} else {
		// 参数错误
		return finalNews, nil
	}
	//返回刚添加这个用户
	err = db.Table("mp_news").Where("id = ?", ID).Find(&finalNews).Error
	return finalNews, err
}


func SetNewsIsShow(IDs []int, IsShow string) (finalNews []News, err error) {

	// 如果隐藏，就把置顶给关掉
	if IsShow == FALSE {
		db.Table("mp_news").Where("id IN ?", IDs).Updates(map[string]interface{}{"news_sort_time": nil, "news_is_show": FALSE})
	} else if IsShow == TRUE {
		// 如果显示，也默认不开启置顶
		db.Table("mp_news").Where("id IN ?", IDs).Updates(map[string]interface{}{"news_sort_time": nil, "news_is_show": TRUE})
	}

	//返回刚添加这个用户
	err = db.Table("mp_news").Find(&finalNews, IDs).Error
	return finalNews, err
}