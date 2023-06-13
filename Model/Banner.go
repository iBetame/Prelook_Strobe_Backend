package Model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Banner struct {
	BannerID       string     `json:"banner_id"`
	BannerTitle          string     `json:"banner_title" gorm:"not null"`
	PicURL         string     `json:"pic_url" gorm:"not null"`
	State          string     `json:"state" gorm:"default:Draft"`
}

type BannerData struct {
	Data []Banner `binding:"required,dive"`
}


// 定义个全局函数
func bannerTable() *gorm.DB {
	return db.Table("pl_banners")
}

//func SearchBannerRelate(KeyWord string) (Results []Banner, err error)  {
//	var myProducts []Product
//	var myNews []News
//	err = db.Table("mp_news").Raw("SELECT * FROM mp_news WHERE CONCAT (news_title, news_summary) like CONCAT('%" + KeyWord + "%')" + " AND news_is_delete = '" + FALSE  + "' AND news_is_show = '" + TRUE + "' LIMIT 30 OFFSET 0").Order("id").Scan(&myNews).Error
//	err = db.Table("mp_products").Raw("SELECT * FROM mp_products WHERE CONCAT (product_title, product_summary) like CONCAT('%"  + KeyWord + "%')" + " AND product_is_delete = '" + FALSE  + "' AND product_is_show = '" + TRUE + "' LIMIT 30 OFFSET 0").Order("id").Scan(&myProducts).Error
//	for _, value := range myProducts {
//		var tempResult Banner
//		tempResult.BannerType = "PG"
//		tempResult.BannerContent = strconv.Itoa(value.ProductID)
//		tempResult.BannerTitle = value.ProductTitle
//		tempResult.UnionID = value.UnionID
//		tempResult.BannerSummary = value.ProductSummary
//		Results = append(Results, tempResult)
//	}
//	for _, value := range myNews {
//		var tempResult Banner
//		tempResult.BannerTitle = value.NewsTitle
//		tempResult.BannerCoverURL = value.NewsCoverURL
//		tempResult.BannerContent = value.NewsContent
//		tempResult.UnionID = value.UnionID
//		tempResult.BannerSummary = value.NewsSummary
//		tempResult.BannerType = value.NewsType
//		Results = append(Results, tempResult)
//	}
//
//	return Results, err
//}
//
//func SearchBannerList(Start string, Limit string, Type string, KeyWord string) (Results []Banner, DataTotal int64, AllTotal int64, IsShowTotal int64, IsHideTotal int64, DeleteTotal int64, err error) {
//	intStart, _ := strconv.Atoi(Start)
//	intLimit, _ := strconv.Atoi(Limit)
//	tableName := "mp_banner"
//
//	concatColumns := "CONCAT_WS (banner_title, banner_desc, banner_content, banner_type)"
//
//	mpTable := db.Table(tableName)
//	// 没有 KeyWords 的情况
//	if KeyWord != "" {
//		mpTable = mpTable.Where(concatColumns + "like ?", "%" + KeyWord + "%")
//	}
//
//	switch Type {
//		case "All":
//			mpTable = db.Table(tableName).Where("banner_is_delete = ?", FALSE)
//			break
//		case "Show":
//			mpTable = db.Table(tableName).Where("banner_is_delete = ? AND banner_is_show = ?", FALSE, TRUE)
//			break
//		case "Hide":
//			mpTable = db.Table(tableName).Where("banner_is_delete = ? AND banner_is_show = ?", FALSE, FALSE)
//			break
//		case "Delete":
//			mpTable = db.Table(tableName).Where("banner_is_delete = ?", TRUE)
//			break
//	}
//
//	err = mpTable.Limit(intLimit).Offset(intStart).Order("banner_sort_time desc").Order("banner_add_time desc").Find(&Results).Error
//
//	// AllTotal
//	err = db.Table(tableName).Where("banner_is_delete = ?", FALSE).Count(&AllTotal).Error
//	// ISShowTotal
//	err = db.Table(tableName).Where("banner_is_delete = ? AND banner_is_show = ?", FALSE, TRUE).Count(&IsShowTotal).Error
//	// ISHideTotal
//	err = db.Table(tableName).Where("banner_is_delete = ? AND banner_is_show = ?", FALSE, FALSE).Count(&IsHideTotal).Error
//	// ISDeleteTotal
//	err = db.Table(tableName).Where("banner_is_delete = ?", TRUE).Count(&DeleteTotal).Error
//	// DataTotal
//	DataTotal = int64(len(Results))
//
//	return Results, DataTotal, AllTotal, IsShowTotal, IsHideTotal, DeleteTotal, err
//}
//
//// 删除商品，数组
//func EmptyBanner() (err error) {
//	var myBanner Banner
//
//	err = db.Table("mp_banner").Delete(&myBanner, "banner_is_delete = ?", TRUE).Error
//
//	return err
//}
//
//func SetBannerIsDelete(IDs []int, IsDelete string) (finalBanner []Banner, err error) {
//
//	if IsDelete == TRUE {
//		// BannerIsDelete 字段设置为 true, IsShow 设置为 False
//		err = db.Table("mp_banner").Where("id IN ?", IDs).Updates(Banner{BannerIsDelete: TRUE, BannerIsShow: FALSE}).Error
//		db.Table("mp_banner").Find(&finalBanner, IDs)
//	} else if IsDelete == FALSE {
//		// BannerIsDelete 字段设置为 false, IsShow 依然设置为 False（总不能刚恢复就显示吧）
//		err = db.Table("mp_banner").Where("id IN ?", IDs).Updates(Banner{BannerIsDelete: FALSE, BannerIsShow: FALSE}).Error
//		db.Table("mp_banner").Find(&finalBanner, IDs)
//	}
//
//	return finalBanner, err
//}

// 添加商品
func UpdateBanner(myBanner Banner) (finalBanner Banner, err error)  {


	if myBanner.BannerID == "-1" {
		// 如果 BannerID 为空字符串，说明需要新增一个 Banner
		// 在数据库中生成一个新的 BannerID
		myBanner.BannerID = uuid.NewV4().String()

		// 在数据库中创建 Banner
		if err := bannerTable().Create(&myBanner).Error; err != nil {
			return Banner{}, err
		}

		// 返回创建成功的 Banner
		return myBanner, nil
	}

	// 如果 BannerID 不为空字符串，说明需要更新一个已有的 Banner
	// 定义一个变量来存放要更新的 Banner 字段
	var updateBanner Banner
	updateBanner.BannerTitle = myBanner.BannerTitle
	updateBanner.PicURL = myBanner.PicURL
	updateBanner.State = myBanner.State

	// 在数据库中更新 Banner
	if err := bannerTable().Where("banner_id = ?", myBanner.BannerID).Updates(updateBanner).Error; err != nil {
		return Banner{}, err
	}

	// 返回更新成功的 Banner
	return myBanner, nil
}

// 用 ID 获取商品
func GetBannerByID(ID string) (finalBanner Banner, err error)  {
	db.Table("mp_banner").Where("id = ?", ID).First(&finalBanner)
	return finalBanner, err
}
