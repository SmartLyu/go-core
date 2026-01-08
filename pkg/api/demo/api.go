package demo

import (
	"github.com/kataras/iris/v12"
	"github.com/yuanJewel/go-core/api"
)

const prefixPath = ""

func Party(app iris.Party, topPath string) {
	app.Get(api.Fusion(topPath, prefixPath, "project"), getProjects).Name = "get-projects"
	app.Put(api.Fusion(topPath, prefixPath, "project"), putProject).Name = "put-project"
	app.Post(api.Fusion(topPath, prefixPath, "project"), postProjects).Name = "post-projects"
	app.Delete(api.Fusion(topPath, prefixPath, "project"), deleteProjects).Name = "delete-projects"

	app.Get(api.Fusion(topPath, "task", "job"), getJobs).Name = "get-task-jobs"
	app.Post(api.Fusion(topPath, "task", "job"), postJob).Name = "post-task-job"
	app.Put(api.Fusion(topPath, "task", "job"), putJob).Name = "put-task-jobs"
	app.Delete(api.Fusion(topPath, "task", "job"), deleteJobs).Name = "delete-task-jobs"
	app.Get(api.Fusion(topPath, "task", "register"), getRegister).Name = "get-task-register"
	app.Get(api.Fusion(topPath, "task", "step"), getSteps).Name = "get-task-steps"
	app.Get(api.Fusion(topPath, "task", "status"), getStepStatus).Name = "get-task-step-status"
}
