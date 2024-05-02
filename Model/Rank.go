package Model

import (
	"errors"
	"strconv"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Rank struct {
	ID          int       `json:"id"`
	RankID      string    `json:"rank_id"`
	Name        string    `json:"name"`
	Description string    `json:"des"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	CreateTime  time.Time `json:"createTime" gorm:"column:create_datetime"`
}

type RankingItem struct {
	ID             int       `json:"id"`
	RankID         string    `json:"rank_id"`
	Brand          string    `json:"brand"`
	Name           string    `json:"name"`
	RankingIndex   string    `json:"ranking_index"`
	Score          string    `json:"score"`
	DrawingPoints  string    `json:"drawing_points" gorm:"type:longtext"`
	CreateDateTime time.Time `json:"create_datetime" gorm:"column:create_datetime"`
}

type RankData struct {
	ID          int    `json:"id"`
	RankID      string `json:"rank_id"`
	Name        string `json:"name"`
	Description string `json:"des"`
	Type        string `json:"type"`
	// Rankings    []RankingItem `json:"rankings"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_datetime"`
}

// 定义个全局函数
func rankRankTable() *gorm.DB {
	return db.Debug().Table("ranking")
}

// 定义个全局函数
func rankRankItemTable() *gorm.DB {
	return db.Table("ranking_item")
}

func GetRankByID(id string) (resultData RankData, err error) {
	var rank Rank

	// 查询排行榜
	err = rankRankTable().Where("id = ?", id).First(&rank).Error
	if err != nil {
		return RankData{}, err
	}

	// 查询关联的排名项
	// var rankingItems []RankingItem
	// err = rankRankItemTable().Where("rank_id = ?", rank.RankID).Find(&rankingItems).Error
	// if err != nil {
	// 	return RankData{}, err
	// }

	// 将关联的排名项添加到排行榜中
	resultData.ID = rank.ID
	resultData.RankID = rank.RankID
	resultData.Name = rank.Name
	resultData.Description = rank.Description
	resultData.Type = rank.Type
	resultData.CreateTime = rank.CreateTime
	// resultData.Rankings = rankingItems

	return resultData, nil
}

// UpdateRank 函数用于更新或创建排行榜
// 参数：myRank 为要更新或创建的排行榜
// 返回值：更新或创建成功后的排行榜实例和可能的错误
func UpdateRank(myRank Rank) (resultRank Rank, err error) {
	rankId, err := gonanoid.Generate("1234567890", 9)
	rankId = "PL-" + rankId
	if err != nil {
		return
	}

	if myRank.Status == "Published" {
		var rankCheck Rank
		rankRankTable().Where("status", "Published").Where("type", myRank.Type).First(&rankCheck)

		if rankCheck.ID != 0 && rankCheck.ID != myRank.ID {
			err = errors.New("该产品类别下已存在发布的排行")
			return
		}
	}
	// 如果 myRank.ID 为 nil，，说明需要新增一个排行榜
	if myRank.ID == 0 {
		// 在数据库中创建排行榜
		myRank.RankID = rankId
		myRank.CreateTime = time.Now()
		err = rankRankTable().Create(&myRank).Error
		resultRank = myRank
	} else {
		// 如果 myRank.RankID 不为 nil，说明需要更新一个已有的排行榜
		// 定义一个变量来存放要更新的排行榜字段
		var updateRank Rank
		updateRank.Name = myRank.Name
		updateRank.Type = myRank.Type
		updateRank.Status = myRank.Status
		updateRank.Description = myRank.Description
		rankRankTable().Where("id =? ", myRank.ID).Updates(updateRank)
		err = rankRankTable().Where("id = ?", myRank.ID).First(&resultRank).Error
	}
	return resultRank, err
}

// UpdateRank 函数用于更新或创建排行榜
// 参数：myRank 为要更新或创建的排行榜
// 返回值：更新或创建成功后的排行榜实例和可能的错误
func UpdateRankItem(myRankItem RankingItem) (resultRank RankingItem, err error) {
	// 如果 myRank.ID 为 nil，，说明需要新增一个排行榜
	if myRankItem.ID == 0 {
		// 在数据库中创建排行榜
		myRankItem.RankID = myRankItem.RankID
		myRankItem.CreateDateTime = time.Now()
		err = rankRankItemTable().Create(&myRankItem).Error
		resultRank = myRankItem
	} else {
		// 如果 myRank.RankID 不为 0，说明需要更新一个已有的排行榜
		// 定义一个变量来存放要更新的排行榜字段
		var updateRankItem RankingItem
		updateRankItem.Name = myRankItem.Name
		updateRankItem.RankingIndex = myRankItem.RankingIndex
		updateRankItem.Brand = myRankItem.Brand
		updateRankItem.Score = myRankItem.Score
		updateRankItem.RankID = myRankItem.RankID
		updateRankItem.DrawingPoints = myRankItem.DrawingPoints
		rankRankItemTable().Where("id =? ", myRankItem.ID).Updates(updateRankItem)
		err = rankRankItemTable().Where("id = ?", myRankItem.ID).First(&resultRank).Error
	}
	if err != nil {
		return
	}
	err = rankRankTable().Where("rank_id = ?", myRankItem.RankID).Update("create_datetime", time.Now()).Error
	return resultRank, err
}

//
//func GetRankByID(ID string) (*Rank, error) {
//	var rank Rank
//	result := rankRankTable().Where("id = ?", ID).First(&rank)
//	if result.Error != nil {
//		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//			return nil, nil
//		}
//		return nil, result.Error
//	}
//	return &rank, nil
//}
//
// SearchRankList 搜索排行榜列表
// 参数：
// start int - 起始记录数
// limit int - 每页记录数
// rank_state string - 状态（All - 所有，Published - 已发布，Draft - 草稿，Scheduled - 已排期，Deleted - 已删除）
// rank_mode string - 模式，是 dock 还是 list
// keyword string - 关键词
// 返回值：
// Results []Rank - 排行榜列表
// DataTotal int64 - 数据总数
// AllTotal int64 - 所有状态的数据总数
// PublishedTotal int64 - 已发布状态的数据总数
// DraftTotal int64 - 草稿状态的数据总数
// ScheduledTotal int64 - 已排期状态的数据总数
// DeleteTotal int64 - 已删除状态的数据总数
// err error - 错误信息
// Mode：List：普通元素；Dock：首页元素；Top：置顶元素

type APIRank struct {
	Rank
	ItemCount int `json:"item_count"`
}

func SearchRankList(start string, limit string, status string, keyword string, order string) (resultRank []APIRank, DataTotal int64, AllTotal int64, PublishedTotal int64, DraftTotal int64, RemovedTotal int64, err error) {
	intStart, _ := strconv.Atoi(start)
	intLimit, _ := strconv.Atoi(limit)

	concatColumns := "CONCAT_WS (rank_id, name)"

	mpTable := rankRankTable()

	// 没有 KeyWords 的情况
	if keyword != "" {
		mpTable = mpTable.Where(concatColumns+"like ?", "%"+keyword+"%")
	}

	if order == "asc" {
		order = "asc"
	} else {
		order = "desc"
	}

	switch status {
	case "All":
		mpTable = mpTable.Where("ranking.status != ?", "Removed")
	case "Published":
		mpTable = mpTable.Where("ranking.status = ? AND ranking.is_delete = ?", "Published", FALSE)
		break
	case "Draft":
		mpTable = mpTable.Where("ranking.status = ? AND ranking.is_delete = ?", "Draft", FALSE)
		break
	case "Removed":
		mpTable = mpTable.Where("ranking.status = ? OR ranking.is_delete = ?", "Removed", TRUE)
		break
	}
	// DataTotal
	mpTable.Count(&DataTotal)

	err = mpTable.Select("ranking.*, count(i.id) as item_count").Joins("left join ranking_item i on i.rank_id=ranking.rank_id and i.is_delete = 'false'").
		Group("ranking.id").
		Limit(intLimit).Offset(intStart).Order("create_datetime " + order).
		Find(&resultRank).Error

	// AllTotal
	err = rankRankTable().Count(&AllTotal).Error
	// Published
	err = rankRankTable().Where("status = ? AND is_delete = ?", "Published", FALSE).Count(&PublishedTotal).Error
	// Draft
	err = rankRankTable().Where("status = ? AND is_delete = ?", "Draft", FALSE).Count(&DraftTotal).Error
	// Removed
	err = rankRankTable().Where("status = ? OR is_delete = ?", "Removed", TRUE).Count(&RemovedTotal).Error

	return resultRank, DataTotal, AllTotal, PublishedTotal, DraftTotal, RemovedTotal, err
}

func SearchRankItemList(rank_id string, start int, limit int, status string, keyword string, order string) (resultRank []RankingItem, DataTotal int64, AllTotal int64, PublishedTotal int64, DraftTotal int64, RemovedTotal int64, err error) {
	concatColumns := "CONCAT_WS (brand, name)"

	mpTable := rankRankItemTable().Where("rank_id", rank_id)

	// 没有 KeyWords 的情况
	if keyword != "" {
		mpTable = mpTable.Where(concatColumns+"like ?", "%"+keyword+"%")
	}

	if order == "asc" {
		order = "asc"
	} else {
		order = "desc"
	}

	switch status {
	case "All":
		mpTable = mpTable.Where("is_delete = ?", FALSE)
	case "Published":
		mpTable = mpTable.Where("status = ? AND is_delete = ?", "Published", FALSE)
		break
	case "Draft":
		mpTable = mpTable.Where("status = ? AND is_delete = ?", "Draft", FALSE)
		break
	case "Removed":
		mpTable = mpTable.Where("status = ? OR is_delete = ?", "Removed", TRUE)
		break
	}

	// DataTotal
	mpTable.Count(&DataTotal)

	err = mpTable.Limit(limit).Offset(start).Order("score " + order).Find(&resultRank).Error

	// AllTotal
	err = rankRankItemTable().Count(&AllTotal).Error
	// Published
	err = rankRankItemTable().Where("status = ? AND is_delete = ?", "Published", FALSE).Count(&PublishedTotal).Error
	// Draft
	err = rankRankItemTable().Where("status = ? AND is_delete = ?", "Draft", FALSE).Count(&DraftTotal).Error
	// Removed
	err = rankRankItemTable().Where("status = ? OR is_delete = ?", "Removed", TRUE).Count(&RemovedTotal).Error

	return resultRank, DataTotal, AllTotal, PublishedTotal, DraftTotal, RemovedTotal, err
}

func ChangeRankState(ids []uint, state string) ([]Rank, error) {
	// 批量查询 rank
	var rank []Rank
	if err := rankRankTable().Where("id IN (?)", ids).Find(&rank).Error; err != nil {
		return nil, err
	}

	// 批量更新状态
	err := rankRankTable().Where("id IN (?)", ids).Updates(map[string]interface{}{"status": state}).Error
	if err != nil {
		return nil, err
	}

	// 批量查询并返回结果
	var updatedRanks []Rank
	if err := rankRankTable().Where("id IN (?)", ids).Find(&updatedRanks).Error; err != nil {
		return nil, err
	}
	return updatedRanks, nil
}

func ChangeRankItemDeleteState(ids []uint, state string) ([]Rank, error) {
	// 批量查询 rank item
	var rank []Rank
	if err := rankRankItemTable().Where("id IN (?)", ids).Find(&rank).Error; err != nil {
		return nil, err
	}

	// 批量更新状态
	err := rankRankItemTable().Where("id IN (?)", ids).Updates(map[string]interface{}{"is_delete": state}).Error
	if err != nil {
		return nil, err
	}

	// 批量查询并返回结果
	var updatedRanks []Rank
	if err := rankRankTable().Where("id IN (?)", ids).Find(&updatedRanks).Error; err != nil {
		return nil, err
	}
	return updatedRanks, nil
}

//
//// 删除 Rank
//func EmptyRank() (Results []Rank, DeleteTotal int, err error) {
//	var myRank Rank
//	mpTable := rankRankTable()
//	mpTable = mpTable.Where("rank_state = ?", "Trash")
//	// 找到他们都是谁
//	mpTable.Order("create_datetime desc").Find(&Results)
//	DeleteTotal = len(Results)
//
//	// 删除
//	err = rankRankTable().Delete(&myRank, "rank_state = ?", "Trash").Error
//
//	return Results, DeleteTotal, err
//}
//
//
//// ChangeRankSectionsState 修改指定分类的状态
//func ChangeRankSectionsState(ids []uint, state string) ([]RankSection, error) {
//	// 批量查询 rankSection
//	var rankSections []RankSection
//	if err := rankRankSectionTable().Where("id IN (?)", ids).Find(&rankSections).Error; err != nil {
//		return nil, err
//	}
//
//	// 如果目标状态为 "Trash"，检查是否有子分类或者关联的rank项目
//	if state == "Trash" {
//		for _, section := range rankSections {
//			var count int64
//			rankRankSectionTable().Where("rank_sec_parent_id = ? and rank_sec_state != ?", section.ID, "Draft").Count(&count)
//			if count > 0 {
//				return nil, fmt.Errorf("分类 %s (%d) 存在子分类，请先删除子分类", section.RankSecName, section.ID)
//			}
//
//			var rankCount int64
//			rankRankTable().Where("rank_sec_id = ? and rank_state != ?", section.ID, "Deleted").Count(&rankCount)
//			if rankCount > 0 {
//				return nil, fmt.Errorf("分类 %s (%d) 存在关联的rank项目，请先删除关联的rank项目", section.RankSecName, section.ID)
//			}
//		}
//	}
//
//	// 批量更新状态
//	err := rankRankSectionTable().Where("id IN (?)", ids).Updates(map[string]interface{}{"rank_sec_state": state}).Error
//	if err != nil {
//		return nil, err
//	}
//
//	// 批量查询并返回结果
//	var updatedSections []RankSection
//	if err := rankRankSectionTable().Where("id IN (?)", ids).Find(&updatedSections).Error; err != nil {
//		return nil, err
//	}
//	return updatedSections, nil
//}
//
//func DeleteTrashRankSections() error {
//	return rankRankSectionTable().Where("rank_sec_state = ?", "Trash").Delete(&RankSection{}).Error
//}
//
//func CheckRankSectionNameExistsAndNotEqual(name string, id int) error {
//	var count int64
//	if id == 0 {
//		rankRankSectionTable().Where("rank_sec_name = ?", name).Count(&count)
//	} else {
//		rankRankSectionTable().Where("rank_sec_name = ? and id != ?", name, id).Count(&count)
//	}
//	if count > 0 {
//		return errors.New("分类名称已存在")
//	}
//	return nil
//}
//
//func CreateOrUpdateRankSection(myRankSection RankSection) (*RankSection, error) {
//	rankSection := RankSection{
//		ID:              myRankSection.ID,
//		RankSecParentID: myRankSection.RankSecParentID,
//		RankSecName:     myRankSection.RankSecName,
//		RankSecState:    myRankSection.RankSecState,
//		CreateDatetime:  time.Now(),
//		UpdateDatetime:  time.Now(),
//	}
//
//	if myRankSection.ID == 0 {
//		if err := rankRankSectionTable().Create(&rankSection).Error; err != nil {
//			return nil, err
//		}
//	} else {
//		if err := rankRankSectionTable().Save(&rankSection).Error; err != nil {
//			return nil, err
//		}
//	}
//
//	var resultRankSection RankSection
//	rankRankSectionTable().Where("id =?", rankSection.ID).Find(&resultRankSection)
//
//	return &resultRankSection, nil
//}
//
//func GetRankSections(start, limit int) ([]RankSection, error) {
//	var rankSecList []RankSection
//	if err := rankRankSectionTable().Offset(start).Limit(limit).Order("create_datetime DESC").Find(&rankSecList).Error; err != nil {
//		return nil, err
//	}
//	return rankSecList, nil
//}
//
//func GetPublishedRankSections(start, limit int) ([]RankSection, error) {
//	var rankSecList []RankSection
//	if err := rankRankSectionTable().Where("rank_sec_state = ?", "Published").Offset(start).Limit(limit).Order("create_datetime DESC").Find(&rankSecList).Error; err != nil {
//		return nil, err
//	}
//	return rankSecList, nil
//}
//
//func GetTrashedRankSections(start, limit int) ([]RankSection, error) {
//	var rankSecList []RankSection
//	if err := rankRankSectionTable().Where("rank_sec_state = ?", "Trash").Offset(start).Limit(limit).Order("create_datetime DESC").Find(&rankSecList).Error; err != nil {
//		return nil, err
//	}
//	return rankSecList, nil
//}
//
//func GetRankSectionsCount() (trashTotal, publishedTotal int64, err error) {
//
//	if err = rankRankSectionTable().Where("rank_sec_state = ?", "Trash").Count(&trashTotal).Error; err != nil {
//		return 0, 0, err
//	}
//	if err = rankRankSectionTable().Where("rank_sec_state = ?", "Published").Count(&publishedTotal).Error; err != nil {
//		return 0, 0, err
//	}
//	return trashTotal, publishedTotal, nil
//}
