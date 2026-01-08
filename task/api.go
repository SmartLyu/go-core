package task

import (
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/kataras/iris/v12"
	"github.com/yuanJewel/go-core/api"
	"github.com/yuanJewel/go-core/db/service"
	"strings"
	"time"
)

func CreateTaskContext(ctx iris.Context) {
	response := api.ResponseInit(ctx)
	dbInstance := service.Instance.WithContext(ctx)
	body, err := ctx.GetBody()
	if err != nil {
		api.ReturnErr(api.GetBodyError, ctx, err, response)
		return
	}
	var job Job

	err = CreateTask(&job, dbInstance, body, api.GetUserName(ctx))
	if err != nil {
		api.ReturnErr(api.CreateTaskError, ctx, err, response)
		return
	}
	api.ResponseBody(ctx, response, job)
}

type StatusObject struct {
	ID        string    `json:"id"`
	Task      string    `json:"task"`
	Job       string    `json:"job"`
	Tag       string    `json:"tag"`
	Stage     int       `json:"stage"`
	StartTime time.Time `json:"start_time"`
	Output    string    `json:"output"`
}

func GetRunningTaskStatus(ctx iris.Context) {
	response := api.ResponseInit(ctx)
	dbInstance := service.Instance.WithContext(ctx)

	var (
		steps   = []Step{}
		results = []StatusObject{}
	)
	_, err := dbInstance.GetAllItems(Step{State: tasks.StateStarted}, &steps)
	if err != nil {
		api.ReturnErr(api.GetTaskStatusError, ctx, err, response)
		return
	}

	for _, step := range steps {
		_status, err := GetStatus(step.ID)
		if err != nil {
			api.ReturnErr(api.GetTaskStatusError, ctx, err, response)
			return
		}
		results = append(results, StatusObject{ID: step.ID, Task: step.Name, Job: step.JobId, Tag: step.Tag,
			Stage: step.Stage, StartTime: step.StartTime, Output: strings.Join(_status, "\n")})
	}
	api.ResponseBody(ctx, response, results)
}
