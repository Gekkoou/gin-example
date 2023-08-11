package core

import (
	"context"
	"errors"
	"fmt"
	"gin-example/config/config"
)

type Queue struct {
	Jobs map[string]Job
	task []TaskInterFace
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
func (q *Queue) NewJob(conf config.Queue) error {
	for _, task := range q.task {
		j := NewJob(task, task.GetConnType(), conf)
		q.Register(j, j.Child.GetName())
	}
	return nil
}

// 提供 push 接口
func (q *Queue) Push(face TaskInterFace, ctx context.Context, message string) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint("推送失败", r))
		}
	}()
	job := q.Jobs[face.GetName()]
	err = job.Push(ctx, message)
	return err
}
