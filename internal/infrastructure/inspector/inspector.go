package inspector

import (
	"api/config"
	"fmt"

	"github.com/hibiken/asynq"
)

var Inspector *asynq.Inspector


func NewInspector() *asynq.Inspector {
	if Inspector != nil {
		return Inspector
	}
	Inspector = asynq.NewInspector(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", config.C.Redis.HOST, config.C.Redis.PORT), DB: 0})
	return Inspector
}
