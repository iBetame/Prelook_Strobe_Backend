package API

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/sunbelife/Prelook_Strobe_Backend/Model"
)

type SearchRankBody struct {
	ID              int       `json:"id,omitempty"`
	Start           int       `json:"start"`
	Limit           int       `json:"limit" binding:"required"`
	State           string    `json:"state" binding:"required"`
	Keyword         string    `json:"keyword"`
	rankSecId       string    `json:"cat1"`
	rankSecParentId string    `json:"cat2"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
}

func UpdateRankItem(c *gin.Context) {
	var myRankItem Model.RankingItem

	if err := c.ShouldBind(&myRankItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
			"data":    err.Error(),
		})
		return
	}

	finalRank, err := Model.UpdateRankItem(myRankItem)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "更新失败，可能是参数错误",
			"data":    err.Error(),
		})
	} else {
		Message := "更新成功"
		if myRankItem.ID == 0 {
			Message = "创建成功"
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": Message,
			"data":    finalRank,
		})
	}
}

func UpdateRank(c *gin.Context) {
	var myRank Model.Rank

	if err := c.ShouldBind(&myRank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误",
			"data":    err.Error(),
		})
		return
	}

	finalRank, err := Model.UpdateRank(myRank)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "更新失败 " + err.Error(),
			"data":    err.Error(),
		})
	} else {
		Message := "更新成功"
		if myRank.ID == 0 {
			Message = "创建成功"
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": Message,
			"data":    finalRank,
		})
	}
}

func GetRankList(c *gin.Context) {
	start := c.Query("start")
	limit := c.Query("limit")
	status := c.Query("type")
	keyword := c.Query("keyword")
	order := c.Query("order")

	// 调用 Model 层获取类型为 Dock 的排行榜数据

	Result, DataTotal, AllTotal, PublishedTotal, DraftTotal, RemovedTotal, err := Model.SearchRankList(start, limit, status, keyword, order)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无法执行，可能参数有误或没有项目",
			"data":    err,
		})
	} else {
		// 输出为 JSON
		c.JSON(http.StatusOK, gin.H{
			"code":            http.StatusOK,
			"message":         "ok, i found them",
			"all_total":       AllTotal,
			"published_total": PublishedTotal,
			"draft_total":     DraftTotal,
			"removed_total":   RemovedTotal,
			"data_total":      DataTotal,
			"data":            Result,
		})
	}

	// 如果 Model 层返回了错误信息，将错误信息返回给客户端
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "获取排行榜数据失败",
			"data":    err.Error(),
		})
		return
	}
}

func GetRankItemList(c *gin.Context) {
	rank_id := c.Query("rank_id")
	start := cast.ToInt(c.Query("start"))
	limit := cast.ToInt(c.Query("limit"))
	status := c.Query("type")
	keyword := c.Query("keyword")
	order := c.Query("order")

	// 调用 Model 层获取类型为 Dock 的排行榜数据

	Result, DataTotal, AllTotal, PublishedTotal, DraftTotal, RemovedTotal, err := Model.SearchRankItemList(rank_id, start, limit, status, keyword, order)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "无法执行，可能参数有误或没有项目",
			"data":    err,
		})
	} else {
		// 排序
		for k := range Result {
			Result[k].RankingIndex = cast.ToString(start + k + 1)
		}
		// 输出为 JSON
		c.JSON(http.StatusOK, gin.H{
			"code":            http.StatusOK,
			"message":         "ok, i found them",
			"all_total":       AllTotal,
			"published_total": PublishedTotal,
			"draft_total":     DraftTotal,
			"removed_total":   RemovedTotal,
			"data_total":      DataTotal,
			"data":            Result,
		})
	}

	// 如果 Model 层返回了错误信息，将错误信息返回给客户端
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "获取排行榜数据失败",
			"data":    err.Error(),
		})
		return
	}
}

func GetRankByID(c *gin.Context) {
	id := c.Query("id")

	// 调用 Model 层获取类型为 Dock 的排行榜数据
	rankData, err := Model.GetRankByID(id)

	// 如果 Model 层返回了错误信息，将错误信息返回给客户端
	if err != nil {
		c.SecureJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "获取排行榜数据失败",
			"data":    err.Error(),
		})
		return
	}

	// 将获取到的类型为 Dock 的排行榜数据返回给客户端
	c.SecureJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "获取排行榜数据成功",
		"data":    rankData,
	})
}

func UpdateRankState(c *gin.Context) {
	var json struct {
		IDs   []uint `json:"ids"`
		State string `json:"state"`
	}
	c.ShouldBind(&json)
	_, err := Model.ChangeRankState(json.IDs, json.State)
	if err != nil {
		c.SecureJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "更新状态失败",
			"data":    err.Error(),
		})
		return
	}
	c.SecureJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "更新状态成功",
	})
}

func UpdateRankItemDeleteState(c *gin.Context) {
	var json struct {
		IDs   []uint `json:"ids"`
		State string `json:"state"`
	}
	c.ShouldBind(&json)
	_, err := Model.ChangeRankItemDeleteState(json.IDs, json.State)
	if err != nil {
		c.SecureJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "更新状态失败",
			"data":    err.Error(),
		})
		return
	}
	c.SecureJSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "更新状态成功",
	})
}

//
//// GetRankDock 根据 RankMode 获取类型为 Dock 的排行榜数据
//// @param c *gin.Context 请求上下文
//func GetRankDock(c *gin.Context) {
//	// 调用 Model 层获取类型为 Dock 的排行榜数据
//	rankList, err := Model.GetRankDock()
//
//	// 如果 Model 层返回了错误信息，将错误信息返回给客户端
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"code":    http.StatusInternalServerError,
//			"message": "获取 Dock 类型的排行榜数据失败",
//			"data":    err.Error(),
//		})
//		return
//	}
//
//	// 将获取到的类型为 Dock 的排行榜数据返回给客户端
//	c.JSON(http.StatusOK, gin.H{
//		"code":    http.StatusOK,
//		"message": "获取 Dock 类型的排行榜数据成功",
//		"data":    rankList,
//	})
//}
//
//func GetRankByID(c *gin.Context) {
//	rankID := c.Query("id")
//	rank, err := Model.GetRankByID(rankID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"code":    http.StatusInternalServerError,
//			"message": "获取 Rank 失败",
//			"data":    err.Error(),
//		})
//		return
//	}
//	if rank == nil {
//		c.JSON(http.StatusNotFound, gin.H{
//			"code":    http.StatusNotFound,
//			"message": "指定的 Rank 不存在",
//			"data":    nil,
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code":    http.StatusOK,
//		"message": "获取 Rank 成功",
//		"data":    rank,
//	})
//}
//
//func SearchRankList(c *gin.Context) {
//	// 检查请求方法是否为 POST
//	if c.Request.Method != "POST" {
//		c.JSON(http.StatusMethodNotAllowed, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "请求错误，请检查",
//			"data":    nil,
//		})
//		c.AbortWithStatus(http.StatusMethodNotAllowed)
//		return
//	}
//
//	// 解析请求体，检查是否为 JSON 格式
//	var searchRankBody SearchRankBody
//	err := c.BindJSON(&searchRankBody)
//
//	if err != nil {
//		c.JSON(http.StatusMethodNotAllowed, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "参数错误，请检查",
//			"data":    nil,
//		})
//		c.AbortWithStatus(http.StatusBadRequest)
//		return
//	}
//
//	// 对请求参数进行处理
//	ID := searchRankBody.ID
//	start := searchRankBody.Start
//	limit := searchRankBody.Limit
//	state := searchRankBody.State
//	keyword := searchRankBody.Keyword
//	rankSecId := searchRankBody.rankSecId
//	rankSecParentId := searchRankBody.rankSecParentId
//	startDate := searchRankBody.StartDate
//	endDate := searchRankBody.EndDate
//
//	Result, DataTotal, AllTotal, PublishedTotal, DraftTotal, ScheduledTotal, TrashTotal, err := Model.SearchRankList(start, limit, ID, rankSecId, rankSecParentId, state, keyword, startDate, endDate)
//
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "无法执行，可能参数有误或没有项目",
//			"data":    err,
//		})
//	} else {
//		// 输出为 JSON
//		c.JSON(http.StatusOK, gin.H{
//			"code":            http.StatusOK,
//			"message":         "ok, i found them",
//			"all_total":       AllTotal,
//			"published_total": PublishedTotal,
//			"draft_total":     DraftTotal,
//			"scheduled_total": ScheduledTotal,
//			"trash_total":     TrashTotal,
//			"data_total":      DataTotal,
//			"data":            Result,
//		})
//	}
//}
//
//// 批量修改分类状态
//func ChangeRankSectionsState(c *gin.Context) {
//	var requestData struct {
//		IDs   []uint `json:"ids" binding:"required"`
//		State string `json:"state" binding:"required"`
//	}
//
//	if err := c.ShouldBindJSON(&requestData); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "参数错误",
//			"data":    err.Error(),
//		})
//		return
//	}
//
//	if requestData.State != "Trash" && requestData.State != "Published" && requestData.State != "Draft" {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "无效的状态",
//			"data":    nil,
//		})
//		return
//	}
//
//	rankSections, err := Model.ChangeRankSectionsState(requestData.IDs, requestData.State)
//	if err != nil {
//		if strings.Contains(err.Error(), "请先删除子分类") {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"code":    http.StatusBadRequest,
//				"message": err.Error(),
//				"data":    nil,
//			})
//			return
//		}
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"code":    http.StatusInternalServerError,
//			"message": err.Error(),
//			"data":    nil,
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code":    http.StatusOK,
//		"message": "批量修改分类状态成功",
//		"data":    rankSections,
//	})
//}
//
//// 清空多个
//func EmptyRank(c *gin.Context) {
//	rankList, deleteTotal, err := Model.EmptyRank()
//
//	if err != nil || deleteTotal == 0 {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "无法执行删除，可能参数有误或没有可删除的项目",
//			"data":    err,
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code":          http.StatusOK,
//		"message":       "删除成功",
//		"deleted_list":  rankList,
//		"deleted_total": deleteTotal,
//	})
//}
//
//// 移至删除
////func SetBannerIsDelete(c *gin.Context)  {
////	var myData Model.Data
////
////	c.ShouldBind(&myData)
////
////	if len(myData.IDs) == 0 {
////		c.JSON(http.StatusBadRequest, gin.H{
////			"Code": http.StatusBadRequest,
////			"message":"set fail, your data is empty?",
////		})
////		return
////	}
////
////	finalBanners, err := Model.SetBannerIsDelete(myData.IDs, myData.Mode)
////
////	if err != nil {
////		c.JSON(http.StatusBadRequest, gin.H{
////			"Code": http.StatusBadRequest,
////			"message":"delete fail, maybe parms error?",
////			"data": err,
////		})
////	} else {
////		// 输出为 JSON
////		c.JSON(http.StatusOK, gin.H{
////			"Code": http.StatusOK,
////			"message": "ok, i have deleted them",
////			"data" : finalBanners,
////		})
////	}
////}
//
////func GetBannerByID(c *gin.Context)  {
////	StrID := c.Query("ID")
////
////	myBanner, err := Model.GetBannerByID(StrID)
////
////	if err != nil || myBanner.ID == 0 {
////		c.JSON(http.StatusBadRequest, gin.H{
////			"Code" : http.StatusBadRequest,
////			"message" : "get fail, maybe parms error?",
////			"data" : err,
////		})
////	} else {
////		// 输出为 JSON
////		c.JSON(http.StatusOK, gin.H{
////			"Code" : http.StatusOK,
////			"message" : "ok, i have found it.",
////			"data" : myBanner,
////		})
////	}
////}
//
//func GetRankSections(c *gin.Context) {
//	// 获取参数
//	state := c.Query("state")
//	start, err := strconv.Atoi(c.DefaultQuery("start", "0"))
//	if err != nil || start < 0 {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "start 参数错误",
//		})
//		return
//	}
//	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
//	if err != nil || limit < 1 {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "limit 参数错误",
//		})
//		return
//	}
//
//	var rankSecList []Model.RankSection
//
//	// 根据 state 获取分类列表
//	switch state {
//	case "All":
//		rankSecList, err = Model.GetRankSections(start, limit)
//	case "Published":
//		rankSecList, err = Model.GetPublishedRankSections(start, limit)
//	case "Trash":
//		rankSecList, err = Model.GetTrashedRankSections(start, limit)
//	default:
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "参数错误",
//		})
//		return
//	}
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"code":    http.StatusInternalServerError,
//			"message": "获取分类列表失败",
//		})
//		return
//	}
//
//	// 获取数据总数
//	trashTotal, publishedTotal, err := Model.GetRankSectionsCount()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"code":    http.StatusInternalServerError,
//			"message": "获取分类数目失败",
//		})
//		return
//	}
//
//	// 返回分类列表和数量
//	c.JSON(http.StatusOK, gin.H{
//		"code":            http.StatusOK,
//		"message":         "获取成功",
//		"data":            rankSecList,
//		"data_total":      len(rankSecList),
//		"trash_total":     trashTotal,
//		"published_total": publishedTotal,
//	})
//}
//
//func CreateOrUpdateRankSection(c *gin.Context) {
//	var myRankSection Model.RankSection
//
//	if err := c.ShouldBind(&myRankSection); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "参数错误",
//			"data":    err.Error(),
//		})
//		return
//	}
//
//	if myRankSection.RankSecName == "" {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "分类名称不能为空",
//			"data":    nil,
//		})
//		return
//	}
//
//	if err := Model.CheckRankSectionNameExistsAndNotEqual(myRankSection.RankSecName, myRankSection.ID); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "分类名称已存在",
//			"data":    err.Error(),
//		})
//		return
//	}
//
//	finalRankSection, err := Model.CreateOrUpdateRankSection(myRankSection)
//
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "创建或更新失败，可能是参数错误",
//			"data":    err.Error(),
//		})
//	} else {
//		message := "更新成功"
//		if myRankSection.ID == 0 {
//			message = "创建成功"
//		}
//		c.JSON(http.StatusOK, gin.H{
//			"code":    http.StatusOK,
//			"message": message,
//			"data":    finalRankSection,
//		})
//	}
//}
//
//// 批量修改分类状态
//func ChangeRankState(c *gin.Context) {
//	var requestData struct {
//		IDs   []uint `json:"ids" binding:"required"`
//		State string `json:"state" binding:"required"`
//	}
//
//	if err := c.ShouldBindJSON(&requestData); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "参数错误",
//			"data":    err.Error(),
//		})
//		return
//	}
//
//	if requestData.State != "Trash" && requestData.State != "Published" {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "无效的状态",
//			"data":    nil,
//		})
//		return
//	}
//
//	rankSections, err := Model.ChangeRankState(requestData.IDs, requestData.State)
//	if err != nil {
//		if strings.Contains(err.Error(), "请先删除子分类") {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"code":    http.StatusBadRequest,
//				"message": err.Error(),
//				"data":    nil,
//			})
//			return
//		}
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"code":    http.StatusInternalServerError,
//			"message": err.Error(),
//			"data":    nil,
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"code":    http.StatusOK,
//		"message": "批量修改分类状态成功",
//		"data":    rankSections,
//	})
//}
//
//func DeleteTrashRankSections(c *gin.Context) {
//	if err := Model.DeleteTrashRankSections(); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"code":    http.StatusInternalServerError,
//			"message": "删除归档分类失败",
//			"data":    err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code":    http.StatusOK,
//		"message": "删除归档分类成功",
//		"data":    nil,
//	})
//}
//
////func UpdateRankSection(c *gin.Context)  {
////	var myRankSection Model.RankSection
////
////	if err := c.ShouldBind(&myRankSection); err != nil {
////		c.JSON(http.StatusBadRequest, gin.H{
////			"Code":    http.StatusBadRequest,
////			"Message": "参数错误",
////			"Data":    err.Error(),
////		})
////		return
////	}
////
////	finalRankSection, err := Model.UpdateRankSection(myRankSection)
////
////	if err != nil {
////		c.JSON(http.StatusBadRequest, gin.H{
////			"Code":    http.StatusBadRequest,
////			"Message": "更新失败，可能是参数错误",
////			"Data":    err.Error(),
////		})
////	} else {
////		Message := "更新成功"
////		if myRankSection.RankSecID == "-1" {
////			Message = "创建成功"
////		}
////		c.JSON(http.StatusOK, gin.H{
////			"Code":    http.StatusOK,
////			"Message": Message,
////			"Data":    finalRankSection,
////		})
////	}
////}
////
//func UpdateRank(c *gin.Context) {
//	var myRank Model.Rank
//
//	if err := c.ShouldBind(&myRank); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "参数错误",
//			"data":    err.Error(),
//		})
//		return
//	}
//
//	finalRank, err := Model.UpdateRank(myRank)
//
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"code":    http.StatusBadRequest,
//			"message": "更新失败，可能是参数错误",
//			"data":    err.Error(),
//		})
//	} else {
//		Message := "更新成功"
//		if myRank.ID == nil {
//			Message = "创建成功"
//		}
//		c.JSON(http.StatusOK, gin.H{
//			"code":    http.StatusOK,
//			"message": Message,
//			"data":    finalRank,
//		})
//	}
//}
