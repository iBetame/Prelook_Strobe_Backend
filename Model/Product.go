package Model

import (
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

type Product struct {
	ID int `gorm:"column:id"`
	UnionID string `gorm:"column:union_id"`
	ProductTitle string `gorm:"column:product_title"`
	ProductType int `gorm:"column:product_type"`
	ProductSummary string `gorm:"column:product_summary"`
	ProductCoverURL string `gorm:"column:product_cover_url"`
	ProductPrice string `gorm:"column:product_price"`
	ProductID int `gorm:"column:product_id"`
	ProductAddTime time.Time `gorm:"column:product_add_time"`
	ProductEditTime time.Time `gorm:"column:product_edit_time"`
	ProductSortTime time.Time `gorm:"column:product_sort_time"`
	// 0 值时 gorm 不给更新，所以要 string
	ProductIsShow string `gorm:"column:product_is_show"`
	ProductIsDelete string `gorm:"column:product_is_delete"`
}

type TopicProduct struct {
	UnionID string `gorm:"column:union_id"`
	ProductTitle string `gorm:"column:product_title"`
	ProductType int `gorm:"column:product_type"`
	ProductSummary string `gorm:"column:product_summary"`
	ProductCoverURL string `gorm:"column:product_cover_url"`
	ProductPrice string `gorm:"column:product_price"`
	ProductID int `gorm:"column:product_id"`
}

type ProductData struct {
	Data []Product `binding:"required,dive"`
}

type OrangeLog struct {
	ID int `gorm:"column:id"`
	TrackNum string `gorm:"column:track_num""`
}

func SearchProductList(Start string, Limit string, Type string, KeyWord string) (Results []Product, DataTotal int64, AllTotal int64, IsShowTotal int64, IsHideTotal int64, DeleteTotal int64, err error) {
	intStart, _ := strconv.Atoi(Start)
	intLimit, _ := strconv.Atoi(Limit)
	tableName := "mp_products"

	concatColumns := "CONCAT_WS (product_title, product_summary, product_price, product_id)"

	mpTable := db.Table(tableName)

	// 没有 KeyWords 的情况
	if KeyWord != "" {
		mpTable = mpTable.Where(concatColumns + "like ?", "%" + KeyWord + "%")
	}

	switch Type {
		case "All":
			mpTable = db.Table(tableName).Where("product_is_delete = ?", FALSE)
			break
		case "Show":
			mpTable = db.Table(tableName).Where("product_is_delete = ? AND product_is_show = ?", FALSE, TRUE)
			break
		case "Hide":
			mpTable = db.Table(tableName).Where("product_is_delete = ? AND product_is_show = ?", FALSE, FALSE)
			break
		case "Delete":
			mpTable = db.Table(tableName).Where("product_is_delete = ?", TRUE)
			break
	}

	err = mpTable.Limit(intLimit).Offset(intStart).Order("product_sort_time desc").Order("product_add_time desc").Find(&Results).Error

	// AllTotal
	err = db.Table(tableName).Where("product_is_delete = ?", FALSE).Count(&AllTotal).Error
	// ISShowTotal
	err = db.Table(tableName).Where("product_is_delete = ? AND product_is_show = ?", FALSE, TRUE).Count(&IsShowTotal).Error
	// ISHideTotal
	err = db.Table(tableName).Where("product_is_delete = ? AND product_is_show = ?", FALSE, FALSE).Count(&IsHideTotal).Error
	// ISDeleteTotal
	err = db.Table(tableName).Where("product_is_delete = ?", TRUE).Count(&DeleteTotal).Error
	// DataTotal
	DataTotal = int64(len(Results))

	return Results, DataTotal, AllTotal, IsShowTotal, IsHideTotal, DeleteTotal, err
}

// 删除商品，数组
func EmptyProduct() (err error) {
	var myProduct Product

	err = db.Table("mp_products").Delete(&myProduct, "product_is_delete = ?", TRUE).Error

	return err
}

func SetProductIsDelete(IDs []int, IsDelete string) (finalProduct []Product, err error) {

	if IsDelete == TRUE {
		// ProductIsDelete 字段设置为 true, IsShow 设置为 False
		err = db.Table("mp_products").Where("id IN ?", IDs).Updates(Product{ProductIsDelete: TRUE, ProductIsShow: FALSE}).Error
		db.Table("mp_products").Find(&finalProduct, IDs)
	} else if IsDelete == FALSE {
		// ProductIsDelete 字段设置为 false, IsShow 依然设置为 False（总不能刚恢复就显示吧）
		err = db.Table("mp_products").Where("id IN ?", IDs).Updates(Product{ProductIsDelete: FALSE, ProductIsShow: FALSE}).Error
		db.Table("mp_products").Find(&finalProduct, IDs)
	}

	return finalProduct, err
}

// 添加商品
func UpdateProduct(myProduct Product) (finalProduct Product, err error)  {

	// -1 为新增，0 为异常，其他则需要修改
	// Product Type 1 是维修商品，0 是普通商品
	if myProduct.ID == -1 {
		myProduct.UnionID = uuid.NewV4().String()
		// 新的加入
		db.Table("mp_products").Select("union_id", "product_type", "product_title", "product_summary", "product_cover_url", "product_price", "product_id").Create(&myProduct)

		// 返回刚添加的这个商品，此处没有 union id 就返回不了了
		err = db.Table("mp_products").Where("union_id = ?", myProduct.UnionID).First(&finalProduct).Error
	} else if myProduct.ID == 0 {
		return finalProduct, err
	} else {
		// 用收到的覆盖掉查询出来的
		db.Table("mp_products").Select("product_title", "product_type", "product_summary", "product_cover_url", "product_price", "product_id").Where("id = ?", myProduct.ID).Updates(&myProduct)

		// 返回覆盖的，传过来的时候没有 union_id，所以得用 id 查
		err = db.Table("mp_products").Where("id = ?", myProduct.ID).First(&finalProduct).Error
	}

	return finalProduct, err
}

// 用 ID 获取商品
func GetProductByID(ID string) (finalProduct Product, err error)  {
	db.Table("mp_products").Where("id = ?", ID).First(&finalProduct)
	return finalProduct, err
}

func SetProductIsTop(ID string, IsTop string) (finalProduct Product, err error) {
	// 置顶则加入 ProductSortTime 时间，反之则清空
	if IsTop == FALSE {
		// 如果不置顶，就把时间清空
		db.Table("mp_products").Where("id = ?", ID).Update("product_sort_time", nil)
	} else if IsTop == TRUE {
		// 如果置顶，就把隐藏给关掉
		ProductSortTime := time.Now().Format(TIME_LAYOUT)
		db.Debug().Table("mp_products").Where("id = ?", ID).Updates(map[string]interface{}{"product_sort_time": ProductSortTime, "product_is_show": TRUE})
	} else {
		// 参数错误
		return finalProduct, nil
	}
	//返回刚添加这个用户
	err = db.Table("mp_products").Where("id = ?", ID).Find(&finalProduct).Error
	return finalProduct, err
}


func SetProductIsShow(IDs []int, IsShow string) (finalProduct []Product, err error) {

	// 如果隐藏，就把置顶给关掉
	if IsShow == FALSE {
		db.Table("mp_products").Where("id IN ?", IDs).Updates(map[string]interface{}{"product_sort_time": nil, "product_is_show": FALSE})
	} else if IsShow == TRUE {
		// 如果显示，也默认不开启置顶
		db.Table("mp_products").Where("id IN ?", IDs).Updates(map[string]interface{}{"product_sort_time": nil, "product_is_show": TRUE})
	}

	//返回刚添加这个用户
	err = db.Table("mp_products").Find(&finalProduct, IDs).Error
	return finalProduct, err
}

func IsTrackNumRepeat(TrackNum string) (myOrangeLog OrangeLog){
	// 找有没有记录过
	db.Table("mp_orange_logs").Where("track_num =? ", TrackNum).First(&myOrangeLog)

	// 没记录过的话
	if myOrangeLog.ID == 0 {
		// 就插入一下
		myOrangeLog.TrackNum = TrackNum
		db.Table("mp_orange_logs").Create(&myOrangeLog)
		// 插入完清空，代表没找到过
		myOrangeLog.ID = 0
	}

	return myOrangeLog
}