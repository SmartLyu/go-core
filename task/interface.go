package task

import (
	"fmt"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/yuanJewel/go-core/db/service"
	"github.com/yuanJewel/go-core/logger"
	"time"
)

type FinishInterface interface {
	Success(string)
	Error(string)
	Abort(string)
}

var finishObject FinishInterface

type FinishStruct struct{}

func (f *FinishStruct) Success(id string) {
	logger.Log.Infof("Task %s finish success !", id)
}

func (f *FinishStruct) Error(id string) {
	logger.Log.Errorf("Task %s finish error !", id)
}

func (f *FinishStruct) Abort(id string) {
	logger.Log.Warnf("Task %s finish abort !", id)
}

func finishAbort(id string) error {
	_, err := service.Instance.UpdateItem(Step{JobId: id, State: tasks.StatePending},
		&Step{
			State: StateAborted,
			Error: fmt.Sprintf("job %s has failed or aborted, terminate this step", id),
		}, -1)
	if err != nil {
		return err
	}
	err = recycleRedisKey(id, 2)
	if err != nil {
		return err
	}

	finishObject.Abort(id)
	return nil
}

func finishError(id string) error {
	if err := finishAbort(id); err != nil {
		return err
	}

	_, err := service.Instance.UpdateItem(Job{ID: id}, &Job{
		JobInfo: JobInfo{
			State:      tasks.StateFailure,
			FinishTime: time.Now(),
		},
	}, 0, 1)
	if err != nil {
		return err
	}
	err = recycleRedisKey(id, 4)
	if err != nil {
		return err
	}

	finishObject.Error(id)
	return nil
}

func finishSuccess(job Job) error {
	_, err := service.Instance.UpdateItem(Job{ID: job.ID}, &Job{
		JobInfo: JobInfo{
			State:       tasks.StateSuccess,
			ActiveStage: job.TotalStage,
			FinishTime:  time.Now(),
		},
	}, 1)
	if err != nil {
		return err
	}
	err = recycleRedisKey(job.ID, 1)
	if err != nil {
		return err
	}

	finishObject.Success(job.ID)
	return nil
}

// SetRecycleKeyExpireTime 设置步骤中任务需要使用的相关key的过期时间
func SetRecycleKeyExpireTime(id string, expire time.Duration) error {
	jobId, err := getJobId(id)
	if err != nil {
		return fmt.Errorf("get job id error: %s", err)
	}
	hasRecycleKey := make(map[string]bool)
	all, err := redisInstance.SMembers(recycleKeyId(jobId))

	fmt.Println("yuan test-1", recycleKeyId(jobId), all)
	expire += finishExpiration
	if err != nil {
		return err
	}

	for _, key := range all {
		if _, ok := hasRecycleKey[key]; ok {
			continue
		}
		hasRecycleKey[key] = true
		exists, err := redisInstance.Exists(key)
		if err != nil {
			return err
		}
		if !exists {
			err = redisInstance.SRem(recycleKeyId(jobId), key)
			if err != nil {
				return err
			}
			continue
		}
		err = redisInstance.Expire(key, expire)
		if err != nil {
			return err
		}
		logger.Log.Debugf("set key(%s) expire to: %v", key, expire)
	}
	return redisInstance.Expire(recycleKeyId(id), expire)
}

func needRecycleKey(stepId, key string) {
	id, err := getJobId(stepId)
	if err != nil {
		logger.Log.Errorln(err)
		return
	}
	err = redisInstance.SAdd(recycleKeyId(id), lockExpiration, key)
	if err != nil {
		return
	}
}

func recycleRedisKey(id string, speed int) error {
	hasRecycleKey := make(map[string]bool)
	all, err := redisInstance.SMembers(recycleKeyId(id))
	if err != nil {
		return err
	}

	for _, key := range all {
		if _, ok := hasRecycleKey[key]; ok {
			continue
		}
		hasRecycleKey[key] = true
		exists, err := redisInstance.Exists(key)
		if err != nil {
			return err
		}
		if !exists {
			continue
		}
		err = redisInstance.Expire(key, finishExpiration*time.Duration(speed))
		if err != nil {
			return err
		}
	}
	return redisInstance.Del(recycleKeyId(id))
}

func recycleKeyId(id string) string {
	return "recycle:key:" + id
}
