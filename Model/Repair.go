package Model

import (
	"github.com/imroc/req/v3"
	"log"

	//"github.com/imroc/req/v3"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"strconv"
	"time"
)

type Thing1 struct {
	Value string `json:"value"`
}
type Phrase4 struct {
	Value string `json:"value"`
}
type Time7 struct {
	Value string `json:"value"`
}

type Data9283 struct {
	Thing1 Thing1 `json:"thing1"`
	Phrase4 Phrase4 `json:"phrase4"`
	Time7 Time7 `json:"time7"`
}

type SubScribeMsgBasic struct {
	Touser string `json:"touser"`
	TemplateID string `json:"template_id"`
	Page string `json:"page"`
	MiniprogramState string `json:"miniprogram_state"`
	Lang string `json:"lang"`
	Data map[string]map[string]interface{} `json:"data"`
}

type TokenSuccessResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn string `json:"expires_in"`
}

type FetchOrder struct {
	ID int `json:"id"`
	RepairID string `json:"repair_id"`
	OrderID string `json:"order_id"`
}

type Repair struct {
	ID int `gorm:"column:id"`
	RepairID string `gorm:"column:repair_id"`
	OrderID string `gorm:"column:order_id"`
	OrderCoverURL string `gorm:"column:order_cover_url"`
	OrderProductName string `gorm:"column:order_product_name"`
	OrderProductID string `gorm:"column:order_product_id"`
	OrderItem string `gorm:"column:order_item"`
	OrderUserMessage string `gorm:"column:order_user_message"`
	RepairStatus string `gorm:"column:repair_status"`
	TicketCode string `gorm:"column:ticket_code"`
	UserOpenId string `gorm:"column:user_open_id"`
	// 0 值时 gorm 不给更新，所以要 string
	UserName string `gorm:"column:user_name"`
	UserAddress string `gorm:"column:user_address"`
	UserTrackNum string `gorm:"column:user_track_num"`
	UserTrackTime time.Time `gorm:"column:user_track_time"`
	RepairIsDelete string `gorm:"column:repair_is_delete"`
	RepairAddress string `gorm:"column:repair_address"`
	RepairTrackNum string `gorm:"column:repair_track_num"`
	RepairTrackTime time.Time `gorm:"column:repair_track_time"`
	RepairCreateTime time.Time `gorm:"column:repair_create_time"`
}

type OrderDetail struct {
	Errcode int `json:"errcode"`
	Order   struct {
		OrderID     int64  `json:"order_id"`
		CreateTime  string `json:"create_time"`
		UpdateTime  string `json:"update_time"`
		Status      int    `json:"status"`
		OrderDetail struct {
			ProductInfos []struct {
				ProductID int    `json:"product_id"`
				SkuID     int    `json:"sku_id"`
				ThumbImg  string `json:"thumb_img"`
				SalePrice int    `json:"sale_price"`
				SkuCnt    int    `json:"sku_cnt"`
				Title     string `json:"title"`
				SkuAttrs  []struct {
					AttrKey   string `json:"attr_key"`
					AttrValue string `json:"attr_value"`
				} `json:"sku_attrs"`
				OnAftersaleSkuCnt     int    `json:"on_aftersale_sku_cnt"`
				FinishAftersaleSkuCnt int    `json:"finish_aftersale_sku_cnt"`
				SkuCode               string `json:"sku_code"`
				MarketPrice           int    `json:"market_price"`
			} `json:"product_infos"`
			PayInfo struct {
				PayMethod     string `json:"pay_method"`
				PrepayID      string `json:"prepay_id"`
				PrepayTime    string `json:"prepay_time"`
				PayTime       string `json:"pay_time"`
				TransactionID string `json:"transaction_id"`
			} `json:"pay_info"`
			PriceInfo struct {
				ProductPrice    int `json:"product_price"`
				OrderPrice      int `json:"order_price"`
				Freight         int `json:"freight"`
				ChangeDownPrice int `json:"change_down_price"`
			} `json:"price_info"`
			DeliveryInfo struct {
				AddressInfo struct {
					UserName     string `json:"user_name"`
					PostalCode   string `json:"postal_code"`
					ProvinceName string `json:"province_name"`
					CityName     string `json:"city_name"`
					CountyName   string `json:"county_name"`
					DetailInfo   string `json:"detail_info"`
					NationalCode string `json:"national_code"`
					TelNumber    string `json:"tel_number"`
				} `json:"address_info"`
				DeliveryMethod      string        `json:"delivery_method"`
				ExpressFee          []interface{} `json:"express_fee"`
				DeliveryProductInfo []interface{} `json:"delivery_product_info"`
				DeliverType         string        `json:"deliver_type"`
				OfflineDeliveryTime int           `json:"offline_delivery_time"`
				OfflinePickupTime   int           `json:"offline_pickup_time"`
			} `json:"delivery_info"`
			CouponInfo struct {
				CouponID []interface{} `json:"coupon_id"`
			} `json:"coupon_info"`
			CouponcodeInfo struct {
				StartTime        int    `json:"start_time"`
				EndTime          int    `json:"end_time"`
				VerifyType       int    `json:"verify_type"`
				PhoneNumber      string `json:"phone_number"`
				VerifierNickname string `json:"verifier_nickname"`
				VerifyTime       int    `json:"verify_time"`
			} `json:"couponcode_info"`
		} `json:"order_detail"`
		AftersaleDetail struct {
			AftersaleOrderList  []interface{} `json:"aftersale_order_list"`
			OnAftersaleOrderCnt int           `json:"on_aftersale_order_cnt"`
		} `json:"aftersale_detail"`
		Openid  string `json:"openid"`
		ExtInfo struct {
			CustomerNotes string `json:"customer_notes"`
			MerchantNotes string `json:"merchant_notes"`
		} `json:"ext_info"`
		OrderType int `json:"order_type"`
	} `json:"order"`
}

type TextMessage struct {
	Msgtype string `json:"msgtype"`
	Text struct {
		Content string `json:"content"`
		MentionedList []string `json:"mentioned_list"`
		MentionedMobileList []string `json:"mentioned_mobile_list"`
	} `json:"text"`
}

func FetchOrderInfoByOrderID(fetchOrder FetchOrder, myOrderDetail OrderDetail) (resultRepair Repair, err error) {
	var ItemStr string

	// 待收货
	resultRepair.RepairStatus = "UnArrive"
	resultRepair.OrderID = fetchOrder.OrderID
	resultRepair.OrderProductName = myOrderDetail.Order.OrderDetail.ProductInfos[0].Title
	skuLen := len(myOrderDetail.Order.OrderDetail.ProductInfos[0].SkuAttrs)
	for index, value := range myOrderDetail.Order.OrderDetail.ProductInfos[0].SkuAttrs {
		if index != skuLen - 1 {
			ItemStr = ItemStr + value.AttrKey + value.AttrValue + "\n"
		} else {
			ItemStr = ItemStr + value.AttrKey + value.AttrValue
		}
	}

	resultRepair.OrderItem =  ItemStr
	resultRepair.OrderUserMessage = myOrderDetail.Order.ExtInfo.CustomerNotes
	resultRepair.RepairID = fetchOrder.RepairID
	resultRepair.OrderCoverURL = myOrderDetail.Order.OrderDetail.ProductInfos[0].ThumbImg
	resultRepair.ID = fetchOrder.ID
	resultRepair.OrderProductID = strconv.Itoa(myOrderDetail.Order.OrderDetail.ProductInfos[0].ProductID)

	db.Table("mp_repair").Select("ID", "repair_id", "order_cover_url", "order_id", "repair_status", "order_product_name", "order_product_id", "order_cover_url","order_item", "order_user_message").Where("repair_id = ?", fetchOrder.RepairID).Updates(&resultRepair)

	return resultRepair, err
}

func SetRepairTrackNumByRepairID(RepairID string, RepairStatus string, RepairTrackNum string) (resultRepair Repair, err error)  {
	resultRepair.RepairStatus = RepairStatus
	resultRepair.RepairTrackNum = RepairTrackNum

	db.Table("mp_repair").Select("repair_id", "repair_status", "repair_track_num").Where("repair_id = ?", RepairID).Updates(&resultRepair)

	return resultRepair, err
}

func SetRepairStatusByRepairID(RepairID string, RepairStatus string, RepairTrackNum string) (resultRepair Repair, err error)  {
	repairCase, err := GetCaseByRepairID(RepairID)

	resultRepair.RepairStatus = RepairStatus
	resultRepair.RepairID = RepairID

	switch RepairStatus {
		case "UnRepair":
			db.Table("mp_repair").Select("repair_id", "repair_status").Where("repair_id = ?", RepairID).Updates(&resultRepair)

			myData := map[string]map[string]interface{}{
				"thing1":{
					"value":repairCase.OrderProductName,
				},
				"phrase4":{
					"value":"产品已收到",
				},
				"time7":{
					"value":time.Now().Format("2006-01-02 15:04:05"),
				},
			}

			PushMsg(repairCase.RepairID, repairCase.UserOpenId, "yyLrF5cJVb_O-jk56S_ZDi_4uMlx2zn9MRPHfDEKJUs", myData)
			break
		case "Repairing":
			db.Table("mp_repair").Select("repair_id", "repair_status").Where("repair_id = ?", RepairID).Updates(&resultRepair)

			myData := map[string]map[string]interface{}{
				"character_string1":{
					"value":repairCase.RepairID,
				},
				"phrase3":{
					"value":"产品维修中",
				},
				"time2":{
					"value":time.Now().Format("2006-01-02 15:04:05"),
				},
			}

			PushMsg(repairCase.RepairID, repairCase.UserOpenId, "r5ZIYimaQ2HXcYS_tnGhQbGuTOA-bE715d5s3xbIVTs", myData)
			break
		case "Finish":
			resultRepair.RepairTrackNum = RepairTrackNum
			resultRepair.RepairTrackTime = time.Now()
			db.Table("mp_repair").Select("repair_id", "repair_status", "repair_track_num", "repair_track_time").Where("repair_id = ?", RepairID).Updates(&resultRepair)

			myData := map[string]map[string]interface{}{
				"character_string1":{
					"value":repairCase.RepairID,
				},
				"thing2":{
					"value":repairCase.OrderProductName,
				},
				"character_string4":{
					"value":RepairTrackNum,
				},
				"date10":{
					"value":resultRepair.RepairTrackTime.Format("2006-01-02 15:04:05"),
				},
			}

			PushMsg(repairCase.RepairID, repairCase.UserOpenId, "xozYOsDp1q_ShUGMbbbmi2yzdS_7laMzihY9pxSxeyY", myData)
			break
		case "PairFail":
			myData := map[string]map[string]interface{}{
				"character_string1":{
					"value":repairCase.RepairID,
				},
				"phrase3":{
					"value":"请修改券码",
				},
				"time2":{
					"value":time.Now().Format("2006-01-02 15:04:05"),
				},
			}
			db.Table("mp_repair").Select("repair_id", "repair_status").Where("repair_id = ?", RepairID).Updates(&resultRepair)

			PushMsg(repairCase.RepairID, repairCase.UserOpenId, "r5ZIYimaQ2HXcYS_tnGhQbGuTOA-bE715d5s3xbIVTs", myData)
		default:
			db.Table("mp_repair").Select("repair_id", "repair_status").Where("repair_id = ?", RepairID).Updates(&resultRepair)
			break
	}

	return resultRepair, err
}

func PushMsg(repairID string, UserOpenId string, TemplateID string, myData map[string]map[string]interface{}) {
	var mySubScribeMsgBasic SubScribeMsgBasic

	mySubScribeMsgBasic.Touser = UserOpenId
	mySubScribeMsgBasic.TemplateID = TemplateID
	mySubScribeMsgBasic.MiniprogramState = "formal"
	mySubScribeMsgBasic.Data = myData
	mySubScribeMsgBasic.Page = "/pages/repair_info/repair_info?id=" + repairID + "&scene=" + "subMsg"

	token := GetAccessToken()
	push_URL := "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=" + token.AccessToken
	client := req.C()
	resp, err := client.R().SetBody(mySubScribeMsgBasic).Post(push_URL)
	if resp.IsSuccess() == true {
		log.Println("is success")
	}
	if err != nil {}
}

// 删除商品，数组
func EmptyRepairCase() (err error) {
	var myRepair Repair

	err = db.Table("mp_repair").Delete(&myRepair, "repair_is_delete = ?", TRUE).Error

	return err
}

// 用 Repair_ID 获取案例
func SetRepairCaseIsDelete(IDs []int, IsDelete string) (err error)  {
	if IsDelete == TRUE {
		err = db.Table("mp_repair").Where("id IN ?", IDs).Updates(Repair{RepairIsDelete: TRUE}).Error
	} else if IsDelete == FALSE {
		err = db.Table("mp_repair").Where("id IN ?", IDs).Updates(Repair{RepairIsDelete: FALSE}).Error
	}
	return err
}

// 用 ID 获取案例
func GetCaseByRepairID(RepairID string) (resultRepair Repair, err error)  {
	err = db.Table("mp_repair").Where("repair_id = ?", RepairID).First(&resultRepair).Error
	return resultRepair, err
}

// 用 UserTrackNum 获取案例
func GetCaseByTrackNum(TrackNum string) (resultRepair Repair, err error)  {
	err = db.Table("mp_repair").Where("user_track_num = ?", TrackNum).First(&resultRepair).Error
	return resultRepair, err
}

func GetRepairListByOpenIDAndProductID(Start string, Limit string, OpenID string, ProductID string) (Results []Repair, DataTotal int64, err error) {
	IntStart, _ := strconv.Atoi(Start)
	IntLimit, _ := strconv.Atoi(Limit)

	err = db.Table("mp_repair").Where("user_open_id = ? AND order_product_id =? ", OpenID, ProductID).Limit(IntLimit).Offset(IntStart).Order("repair_create_time desc").Find(&Results).Count(&DataTotal).Error

	return Results, DataTotal, err
}

func GetRepairListByOpenID(Start string, Limit string, OpenID string) (Results []Repair, DataTotal int64, err error) {

	IntStart, _ := strconv.Atoi(Start)
	IntLimit, _ := strconv.Atoi(Limit)

	err = db.Table("mp_repair").Where("user_open_id = ? AND repair_is_delete = ?", OpenID, FALSE).Limit(IntLimit).Offset(IntStart).Order("repair_create_time desc").Find(&Results).Count(&DataTotal).Error

	return Results, DataTotal, err
}

func UpdateRepair(myRepair Repair) (finalRepair Repair, err error)  {
	repairId, err := gonanoid.Generate("1234567890", 9)
	// 加上前缀
	repairId = "ZY-" + repairId

	// -1 为新增，0 为异常，其他则需要修改
	if myRepair.ID == -1 {
		myRepair.RepairID = repairId
		mySettings, _ := GetSettings(("repair_address"))
		myRepair.RepairAddress = mySettings.SettingsContent
		myRepair.RepairStatus = "UnPaired"
		// 新的加入
		db.Table("mp_repair").Select("repair_id", "repair_status", "user_open_id", "ticket_code", "user_track_num", "user_name", "user_address", "repair_address").Create(&myRepair)

		PushURL := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=a82fbd20-2181-437d-b15b-9ddeb36dedbe"
		var myMsg TextMessage
		myMsg.Msgtype = "text"
		myMsg.Text.Content = myRepair.RepairID
		//jsons, err := json.Marshal(myMsg)
		req.Post(PushURL)

		// 返回刚添加的，此处没有 union id 就返回不了了
		err = db.Table("mp_repair").Where("repair_id = ?", myRepair.RepairID).First(&finalRepair).Error
	} else if myRepair.ID == 0 {
		return finalRepair, err
	} else {
		// 重新设置券码，状态记得变一下，其他不要变，给手机小程序用的
		if myRepair.TicketCode != "" {
			myRepair.RepairStatus = "UnPaired"
			db.Table("mp_repair").Select("repair_status", "ticket_code").Where("id = ?", myRepair.ID).Updates(Repair{RepairStatus:  myRepair.RepairStatus, TicketCode: myRepair.TicketCode})
		} else {
			// 不设置券码，设置其他信息
			db.Table("mp_repair").Select("order_id", "user_track_num", "user_address", "repair_address", "repair_track_num", "order_user_message").Where("id = ?", myRepair.ID).Updates(&myRepair)
		}

		// 返回覆盖的，传过来的时候没有 union_id，所以得用 id 查
		err = db.Table("mp_repair").Where("id = ?", myRepair.ID).First(&finalRepair).Error
	}

	return finalRepair, err
}

func SearchRepairList(Start string, Limit string, Type string, KeyWord string) (Results []Repair, DataTotal int64, AllTotal int64, UnPairedTotal int64, UnArriveTotal int64, UnRepairTotal int64, RepairingTotal int64, FinishTotal int64, DeleteTotal int64, err error) {
	intStart, _ := strconv.Atoi(Start)
	intLimit, _ := strconv.Atoi(Limit)
	tableName := "mp_repair"

	concatColumns := "CONCAT_WS (repair_id, order_id, order_product_name, order_item, ticket_code, user_name, user_track_num, user_address, repair_track_num)"

	mpTable := db.Table(tableName)
	// 没有 KeyWords 的情况
	if KeyWord != "" {
		mpTable = mpTable.Where(concatColumns + "like ?", "%" + KeyWord + "%")
	}

	// All 是管理端获取全部任何状态
	// RepairAll 只获取核验完毕的
	// UnPaired 未绑定的
	// UnArrive 待收货
	// UnRepair 待维修
	// Repairing 维修中
	// Finish 完成维修
	// PairFail 绑定失败

	switch Type {
		case "All":
			break
		case "RepairAll":
			mpTable = db.Table(tableName).Where("repair_status != ? AND repair_status != ? AND repair_is_delete = ?", "UnPaired", "PairFail", FALSE)
			break
		case "UnPaired":
			mpTable = db.Table(tableName).Where("repair_status = ? AND repair_is_delete = ?", "UnPaired", FALSE)
			break
		case "UnArrive":
			mpTable = db.Table(tableName).Where("repair_status = ? AND repair_is_delete = ?", "UnArrive", FALSE)
			break
		case "UnRepair":
			mpTable = db.Table(tableName).Where("repair_status = ? AND repair_is_delete = ?", "UnRepair", FALSE)
			break
		case "Repairing":
			mpTable = db.Table(tableName).Where("repair_status = ? AND repair_is_delete = ?", "Repairing", FALSE)
			break
		case "Finish":
			mpTable = db.Table(tableName).Where("repair_status = ? AND repair_is_delete = ?", "Finish", FALSE)
			break
		case "Delete":
			mpTable = db.Table(tableName).Where("repair_is_delete = ?", TRUE)
			break
	}

	err = mpTable.Limit(intLimit).Offset(intStart).Order("repair_create_time desc").Find(&Results).Error

	// AllTotal
	err = db.Table(tableName).Count(&AllTotal).Error
	// UnPairedTotal
	err = db.Table(tableName).Where("repair_status = ? AND repair_is_delete = ?", "UnPaired", FALSE).Count(&UnPairedTotal).Error
	// UnArriveTotal
	err = db.Table(tableName).Where("repair_status = ? AND repair_is_delete = ?", "UnArrive", FALSE).Count(&UnArriveTotal).Error
	// UnRepairTotal
	err = db.Table(tableName).Where("repair_status = ? AND repair_is_delete = ?", "UnRepair", FALSE).Count(&UnRepairTotal).Error
	// RepairingTotal
	err = db.Table(tableName).Where("repair_status = ? AND repair_is_delete = ?", "Repairing", FALSE).Count(&RepairingTotal).Error
	// FinishTotal
	err = db.Table(tableName).Where("repair_status = ? AND repair_is_delete = ?", "Finish", FALSE).Count(&FinishTotal).Error
	// DeleteTotal
	err = db.Table(tableName).Where("repair_is_delete = ?", "true").Count(&DeleteTotal).Error
	// DataTotal
	DataTotal = int64(len(Results))

	return Results, DataTotal, AllTotal, UnPairedTotal, UnArriveTotal, UnRepairTotal, RepairingTotal, FinishTotal, DeleteTotal, err
}