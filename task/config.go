package task

import (
	"os"
	"strconv"

	"github.com/yuanJewel/go-core/config"
)

type Task struct {
	Tag         string       `required:"true" yaml:"tag" json:"tag" env:"task.tag"`
	Concurrency int          `default:"10" yaml:"concurrency" json:"concurrency" env:"task.concurrency"`
	IsWorker    bool         `default:"false" yaml:"worker" json:"worker" env:"task.worker"`
	Redis       config.Redis `yaml:"redis" json:"redis"`
	RabbitMq    RabbitMq     `yaml:"mq" json:"mq"`
	// RunExpiration is Task running expiration time.
	// After the task start, the parameters and lock retention time will be recycled according to this time.
	// If the task fails and is blocked, the recycling time will increase exponentially.
	// The storage time is evaluated by the redis service pressure.
	// It is generally set to 5 minutes and is recommended to be no less than 60 seconds
	// and is recommended to be no more than 30 minute.
	RunExpiration int `default:"300" yaml:"results_expiration" json:"results_expiration" env:"task.results_expiration"`
	// LockExpiration is Task atomic protection lock expiration time.
	// Set according to the estimated maximum time for the task. You can set a longer time to enhance protection.
	LockExpiration int `default:"18000" yaml:"lock_expiration" json:"lock_expiration" env:"task.lock_expiration"`
	// VarExpiration is Task parameter expiration time.
	VarExpiration int `default:"3600" yaml:"var_expiration" json:"var_expiration" env:"task.var_expiration"`
}

type RabbitMq struct {
	Host     string `required:"true" yaml:"host" json:"host" env:"task.mq.host"`
	Port     string `default:"5672" yaml:"port" json:"port" env:"task.mq.port"`
	Username string `required:"true" yaml:"username" json:"username" env:"task.mq.username"`
	Password string `required:"true" yaml:"password" json:"password" env:"task.mq.password"`
	Queue    string `default:"machinery_tasks" yaml:"queue" json:"queue" env:"task.mq.queue"`
	Exchange string `default:"machinery_exchange" yaml:"exchange" json:"exchange" env:"task.mq.exchange"`
}

func maxConcurrency() int {
	concurrency := os.Getenv("TASK_MAX_CONCURRENCY")
	if concurrency != "" {
		concurrencyNumber, err := strconv.Atoi(concurrency)
		if err == nil {
			return concurrencyNumber
		}
	}
	return 5
}
