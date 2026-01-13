package task

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/logger"
	"time"
)

func createTask(job string, stage int) error {
	var (
		redisKeyList = []interface{}{}
		steps        []Step
		signatures   []*tasks.Signature
	)
	logger.Log.Debugf("start create task in job(%s[%d])", job, stage)
	number, err := service.Instance.GetAllItems(Step{JobId: job, StepInfo: StepInfo{Stage: stage}}, &steps)
	if err != nil {
		return err
	}
	if number == 0 {
		return finishTask(job, stage)
	}

	for _, step := range steps {
		if step.State == tasks.StateSuccess {
			continue
		}
		stepId := signatureId(step.ID)
		redisKeyList = append(redisKeyList, stepId)
		signatures = append(signatures, &tasks.Signature{
			UUID: stepId,
			Name: step.Tag,
			Args: []tasks.Arg{{
				Type:  "string",
				Value: step.ID,
			}, {
				Type:  "string",
				Value: step.Option,
			}},
			IgnoreWhenTaskNotRegistered: true,
			OnSuccess:                   newSignature(step.ID, "success"),
			OnError:                     newSignature(step.ID, "error"),
			RetryCount:                  0,
		})
	}

	_, err = service.Instance.UpdateItem(Step{JobId: job, StepInfo: StepInfo{Stage: stage}}, &Step{
		State:     tasks.StateStarted,
		StartTime: time.Now(),
	}, int64(len(steps)))
	if err != nil {
		return err
	}
	group, err := newGroup(signatures...)
	if err != nil {
		return err
	}
	redisKeyList = append(redisKeyList, group.GroupUUID)
	err = redisInstance.SAdd(recycleKeyId(job), finishExpiration, redisKeyList...)
	if err != nil {
		return fmt.Errorf("add recycle redis key error: %v", err)
	}

	finishId := fmt.Sprintf("finish:%s", group.GroupUUID)
	chord, err := tasks.NewChord(group, &tasks.Signature{
		UUID: finishId,
		Name: "finish",
		Args: []tasks.Arg{{
			Type:  "string",
			Value: job,
		}, {
			Type:  "int",
			Value: stage,
		}},
		IgnoreWhenTaskNotRegistered: true,
		RetryCount:                  0,
		OnSuccess:                   []*tasks.Signature{},
		OnError: []*tasks.Signature{{
			UUID: finishId + ":error",
			Name: "error",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: finishId + "-job-" + job,
				},
			},
			RetryCount: 0,
		}},
	})
	if err != nil {
		return err
	}
	_, err = machineryInstance.SendChord(chord, maxConcurrency())
	if err == nil {
		logger.Log.Debugf("create task(%v) successfully", steps)
	}
	return err
}

func newGroup(signatures ...*tasks.Signature) (*tasks.Group, error) {
	groupUUID := uuid.New().String()
	groupID := fmt.Sprintf("group:%v", groupUUID)

	for _, signature := range signatures {
		if signature.UUID == "" {
			signatureID := uuid.New().String()
			signature.UUID = fmt.Sprintf("step:%v", signatureID)
		}
		signature.GroupUUID = groupID
		signature.GroupTaskCount = len(signatures)
	}

	return &tasks.Group{
		GroupUUID: groupID,
		Tasks:     signatures,
	}, nil
}

func finishTask(jobId string, stage int, _ ...interface{}) error {
	var (
		job   Job
		steps []Step
	)

	logger.Log.Debugf("start finish task in job(%s[%d])", jobId, stage)
	exists, err := service.Instance.GetItem(Job{ID: jobId}, &job)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("job(%s) is not exists", jobId)
	}
	if job.State == StateAborted {
		return finishAbort(jobId)
	}

	if stage >= job.TotalStage {
		return finishSuccess(job)
	}

	_, err = service.Instance.GetAllItems(Step{JobId: jobId, StepInfo: StepInfo{Stage: stage}}, &steps)
	if err != nil {
		return err
	}

	for _, s := range steps {
		if s.State != tasks.StateSuccess {
			if s.State != tasks.StateStarted {
				logger.Log.Debugf("task(%s) is not success in job(%s[%d]), %v", s.ID, jobId, stage, s)
				return finishError(jobId)
			}
			task, err := machineryInstance.GetBackend().GetState(signatureId(s.ID))
			if err != nil {
				logger.Log.Debugf("get task(%s) state error: %v", s.ID, err)
				return err
			}
			if task.State != tasks.StateSuccess {
				logger.Log.Debugf("task(%s) is not success in job(%s[%d]), %v", s.ID, jobId, stage, s)
				return finishError(jobId)
			}
		}
	}

	_, err = service.Instance.UpdateItem(Job{ID: jobId}, &Job{
		JobInfo: JobInfo{
			ActiveStage: job.ActiveStage + 1,
		},
	}, 1)
	if err != nil {
		return err
	}
	logger.Log.Debugf("finish task in job(%s[%d]) successfully", jobId, stage)
	return createTask(jobId, stage+1)
}

func newSignature(id string, status string) []*tasks.Signature {
	if status != "success" && status != "error" {
		return nil
	}
	return []*tasks.Signature{{
		UUID: status + ":" + id,
		Name: status,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: id,
			},
		},
		RetryCount: 0,
	}}
}

func signatureId(id string) string {
	return fmt.Sprintf("step:%v", id)
}
