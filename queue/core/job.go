package core

import (
	"context"
	"fmt"
	"gin-example/config/config"
	"gin-example/queue/drive"
	"log"
)

type Job struct {
	Conn           drive.Interface
	Name           string
	Child          TaskInterFace
	ConsumerNumber int
}

func NewJob(child TaskInterFace, t ConnType, conf config.Queue) *Job {
	j := &Job{Child: child}
	j.SetType(t, child.GetName(), conf)
	j.ConsumerNumber = child.GetConsumerNumber()
	return j
}

func (j *Job) SetType(t ConnType, name string, conf config.Queue) {
	switch t {
	case Kafka:
		j.Conn = drive.NewKafka(name, conf.Kafka, conf.Prefix)
	case Redis:
		j.Conn = drive.NewRedis(name, conf.Redis, conf.Prefix)
	}
}

func (t *Job) Push(ctx context.Context, message string) error {
	err := t.Conn.Push(ctx, message)
	return err
}

func (t *Job) Run() {
	for i := 0; i < t.ConsumerNumber; i++ {
		t.RunHandel()
	}
}
func (t *Job) RunHandel() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("消费失败", r)
			}
		}()
		ctx := context.Background()
		for {
			m, err := t.Conn.GetMessage(ctx)
			if err != nil {
				log.Fatalln(err)
				// break
			}
			if err := t.Child.Handel(m); err != nil {
				log.Fatalln(err)
				return
			}
			t.Conn.CommitMessage(ctx)
		}
	}()
}
