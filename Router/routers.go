package Router

import (
	"github.com/gin-gonic/gin"
	"github.com/sunbelife/Prelook_Strobe_Backend/API"
	"github.com/sunbelife/Prelook_Strobe_Backend/MiddleWare"
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	// AdminUser
	r.POST("/AdminLoginIn", API.AdminLoginIn)
	// WxUser
	//r.POST("/WxLoginIn", API.WxLoginIn)

	apiv1manage := r.Group("v1/manage")
	apiv1manage.Use(MiddleWare.JWT())
	{
		// Rank
		apiv1manage.GET("/getRankList", API.GetRankList)
		apiv1manage.GET("/getRankItemList", API.GetRankItemList)
		apiv1manage.GET("/getRankByID", API.GetRankByID)
		apiv1manage.POST("/updateRank", API.UpdateRank)
		apiv1manage.POST("/updateRankItem", API.UpdateRankItem)
		apiv1manage.POST("/updateRankState", API.UpdateRankState)
		apiv1manage.POST("/updateRankItemState", API.UpdateRankItemDeleteState)
	}

	//apiv1wx := r.Group("/wx")
	//apiv1wx.Use()
	//{
	//	// Index 首页 = 排行榜图片 * 4 + 头图 * 4 + 大牌直降 * 1 + 正在团购 * 1 + 信息流 * 10
	//	//apiv1wx.GET("/GetPriceDownDock", API.GetPriceDownDock)
	//	//apiv1wx.GET("/GetGroupBuyDock", API.GetGroupBuyDock)
	//	//apiv1wx.GET("/GetMsgCard", API.GetMsgCard)
	//
	//	// User
	//
	//	// Rankings
	//	// 创建/修改 Rank 信息
	//	apiv1wx.POST("/Rank/Update", API.UpdateRank)
	//	// 搜索 Rank
	//	apiv1wx.POST("/Rank/Search", API.SearchRankList)
	//	// 获取 Rank Dock
	//	apiv1wx.GET("/Rank/GetDock", API.GetRankDock)
	//	// 查询单独 Rank
	//	apiv1wx.GET("/Rank/GetByID", API.GetRankByID)
	//	// 变更 Rank 为 Draft
	//	apiv1wx.POST("/Rank/Trash", API.ChangeRankState)
	//	// 清空 Rank
	//	apiv1wx.GET("/Rank/Empty", API.EmptyRank)
	//	// 获取全部分类
	//	apiv1wx.GET("/RankSec/GetAll", API.GetRankSections)
	//	// 创建/修改 Rank 分类
	//	apiv1wx.POST("/RankSec/Update", API.CreateOrUpdateRankSection)
	//	// 变更 Rank 分类为 Draft
	//	apiv1wx.POST("/RankSec/Trash", API.ChangeRankSectionsState)
	//	// 清空 Draft 状态 Rank 分类
	//	apiv1wx.GET("/RankSec/Empty", API.DeleteTrashRankSections)
	//	// Product
	//	//apiv1wx.GET("/IsTrackNumRepeat", API.IsTrackNumRepeat)
	//	//apiv1wx.POST("/UpdateProduct", API.UpdateProduct)
	//	//apiv1wx.GET("/GetProductByID", API.GetProductByID)
	//	////apiv1wx.GET("/EmptyProduct", API.EmptyProduct)
	//	////apiv1wx.POST("/SetProductIsDelete", API.SetProductIsDelete)
	//	////apiv1wx.GET("/SetProductIsTop/:ID/:Mode", API.SetProductIsTop)
	//	////apiv1wx.POST("/SetProductIsShow", API.SetProductIsShow)
	//	////apiv1wx.GET("/GetProductList", API.GetProductList)
	//	//apiv1wx.GET("/GetProductList", API.SearchProductList)
	//	//
	//	//// News
	//	////apiv1wx.POST("/UpdateNews", API.UpdateNews)
	//	//apiv1wx.GET("/GetNewsByUnionID", API.GetNewsByUnionID)
	//	////apiv1wx.GET("/EmptyNews", API.EmptyNews)
	//	////apiv1wx.POST("/SetNewsIsDelete", API.SetNewsIsDelete)
	//	////apiv1wx.GET("/SetNewsIsTop/:ID/:Mode", API.SetNewsIsTop)
	//	////apiv1wx.POST("/SetNewsIsShow", API.SetNewsIsShow)
	//	////apiv1wx.GET("/GetNewsList", API.GetNewsList)
	//	//apiv1wx.GET("/GetNewsList", API.SearchNewsList)
	//
	//	// Banner
	//	apiv1wx.POST("/Banner/Update", API.UpdateBanner)
	//	//apiv1wx.GET("/GetBannerByID", API.GetBannerByID)
	//	////apiv1wx.GET("/EmptyBanner", API.EmptyBanner)
	//	////apiv1wx.POST("/SetBannerIsDelete", API.SetBannerIsDelete)
	//	////apiv1wx.GET("/SetBannerIsTop/:ID/:Mode", API.SetBannerIsTop)
	//	////apiv1wx.POST("/SetBannerIsShow", API.SetBannerIsShow)
	//	////apiv1wx.GET("/GetBannerList", API.GetBannerList)
	//	//apiv1wx.GET("/GetBannerList", API.SearchBannerList)
	//	//
	//	//// Settings
	//	//apiv1wx.GET("/GetAllSettings", API.GetAllSettings)
	//	//apiv1wx.GET("/GetSettings", API.GetSettings)
	//	////apiv1wx.POST("/SetAllSettings", API.SetAllSettings)
	//	//
	//	//// User
	//	//apiv1wx.POST("/GetUserSettings", API.GetUserSettings)
	//	//apiv1wx.POST("/SetUserSettings", API.SetUserSettings)
	//	//apiv1wx.POST("/CheckSubscribeState", API.CheckSubscribeState)
	//	//apiv1wx.POST("/DeleteSubscribeState", API.DeleteSubscribeState)
	//	//apiv1wx.POST("/AddSubscribeTimes", API.AddSubscribeTimes)
	//	//apiv1wx.POST("/ResetSubscribeTimes", API.ResetSubscribeTimes)
	//	//
	//	//// Topic
	//	//apiv1wx.GET("/GetTopicList", API.SearchTopicList)
	//	//apiv1wx.GET("/GetTopicByUnionID", API.GetTopicByUnionID)
	//	//
	//	//// Repair
	//	//apiv1wx.GET("/GetRepairList", API.SearchRepairList)
	//	//apiv1wx.GET("/GetRepairListByOpenID", API.GetRepairListByOpenID)
	//	//apiv1wx.GET("/GetRepairListByOpenIDAndProductID", API.GetRepairListByOpenIDAndProductID)
	//	//apiv1wx.GET("/GetCaseByRepairID", API.GetCaseByRepairID)
	//	//apiv1wx.GET("/GetCaseByTrackNum", API.GetCaseByTrackNum)
	//	//apiv1wx.GET("/SetRepairStatusByRepairID", API.SetRepairStatusByRepairID)
	//	//apiv1wx.POST("/UpdateRepair", API.UpdateRepair)
	//}

	return r
}
