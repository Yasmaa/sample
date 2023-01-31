package schedular

import (
	"api/config"
	"fmt"

	"github.com/hibiken/asynq"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Scheduler *asynq.Scheduler


func NewScheduler() *asynq.Scheduler {
	if Scheduler != nil {
		return Scheduler
	}
	Scheduler = asynq.NewScheduler(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", config.C.Redis.HOST, config.C.Redis.PORT), DB: 0}, nil)
	return Scheduler
}
