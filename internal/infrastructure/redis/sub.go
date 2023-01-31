package redis


type Task struct {
	BackupId int    `json:"backup_id"`
	TaskId   string `json:"task_id"`
	EntryId  string `json:"entry_id"`
}

func StateSub() {

	//
}
