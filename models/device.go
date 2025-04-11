package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type Device struct {
	Id           int    `orm:"auto"`
	Name         string `orm:"size(100)"`
	Manufacturer string `orm:"size(100)"`
	Type         string `orm:"size(50)"`
	Building     string `orm:"size(10)"`
}

func init() {
	orm.RegisterModel(new(Device))
}

func GetAllDevices() ([]*Device, error) {
	o := orm.NewOrm()
	var devices []*Device
	_, err := o.QueryTable("device").All(&devices)
	return devices, err
}

// GetDeviceCount returns the total number of devices
func GetDeviceCount() (int64, error) {
	o := orm.NewOrm()
	count, err := o.QueryTable("device").Count()
	return count, err
}

func GetDeviceById(id int) (*Device, error) {
	o := orm.NewOrm()
	device := &Device{Id: id}
	err := o.Read(device)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return device, err
}

func AddDevice(device *Device) error {
	o := orm.NewOrm()
	_, err := o.Insert(device)
	return err
}

func UpdateDevice(device *Device) error {
	o := orm.NewOrm()
	_, err := o.Update(device)
	return err
}

func DeleteDevice(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Device{Id: id})
	return err
}
