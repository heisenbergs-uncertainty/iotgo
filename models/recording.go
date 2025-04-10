package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Recording represents a time-series recording of a node's data
type Recording struct {
	Id            int       `orm:"auto"`
	DeviceId      int       `orm:"index"`
	IntegrationId int       `orm:"index"`
	NodeId        string    `orm:"size(255)"` // OPC UA Node ID
	NodeName      string    `orm:"size(100)"` // Node browse name
	Timestamp     time.Time `orm:"auto_now_add;type(datetime)"`
	Value         string    `orm:"size(255)"` // Node value as string (can be converted to appropriate type)
}

func init() {
	orm.RegisterModel(new(Recording))
}

// Recording CRUD operations
func AddRecording(recording *Recording) error {
	o := orm.NewOrm()
	_, err := o.Insert(recording)
	return err
}

func GetRecordingsByDeviceId(deviceId int) ([]*Recording, error) {
	o := orm.NewOrm()
	var recordings []*Recording
	_, err := o.QueryTable("recording").Filter("device_id", deviceId).OrderBy("-timestamp").All(&recordings)
	return recordings, err
}
