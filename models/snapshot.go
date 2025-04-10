package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Snapshot represents a single snapshot of node values at a specific time
type Snapshot struct {
	Id            int       `orm:"auto"`
	DeviceId      int       `orm:"index"`
	IntegrationId int       `orm:"index"`
	Timestamp     time.Time `orm:"auto_now_add;type(datetime)"`
	Nodes         string    `orm:"type(text)"` // JSON string of node ID-value pairs
}

func init() {
	orm.RegisterModel(new(Snapshot))
}

// Snapshot CRUD operations
func AddSnapshot(snapshot *Snapshot) error {
	o := orm.NewOrm()
	_, err := o.Insert(snapshot)
	return err
}

func GetSnapshotsByDeviceId(deviceId int) ([]*Snapshot, error) {
	o := orm.NewOrm()
	var snapshots []*Snapshot
	_, err := o.QueryTable("snapshot").Filter("device_id", deviceId).OrderBy("-timestamp").All(&snapshots)
	return snapshots, err
}
