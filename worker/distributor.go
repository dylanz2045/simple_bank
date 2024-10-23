package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistrbuteTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

// 初始化一个redis的任务队列库
func NewRedisTaskDistributor(redisDpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisDpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
