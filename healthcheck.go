package main

import (
	"errors"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/admin"
)

type DatabaseCheck struct{}

func (dc *DatabaseCheck) Check() error {
	o := orm.NewOrm()
	_, err := o.Raw("SELECT 1").Exec()
	if err != nil {
		return errors.New("can't connect to database: " + err.Error())
	}
	return nil
}

func initHealthChecks() {
	admin.AddHealthCheck("database", &DatabaseCheck{})
}
