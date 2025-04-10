package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type Integration struct {
	Id              int     `orm:"auto"`
	Device          *Device `orm:"rel(fk)"`
	IntegrationType string  `orm:"size(50)"`
	Identifier      string  `orm:"size(100)"`
	Host            string  `orm:"size(255)"` // e.g., "example.com" or "192.168.1.10"
	Port            string  `orm:"size(10)"`  // e.g., "4840" for OPC UA, "8080" for REST
	Protocol        string  `orm:"size(50)"`  // e.g., "opc.tcp", "http", "https"
}

func init() {
	orm.RegisterModel(new(Integration))
}

func GetIntegrationsByDeviceId(deviceId int) ([]*Integration, error) {
	o := orm.NewOrm()
	var integrations []*Integration
	_, err := o.QueryTable("integration").Filter("device_id", deviceId).RelatedSel().All(&integrations)
	return integrations, err
}

func AddIntegration(integration *Integration) error {
	o := orm.NewOrm()
	_, err := o.Insert(integration)
	return err
}

func GetIntegrationById(id int) (*Integration, error) {
	o := orm.NewOrm()
	integration := &Integration{Id: id}
	err := o.Read(integration)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return integration, err
}

func GetAllIntegrations() ([]*Integration, error) {
	o := orm.NewOrm()
	var integrations []*Integration
	_, err := o.QueryTable("integration").All(&integrations)
	return integrations, err
}

func UpdateIntegration(integration *Integration) error {
	o := orm.NewOrm()
	_, err := o.Update(integration)
	return err
}

func DeleteIntegration(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Integration{Id: id})
	return err
}
