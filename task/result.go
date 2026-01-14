package task

import (
	"encoding/json"
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/logger"
	"regexp"
	"strings"
	"time"
)

var DesensitizationKeyList = []string{"password", "passwd", "token", "auth"}

func GetRegisteredTaskNames() []string {
	return machineryInstance.GetRegisteredTaskNames()
}

func LockError(id string) error {
	return fmt.Errorf("任务%s已被锁定无法执行", id)
}

// LockTaskState Non-idempotent tasks require additional protection locks to use this module
func LockTaskState(id string, suffix ...string) error {
	lockId := lockTaskKey(id, suffix...)
	lock, err := redisInstance.Exists(lockId)
	if err != nil {
		logger.Log.Errorf("redis for task is error, please check %v", err)
		return nil
	}
	if lock {
		logger.Log.Debugf("task(%s) has been locked: %s", id, lockId)
		return LockError(id)
	}
	err = redisInstance.Set(lockId, time.Now().String(), lockExpiration)
	if err != nil {
		return err
	}
	defer needRecycleKey(id, lockId)
	return nil
}

func lockTaskKey(id string, suffix ...string) string {
	return fmt.Sprintf("step:%s:%s:lock", id, strings.Join(suffix, ":"))
}

func resultToDb(id string, _ ...interface{}) error {
	task, err := machineryInstance.GetBackend().GetState(signatureId(id))
	if err != nil {
		return err
	}
	for _, r := range task.Results {
		if r.Type == "string" {
			v, ok := r.Value.(string)
			if !ok {
				break
			}
			if v == LockError(id).Error() {
				logger.Log.Warnf("step %s has been locked", id)
				return nil
			}
		}
	}

	t := LockTaskState(id, "end")
	if t != nil {
		return nil
	}
	step := Step{}

	exists, err := service.Instance.GetItem(Step{ID: id}, &step)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("cannot find step: %s", id)
	}
	if step.State == tasks.StateSuccess {
		logger.Log.Errorf("step %s has run over once", id)
		return nil
	}

	_, err = service.Instance.UpdateItem(Step{ID: id}, &Step{
		State:      task.State,
		Result:     HumanReadableResults(task.Results),
		FinishTime: time.Now(),
	}, 1)
	if err != nil {
		return err
	}
	logger.Log.Debugf("result success step %s successfully", id)
	return nil
}

func errorToDb(errorStr, id string, _ ...interface{}) error {
	t := LockTaskState(id, "end")
	if t != nil {
		return nil
	}
	step := Step{}
	signatureKey := signatureId(id)

	if strings.HasPrefix(id, "finish:") && strings.Contains(id, "-job-") {
		idList := strings.Split(id, "-job-")
		signatureKey = idList[0]
		id = idList[0]
		finishObject.Error(idList[1])
		_, err := service.Instance.AddItem(&Step{
			ID:         id,
			JobId:      idList[1],
			CreateDate: time.Now(),
			StepInfo: StepInfo{
				Name:      "任务完结模块抛出异常",
				Tag:       "finish",
				Stage:     0,
				StartTime: time.Now(),
				Option:    "",
			},
		}, 1)
		if err != nil {
			return err
		}
	}
	exists, err := service.Instance.GetItem(Step{ID: id}, &step)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("cannot find step: %s", id)
	}

	// Special logic: When a task is interrupted and a retry is triggered.
	if errorStr == AbortedRetryError {
		jobId, err := getJobId(id)
		if err != nil {
			return err
		}
		err = redisInstance.Del(lockTaskKey(id, "end"))
		if err != nil {
			logger.Log.Errorf("删除异常完结锁失败: %v", err)
		}
		return createTask(jobId, step.Stage)
	}

	if step.State == tasks.StateFailure {
		logger.Log.Errorf("step %s has run over once", id)
		return nil
	}

	task, err := machineryInstance.GetBackend().GetState(signatureKey)
	if err != nil {
		return err
	}

	_, err = service.Instance.UpdateItem(Step{ID: id}, &Step{
		State:      task.State,
		Result:     HumanReadableResults(task.Results),
		Error:      errorStr,
		FinishTime: time.Now(),
	}, 1)
	if err != nil {
		return err
	}

	logger.Log.Debugf("result error step %s successfully", id)
	return finishError(step.JobId)
}

func HumanReadableResults(taskResults []*tasks.TaskResult) string {
	resultValues := make([]string, 0, 0)
	flagError := false
	for _, taskResult := range taskResults {
		if taskResult.Type == "string" {
			v, ok := taskResult.Value.(string)
			if !ok {
				flagError = true
				break
			}
			resultValues = append(resultValues, desensitizationResults(v))
		} else if taskResult.Type == "[]string" {
			v, ok := taskResult.Value.([]interface{})
			if !ok {
				flagError = true
				break
			}
			for _, _v := range v {
				__v, ok := _v.(string)
				if !ok {
					flagError = true
					break
				}
				resultValues = append(resultValues, desensitizationResults(__v))
			}
		} else {
			flagError = true
			break
		}
	}
	if flagError {
		v, err := tasks.ReflectTaskResults(taskResults)
		if err != nil {
			return "reflect task results error"
		}
		return tasks.HumanReadableResults(v)
	}
	output, err := json.Marshal(resultValues)
	if err != nil {
		return strings.Join(resultValues, "\n")
	}
	return string(output)
}

func desensitizationResults(input string) string {
	for _, key := range DesensitizationKeyList {
		reg := regexp.MustCompile(fmt.Sprintf("\"%s\":\"[^\"]*\"", key))
		input = reg.ReplaceAllString(input, fmt.Sprintf("\"%s\":\"***\"", key))
	}
	return input
}
