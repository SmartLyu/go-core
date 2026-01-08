package task

import (
	"fmt"
	"reflect"
	"time"
)

func registerTasks(namedTaskFuncs map[string]Func) error {
	taskMap := make(map[string]interface{})
	for key, f := range namedTaskFuncs {
		if key == "success" || key == "error" || key == "finish" {
			taskMap[key] = f.F
			continue
		}
		taskMap[key] = wrapWithLogic(f)
	}
	return machineryInstance.RegisterTasks(taskMap)
}

func wrapWithLogic(f Func) interface{} {
	fnValue := reflect.ValueOf(f.F)
	return reflect.MakeFunc(fnValue.Type(), func(args []reflect.Value) (results []reflect.Value) {
		makeErrorResult := func(err error) []reflect.Value {
			outType0 := fnValue.Type().Out(0) // string 类型
			results = make([]reflect.Value, 2)
			results[0] = reflect.New(outType0).Elem()
			results[1] = reflect.ValueOf(err)
			return results
		}

		id := args[0].String()
		if args[0].Kind() != reflect.String {
			return makeErrorResult(fmt.Errorf("参数类型错误"))
		}
		err := SetStatus(id, "开始执行任务")
		if err != nil {
			return makeErrorResult(fmt.Errorf("设置任务开始状态失败: %v", err))
		}
		err = LockTaskState(id, "task")
		if err != nil {
			return makeErrorResult(fmt.Errorf("获取任务执行锁失败: %v", err))
		}

		if f.Cancel {
			cancelTask[id] = make(chan int)
			go func() {
				results = fnValue.Call(args)
				cancelTask[id] <- 1
			}()
			end := <-cancelTask[id]
			if end == 1 {
				delete(cancelTask, id)
			} else {
				_ = SetStatus(id, "任务中止")
				return makeErrorResult(fmt.Errorf(AbortedRetryError))
			}
		} else {
			results = fnValue.Call(args)
		}

		err = SetStatus(id, "任务执行完成")
		if err != nil {
			results[1] = reflect.ValueOf(fmt.Errorf("设置任务完结状态失败: %v", err))
		}
		return results
	}).Interface()
}

// SetStatus 设置当前执行信息
func SetStatus(id, value string) error {
	_, err := getJobId(id)
	if err != nil {
		return err
	}
	key := statusKey(id)
	defer needRecycleKey(id, key)
	return redisInstance.RPush(key,
		fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), value), lockExpiration)
}

// GetStatus 获取当前执行信息
func GetStatus(id string) ([]string, error) {
	_, err := getJobId(id)
	if err != nil {
		return nil, err
	}
	key := statusKey(id)
	ok, err := redisInstance.Exists(key)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return redisInstance.LRange(key, 0, -1)
}

func statusKey(id string) string {
	return fmt.Sprintf("step:%s:status", id)
}
