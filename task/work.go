package task

import (
	"errors"
	"fmt"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/sirupsen/logrus"
	"github.com/yuanJewel/go-core/db/redis"
	"github.com/yuanJewel/go-core/logger"
	"net"
	"sync"
	"time"
)

const (
	AbortedRetryError = "任务被取消"
	StateAborted      = "ABORTED"
)

var (
	machineryInstance *machinery.Server
	redisInstance     *redis.Store
	lockExpiration    time.Duration
	varExpiration     time.Duration
	finishExpiration  time.Duration
	worker            *machinery.Worker
	workerChannel     = make(chan error)
	cancelTask        = make(map[string]chan int)
	stepToJob         sync.Map
)

type Func struct {
	F      interface{}
	Cancel bool
}

func InitWork(task Task, taskMap map[string]Func, f FinishInterface) (err error) {
	lockExpiration = time.Duration(task.LockExpiration) * time.Second
	varExpiration = time.Duration(task.VarExpiration) * time.Second
	finishExpiration = time.Duration(task.RunExpiration) * time.Second
	stepToJob = sync.Map{}
	machineryInstance, err = machinery.NewServer(&config.Config{
		Broker:          fmt.Sprintf("amqp://%s:%s@%s:%s", task.RabbitMq.Username, task.RabbitMq.Password, task.RabbitMq.Host, task.RabbitMq.Port),
		DefaultQueue:    task.RabbitMq.Queue,
		ResultBackend:   fmt.Sprintf("redis://%s@%s:%s/%d", task.Redis.Password, task.Redis.Host, task.Redis.Port, task.Redis.Db),
		ResultsExpireIn: task.RunExpiration,
		Redis: &config.RedisConfig{
			MaxIdle:      task.Redis.PoolSize,
			ReadTimeout:  task.Redis.Timeout,
			WriteTimeout: task.Redis.Timeout,
		},
		AMQP: &config.AMQPConfig{
			Exchange:      task.RabbitMq.Exchange,
			ExchangeType:  "direct",
			BindingKey:    task.RabbitMq.Queue,
			PrefetchCount: task.Concurrency,
		},
	})
	if err != nil {
		return
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(task.RabbitMq.Host, task.RabbitMq.Port), 3*time.Second)
	if err != nil || conn == nil {
		return fmt.Errorf("cannot connect task.RabbitMq(%s:%s), error: %v", task.RabbitMq.Host, task.RabbitMq.Port, err)
	}
	_ = conn.Close()
	conn, err = net.DialTimeout("tcp", net.JoinHostPort(task.Redis.Host, task.Redis.Port), 3*time.Second)
	if err != nil || conn == nil {
		return fmt.Errorf("cannot connect redis(%s:%s), error: %v", task.Redis.Host, task.Redis.Port, err)
	}
	_ = conn.Close()

	redisInstance, err = redis.GetRedisInstance(&task.Redis)
	if err != nil {
		return
	}

	taskMap["success"] = Func{resultToDb, false}
	taskMap["error"] = Func{errorToDb, false}
	taskMap["finish"] = Func{finishTask, false}

	err = registerTasks(taskMap)
	if err != nil {
		return
	}
	log.SetInfo(&wrapper{logrus.InfoLevel})
	log.SetDebug(&wrapper{logrus.DebugLevel})
	log.SetError(&wrapper{logrus.ErrorLevel})
	log.SetWarning(&wrapper{logrus.WarnLevel})

	logger.Log.Infof("Complete task system registration !")
	if task.IsWorker {
		worker = machineryInstance.NewWorker(task.Tag, task.Concurrency)
		logger.Log.Infof("Start one worker !")
		worker.LaunchAsync(workerChannel)
	}
	finishObject = f
	return
}

func BeforeExit() {
	logger.Log.Infof("正在优雅退出...")
	for id, channel := range cancelTask {
		channel <- 2
		err := redisInstance.Del(lockTaskKey(id, "task"))
		if err != nil {
			logger.Log.Errorf("删除任务锁失败: %v", err)
		}
	}
	err := <-workerChannel
	if errors.Is(err, machinery.ErrWorkerQuitGracefully) {
		logger.Log.Infof("任务正常退出")
	} else {
		logger.Log.Errorf("任务异常退出: %v", err)
	}
}
