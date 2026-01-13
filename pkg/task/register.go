package task

import (
	"fmt"
	"github.com/yuanJewel/go-core/task"
	"time"
)

var RegisteredTask = map[string]task.Func{
	"test success": {F: testSuccess, Cancel: true},
	"test error":   {F: testError, Cancel: false},
}

func testError(id string, data ...interface{}) (string, error) {
	err := task.SetVariable(id, "test-str", fmt.Sprintf("%v", data[0]))
	if err != nil {
		return "", err
	}
	s := fmt.Sprintf("task(%s) start in %s, input is %v, variable is %s", id, time.Now().String(), data,
		task.GetVariable(id, "test-str"))
	fmt.Println("yuanTag Error " + s)
	time.Sleep(1 * time.Second)
	return s, fmt.Errorf("%v", data)
}

func testSuccess(id string, data ...interface{}) (string, error) {
	s := fmt.Sprintf("task(%s) start in %s, input is %v, variable is %s", id, time.Now().String(),
		data, task.GetVariable(id, "test-list"))
	err := task.AppendVariable(id, "test-list", fmt.Sprintf("%v", data[0]))
	if err != nil {
		return "", err
	}
	err = task.SetRecycleKeyExpireTime(id, 100*time.Second)
	if err != nil {
		return "", err
	}

	for i := 0; i < 10; i++ {
		time.Sleep(10000 * time.Millisecond)
		fmt.Printf("yuanTag Success [%s]: %d\n", s, i)
	}
	return s, nil
}
