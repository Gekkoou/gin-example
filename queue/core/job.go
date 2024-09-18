package core

import (
	"context"
	"fmt"
	"gin-example/config/config"
	"gin-example/queue/drive"
	"github.com/bytedance/sonic"
	"go.uber.org/ratelimit"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// 一个协程最多处理任务数，防止协程泄露
var _maxTaskNum = 1000

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
	// 并发限制
	ch := make(chan struct{}, t.Child.GetConsumerNumber())
	// 流速限制 ,每秒
	rlc := t.Child.GetRateLimit()
	var rl ratelimit.Limiter
	if rlc > 0 {
		rl = ratelimit.New(rlc) // per second
	}
	for {
		select {
		case ch <- struct{}{}:
			go t.RunHandel(ch, rl)
		}
	}
}

func (t *Job) RunHandel(ch chan struct{}, rl ratelimit.Limiter) {

	defer func() {
		if r := recover(); r != nil {
			t.ErrorLogger.Printf(fmt.Sprintln(t.Child.GetName(), "消费失败", r))
		}
		<-ch
	}()

	ctx := context.Background()
	retryCount := t.Child.GetRetryCount()
	gid := getGID()
	maxTask := _maxTaskNum
	for {

		if maxTask < 0 {
			return
		}
		maxTask = maxTask - 1

		// 限流
		if rl != nil {
			rl.Take()
		}

		m, err := t.Conn.GetMessage(ctx)
		fmt.Println("协程ID:", gid, ",message:", m, ",", t.Child.GetName(), "，时间：", time.Now().Format(time.DateTime))
		if err != nil {
			t.ErrorLogger.Printf(fmt.Sprintln(t.Child.GetName(), "拉取信息失败", err))
			time.Sleep(5 * time.Second)
			continue
		}

		/*// 测试并协程数
		rn := rand.Intn(10)
		if rn > 5 {
			panic("随机抛异常")
		}*/

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
}

func getGID() int64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	goidStr := strings.TrimPrefix(string(b), "goroutine ")
	goidStr = goidStr[:strings.Index(goidStr, " ")]
	gid, err := strconv.ParseInt(goidStr, 10, 64)
	if err != nil {
		return -1
	}
	return gid
}
