package api

import (
	"buding-job/orm"
	"buding-job/orm/do"
	"buding-job/web/dto"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"strconv"
)

// JobInfoApi /*todo 除了添加的方法其他方法都差个更新缓存以及替换cron操作*/
type JobInfoApi struct {
	router *gin.RouterGroup
}

func NewJobInfoApi(router *gin.Engine) *JobInfoApi {
	group := router.Group("/jobInfo")
	return &JobInfoApi{
		router: group,
	}
}

func (jobInfoApi *JobInfoApi) Router() {
	jobInfoApi.router.POST("/add", jobInfoApi.Add)
	jobInfoApi.router.POST("/update", jobInfoApi.update)
	jobInfoApi.router.GET("/stop", jobInfoApi.stop)
	jobInfoApi.router.GET("/start", jobInfoApi.start)
	jobInfoApi.router.GET("/delete", jobInfoApi.delete)
}

func (jobInfoApi *JobInfoApi) update(ctx *gin.Context) {
	var jobDto dto.JobInfoDto
	ctx.BindJSON(&jobDto)
	if jobDto.Id == 0 {
		ctx.JSON(500, gin.Error{
			Meta: "修改任务需要绑定id",
		})
		ctx.Done()
		return
	}
	var jobInfoDo do.JobInfoDo
	copier.Copy(&jobInfoDo, &jobDto)
	orm.DB.Updates(&jobInfoDo)

	ctx.JSON(200, map[string]string{
		"msg": "ok",
	})
}

func (jobInfoApi *JobInfoApi) Add(ctx *gin.Context) {
	var jobDto dto.JobInfoDto
	ctx.BindJSON(&jobDto)
	//管理器
	var jobManagement do.JobManagementDo
	orm.DB.First(&jobManagement, jobDto.ManageId)
	if jobManagement.Id == 0 {
		ctx.JSON(500, gin.Error{
			Meta: "任务管理器不存在",
		})
		ctx.Done()
		return
	}
	//添加
	var jobInfoDo do.JobInfoDo
	copier.Copy(&jobInfoDo, &jobDto)
	jobInfoDo.Enable = false
	orm.DB.Create(&jobInfoDo)
	//添加缓存
	ctx.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (jobInfoApi *JobInfoApi) stop(context *gin.Context) {
	i := context.Query("id")
	//获取
	id, _ := strconv.ParseInt(i, 10, 64)
	var jobInfoDo do.JobInfoDo
	orm.DB.First(&jobInfoDo, id)
	if jobInfoDo.Id == 0 {
		context.JSON(500, gin.Error{
			Meta: "条目不存在",
		})
		context.Done()
		return
	}
	if !jobInfoDo.Enable {
		context.JSON(200, map[string]interface{}{
			"message": "ok",
		})
		context.Done()
		return
	}
	orm.DB.Model(&do.JobInfoDo{}).Where("id = ?", id).Update("is_enable", 0)
	//添加任务
	context.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (jobInfoApi *JobInfoApi) start(context *gin.Context) {
	i := context.Query("id")
	//获取
	id, _ := strconv.ParseInt(i, 10, 64)
	var jobInfoDo do.JobInfoDo
	orm.DB.First(&jobInfoDo, id)
	if jobInfoDo.Id == 0 {
		context.JSON(500, gin.Error{
			Meta: "条目不存在",
		})
		context.Done()
		return
	}
	if jobInfoDo.Enable {
		context.JSON(200, map[string]interface{}{
			"message": "ok",
		})
		context.Done()
		return
	}
	orm.DB.Model(&do.JobInfoDo{}).Where("id = ?", id).Update("is_enable", 1)
	//添加任务
	context.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}

func (jobInfoApi *JobInfoApi) delete(context *gin.Context) {
	i := context.Query("id")
	//获取
	id, _ := strconv.ParseInt(i, 10, 64)
	var jobInfoDo do.JobInfoDo
	orm.DB.First(&jobInfoDo, id)
	if jobInfoDo.Id == 0 {
		context.JSON(500, gin.Error{
			Meta: "条目不存在",
		})
		context.Done()
		return
	}
	orm.DB.Delete(&jobInfoDo)
	//删除缓存
	context.JSON(200, map[string]interface{}{
		"message": "ok",
	})
}
