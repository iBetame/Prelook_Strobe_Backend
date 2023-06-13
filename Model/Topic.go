package Model

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"log"
	"strconv"
	"time"
)

type Topic struct {
	ID int `gorm:"column:id"`
	UnionID string `gorm:"column:union_id"`
	TopicTitle string `gorm:"column:topic_title"`
	TopicSummary string `gorm:"column:topic_summary"`
	TopicGoods string `gorm:"column:topic_goods"`
	// 0 值时 gorm 不给更新，所以要 string
	TopicIsShow string `gorm:"column:topic_is_show"`
	TopicIsDelete string `gorm:"column:topic_is_delete"`
	TopicSortTime time.Time `gorm:"column:topic_sort_time"`
	TopicAddTime time.Time `gorm:"column:topic_add_time"`
}

type TopicProducts struct {
	Topic
	TopicGoods []TopicProduct
}

// 用 ID 获取商品
func GetTopicByUnionID(UnionID string) (resultTopics TopicProducts, err error)  {
	var myTopic Topic
	var goodsList []interface{}
	var myProducts []TopicProduct

	db.Table("mp_topic").Where("union_id = ?", UnionID).First(&myTopic)

	err = json.Unmarshal([]byte(myTopic.TopicGoods), &goodsList)

	if err != nil {
		log.Println(err)
	}

	db.Table("mp_products").Where("union_id in ?", goodsList).Find(&myProducts)

	// 赋值给 Result
	resultTopics.Topic = myTopic
	resultTopics.TopicGoods = myProducts

	return resultTopics, err
}

// 添加话题
func UpdateTopic(myTopic Topic) (finalTopic Topic, err error)  {
	// -1 为新增，0 为异常，其他则需要修改
	if myTopic.ID == -1 {
		myTopic.UnionID = uuid.NewV4().String()
		// 新的加入
		db.Table("mp_topic").Select("union_id", "topic_title", "topic_summary", "topic_goods").Create(&myTopic)

		// 返回刚添加的这个话题，此处没有 union id 就返回不了了
		err = db.Table("mp_topic").Where("union_id = ?", myTopic.UnionID).First(&finalTopic).Error
	} else if myTopic.ID == 0 {
		return finalTopic, err
	} else {
		// 用收到的覆盖掉查询出来的
		db.Table("mp_topic").Select( "topic_title", "topic_summary", "topic_goods").Where("id = ?", myTopic.ID).Updates(&myTopic)

		// 返回覆盖的，传过来的时候没有 union_id，所以得用 id 查
		err = db.Table("mp_topic").Where("id = ?", myTopic.ID).First(&finalTopic).Error
	}

	return finalTopic, err
}


func SearchTopicList(Start string, Limit string, Type string, KeyWord string) (Results []Topic, DataTotal int64, AllTotal int64, IsShowTotal int64, IsHideTotal int64, DeleteTotal int64, err error) {
	intStart, _ := strconv.Atoi(Start)
	intLimit, _ := strconv.Atoi(Limit)
	tableName := "mp_topic"

	concatColumns := "CONCAT_WS (topic_title, topic_summary, topic_goods)"

	mpTable := db.Table(tableName)
	// 没有 KeyWords 的情况
	if KeyWord != "" {
		mpTable = mpTable.Where(concatColumns + "like ?", "%" + KeyWord + "%")
	}

	switch Type {
		case "All":
			mpTable = db.Table(tableName).Where("topic_is_delete = ?", FALSE)
		case "Show":
			mpTable = db.Table(tableName).Where("topic_is_delete = ? AND topic_is_show = ?", FALSE, TRUE)
		case "Hide":
			mpTable = db.Table(tableName).Where("topic_is_delete = ? AND topic_is_show = ?", FALSE, FALSE)
		case "Delete":
			mpTable = db.Table(tableName).Where("topic_is_delete = ?", TRUE)
	}

	err = mpTable.Limit(intLimit).Offset(intStart).Order("topic_sort_time desc").Order("topic_add_time desc").Find(&Results).Error

	// AllTotal
	err = db.Table(tableName).Where("topic_is_delete = ?", FALSE).Count(&AllTotal).Error
	// ISShowTotal
	err = db.Table(tableName).Where("topic_is_delete = ? AND topic_is_show = ?", FALSE, TRUE).Count(&IsShowTotal).Error
	// ISHideTotal
	err = db.Table(tableName).Where("topic_is_delete = ? AND topic_is_show = ?", FALSE, FALSE).Count(&IsHideTotal).Error
	// ISDeleteTotal
	err = db.Table(tableName).Where("topic_is_delete = ?", TRUE).Count(&DeleteTotal).Error
	// DataTotal
	DataTotal = int64(len(Results))

	return Results, DataTotal, AllTotal, IsShowTotal, IsHideTotal, DeleteTotal, err
}

// 删除话题，数组
func EmptyTopic() (err error) {
	var myTopic Topic

	err = db.Table("mp_topic").Delete(&myTopic, "topic_is_delete = ?", TRUE).Error

	return err
}

func SetTopicIsDelete(IDs []int, IsDelete string) (finalTopic []Topic, err error) {

	if IsDelete == TRUE {
		// TopicIsDelete 字段设置为 true, IsShow 设置为 False
		err = db.Table("mp_topic").Where("id IN ?", IDs).Updates(Topic{TopicIsDelete: TRUE, TopicIsShow: FALSE}).Error
		db.Table("mp_topic").Find(&finalTopic, IDs)
	} else if IsDelete == FALSE {
		// TopicIsDelete 字段设置为 false, IsShow 依然设置为 False（总不能刚恢复就显示吧）
		err = db.Table("mp_topic").Where("id IN ?", IDs).Updates(Topic{TopicIsDelete: FALSE, TopicIsShow: FALSE}).Error
		db.Table("mp_topic").Find(&finalTopic, IDs)
	}

	return finalTopic, err
}

// 用 ID 获取话题
func GetTopicByID(ID string) (finalTopic Topic, err error)  {
	db.Table("mp_topic").Where("id = ?", ID).First(&finalTopic)
	return finalTopic, err
}

func SetTopicIsTop(ID string, IsTop string) (finalTopic Topic, err error) {
	// 置顶则加入 TopicSortTime 时间，反之则清空
	if IsTop == FALSE {
		// 如果不置顶，就把时间清空
		db.Table("mp_topic").Where("id = ?", ID).Update("topic_sort_time", nil)
	} else if IsTop == TRUE {
		// 如果置顶，就把隐藏给关掉
		TopicSortTime := time.Now().Format(TIME_LAYOUT)
		db.Debug().Table("mp_topic").Where("id = ?", ID).Updates(map[string]interface{}{"topic_sort_time": TopicSortTime, "topic_is_show": TRUE})
	} else {
		// 参数错误
		return finalTopic, nil
	}
	//返回刚添加这个用户
	err = db.Table("mp_topic").Where("id = ?", ID).Find(&finalTopic).Error
	return finalTopic, err
}


func SetTopicIsShow(IDs []int, IsShow string) (finalTopic []Topic, err error) {

	// 如果隐藏，就把置顶给关掉
	if IsShow == FALSE {
		db.Table("mp_topic").Where("id IN ?", IDs).Updates(map[string]interface{}{"topic_sort_time": nil, "topic_is_show": FALSE})
	} else if IsShow == TRUE {
		// 如果显示，也默认不开启置顶
		db.Table("mp_topic").Where("id IN ?", IDs).Updates(map[string]interface{}{"topic_sort_time": nil, "topic_is_show": TRUE})
	}

	//返回刚添加这个用户
	err = db.Table("mp_topic").Find(&finalTopic, IDs).Error
	return finalTopic, err
}