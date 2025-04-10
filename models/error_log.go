package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// ErrorLog represents a logged error in the application
type ErrorLog struct {
	Id        int       `orm:"auto"`
	Message   string    `orm:"size(255)"`
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
	UserId    int       `orm:"index"` // Link to the user who encountered the error
}

func init() {
	orm.RegisterModel(new(ErrorLog))
}

func AddErrorLog(message string, userId int) error {
	o := orm.NewOrm()
	errorLog := &ErrorLog{
		Message:   message,
		UserId:    userId,
		Timestamp: time.Now(),
	}
	_, err := o.Insert(errorLog)
	return err
}

func GetErrorLogs() ([]*ErrorLog, error) {
	o := orm.NewOrm()
	var logs []*ErrorLog
	_, err := o.QueryTable("error_log").OrderBy("-timestamp").All(&logs)
	return logs, err
}
