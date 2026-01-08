package demo

import (
	"github.com/kataras/iris/v12"
	"github.com/yuanJewel/go-core/api"
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/task"
)

// @Summary 获取作业任务
// @Description 获取作业任务
// @Param id query string false "job"
// @Param id header string false "job"
// @Param page header string false "page"
// @Param body body task.JobInfo false "Info"
// @tags job
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/task/job [get]
func getJobs(ctx iris.Context) {
	service.GetDbInfoByIdsAndOrder(ctx, task.Job{}, &[]task.Job{}, "date desc")
}

// @Summary 创建作业任务
// @Description 创建作业任务
// @Param body body []task.StepInfo true "Info"
// @tags job
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/task/job [post]
func postJob(ctx iris.Context) {
	task.CreateTaskContext(ctx)
}

// @Summary 修改作业任务
// @Description 修改作业任务
// @Param id query string false "job"
// @Param id header string false "job"
// @Param body body task.JobInfo true "Info"
// @tags job
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/task/job [put]
func putJob(ctx iris.Context) {
	// State 修改为 ABORTED 会阻断任务继续运行
	service.PutDbInfoById(ctx, "get-jobs", &task.Job{}, api.NormalSpecialTask)
}

// @Summary 删除作业任务
// @Description 删除作业任务
// @Param id query string false "job"
// @Param id header string false "job"
// @tags job
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/task/job [delete]
func deleteJobs(ctx iris.Context) {
	service.DeleteDb(ctx, "get-jobs", &task.Job{})
}

// @Summary 获取作业任务
// @Description 获取作业任务
// @Param id query string false "step"
// @Param id header string false "step"
// @Param job_id query string true "job_id"
// @Param job_id header string true "job_id"
// @Param page header string false "page"
// @Param body body task.StepInfo false "Info"
// @tags job
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/task/step [get]
func getSteps(ctx iris.Context) {
	service.GetDbInfo(ctx, task.Step{}, &[]task.Step{}, []string{"id", "job_id"},
		[]string{"date desc", "start_time desc"})
}

// @Summary 获取作业任务
// @Description 获取作业任务
// @Param id query string false "step"
// @Param id header string false "step"
// @Param page header string false "page"
// @Param body body task.StepInfo false "Info"
// @tags job
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/task/register [get]
func getRegister(ctx iris.Context) {
	response := api.ResponseInit(ctx)
	api.ResponseBody(ctx, response, task.GetRegisteredTaskNames())
}

// @Summary 获取作业状态
// @Description 获取作业状态
// @tags job
// @Accept json
// @Produce json
// @Success 200 {object} api.Response "ok"
// @Failure 401 string string "未授权"
// @Failure 403 {object} api.Response "权限不足"
// @Failure 501 {object} api.Response "处理存在异常"
// @Security ApiKeyAuth
// @Router /api/v1/task/status [get]
func getStepStatus(ctx iris.Context) {
	task.GetRunningTaskStatus(ctx)
}
