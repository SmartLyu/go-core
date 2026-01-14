package task

import "time"

type Job struct {
	ID         string    `gorm:"column:id;type:varchar(255);primary_key" json:"id"`
	CreateDate time.Time `gorm:"column:date;type:datetime" json:"date"`
	JobInfo
}

type JobInfo struct {
	User        string    `gorm:"column:user;type:varchar(255)" json:"user"`
	ActiveStage int       `gorm:"column:active_stage;type:int" json:"active_stage"`
	TotalStage  int       `gorm:"column:total_stage;type:int" json:"total_stage"`
	Option      string    `gorm:"column:option;type:text" json:"option"`
	State       string    `gorm:"column:state;type:varchar(255)" json:"state"`
	FinishTime  time.Time `gorm:"column:finish_time;type:datetime" json:"finish_time"`
}

type Step struct {
	ID         string    `gorm:"column:id;type:varchar(255);primary_key" json:"id"`
	JobId      string    `gorm:"column:job_id;type:varchar(255)" json:"job_id"`
	Job        Job       `gorm:"foreignKey:JobId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreateDate time.Time `gorm:"column:date;type:datetime" json:"date"`
	FinishTime time.Time `gorm:"column:finish_time;type:datetime" json:"finish_time"`
	Result     string    `gorm:"column:result;type:text" json:"result"`
	Error      string    `gorm:"column:error;type:text" json:"error"`
	State      string    `gorm:"column:state;type:varchar(255)" json:"state"`
	StepInfo
}

type StepInfo struct {
	Name      string    `gorm:"column:name;type:varchar(255)" json:"name"`
	Tag       string    `gorm:"column:tag;type:varchar(255)" json:"tag"`
	Stage     int       `gorm:"column:stage;type:int" json:"stage"`
	Option    string    `gorm:"column:option;type:text" json:"option"`
	StartTime time.Time `gorm:"column:start_time;type:datetime" json:"start_time"`
}
