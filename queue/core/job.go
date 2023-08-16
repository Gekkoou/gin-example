package core

import (
	"context"
	"fmt"
	"gin-example/config/config"
	"gin-example/queue/drive"
	"github.com/bytedance/sonic"
	"time"
)

type JobErr struct {
	Name    string
	Err     string
	Message string
}

type Job struct {
	Conn        drive.Interface
	Name        string
	Child       TaskInterFace
	Logger      Logger
	ErrorLogger Logger
}

func NewJob(child TaskInterFace, cfg config.Queue, logger Logger, errorLogger Logger) (*Job, error) {
	j := &Job{Child: child, Logger: logger, ErrorLogger: errorLogger}
	err := j.SetType(child.GetConnType(), child.GetName(), cfg)
	return j, err
}

func (j *Job) SetType(t ConnType, name string, cfg config.Queue) (err error) {
	switch t {
	case Kafka:
		j.Conn, err = drive.NewKafka(name, cfg.Kafka, cfg.Prefix, cfg.FailureSuffix)
	case Redis:
		j.Conn, err = drive.NewRedis(name, cfg.Redis, cfg.Prefix, cfg.FailureSuffix)
	}
	return err
}

func (t *Job) Push(ctx context.Context, message string) error {
	err := t.Conn.Push(ctx, message)
	return err
}

func (t *Job) Run() {
	for i := 0; i < t.Child.GetConsumerNumber(); i++ {
		t.RunHandel()
	}
}

func (t *Job) RunHandel() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.ErrorLogger.Printf(fmt.Sprintln(t.Child.GetName(), "消费失败", r))
			}
		}()
		ctx := context.Background()
		retryCount := t.Child.GetRetryCount()
		for {
			m, err := t.Conn.GetMessage(ctx)
			if err != nil {
				t.ErrorLogger.Printf(fmt.Sprintln(t.Child.GetName(), "拉取信息失败", err))
				time.Sleep(5 * time.Second)
				continue
			}
			for i := 1; i <= retryCount; i++ {
				if err = t.Child.Handel(m); err == nil {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			if err != nil {
				t.ErrorLogger.Printf(fmt.Sprintf("%s 消费信息失败, msg: %+v, err: %+v", t.Child.GetName(), m, err))

				jobErrString, _ := sonic.MarshalString(&JobErr{
					Name:    t.Child.GetName(),
					Err:     err.Error(),
					Message: m,
				})
				if pfErr := t.Conn.PushFailure(ctx, jobErrString); pfErr != nil {
					t.ErrorLogger.Printf(fmt.Sprintf("%s PushFailure 失败, msg: %+v, err: %+v", t.Child.GetName(), m, pfErr))
				}
			}
			t.Conn.CommitMessage(ctx)
		}
	}()
}
