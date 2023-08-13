package core

import (
	"context"
	"errors"
	"fmt"
	"gin-example/config/config"
	"go.uber.org/zap"
)

type Queue struct {
	Jobs   map[string]Job
	task   []TaskInterFace
	Logger *zap.Logger
}

var QueueApp = &Queue{Jobs: map[string]Job{}}

// 绑定 自定义的 TaskInterFace
func (q *Queue) Bind(taskI TaskInterFace) {
	q.task = append(q.task, taskI)
}

// 将绑定的 自定义 TaskInterFace 绑定到 job，监听 消费
func (q *Queue) Register(j *Job, name string) {
	q.Jobs[name] = *j
	if j.Child.Enable() {
		go j.Run()
	}
}

// 实例化 job
func (q *Queue) NewJob(cfg config.Queue) error {
	for _, task := range q.task {
		j, err := NewJob(task, cfg, q.Logger)
		if err != nil {
			return err
		}
		q.Register(j, j.Child.GetName())
	}
	return nil
}

// 提供 push 接口
func (q *Queue) Push(queue string, ctx context.Context, message string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint("推送失败", r))
		}
	}()
	job := q.Jobs[queue]
	err = job.Push(ctx, message)
	return err
}
