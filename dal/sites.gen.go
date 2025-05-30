// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dal

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"app/model"
)

func newSite(db *gorm.DB, opts ...gen.DOOption) site {
	_site := site{}

	_site.siteDo.UseDB(db, opts...)
	_site.siteDo.UseModel(&model.Site{})

	tableName := _site.siteDo.TableName()
	_site.ALL = field.NewAsterisk(tableName)
	_site.ID = field.NewUint(tableName, "id")
	_site.CreatedAt = field.NewTime(tableName, "created_at")
	_site.UpdatedAt = field.NewTime(tableName, "updated_at")
	_site.DeletedAt = field.NewField(tableName, "deleted_at")
	_site.Name = field.NewString(tableName, "name")
	_site.Description = field.NewString(tableName, "description")
	_site.Address = field.NewString(tableName, "address")
	_site.City = field.NewString(tableName, "city")
	_site.State = field.NewString(tableName, "state")
	_site.Country = field.NewString(tableName, "country")
	_site.Metadata = field.NewString(tableName, "metadata")
	_site.Devices = siteHasManyDevices{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("Devices", "model.Device"),
		Site: struct {
			field.RelationField
			Devices struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("Devices.Site", "model.Site"),
			Devices: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Devices.Site.Devices", "model.Device"),
			},
		},
		ValueStream: struct {
			field.RelationField
			Devices struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("Devices.ValueStream", "model.ValueStream"),
			Devices: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Devices.ValueStream.Devices", "model.Device"),
			},
		},
		Platforms: struct {
			field.RelationField
			Resources struct {
				field.RelationField
			}
			Devices struct {
				field.RelationField
			}
		}{
			RelationField: field.NewRelation("Devices.Platforms", "model.Platform"),
			Resources: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Devices.Platforms.Resources", "model.Resource"),
			},
			Devices: struct {
				field.RelationField
			}{
				RelationField: field.NewRelation("Devices.Platforms.Devices", "model.Device"),
			},
		},
	}

	_site.fillFieldMap()

	return _site
}

type site struct {
	siteDo

	ALL         field.Asterisk
	ID          field.Uint
	CreatedAt   field.Time
	UpdatedAt   field.Time
	DeletedAt   field.Field
	Name        field.String
	Description field.String
	Address     field.String
	City        field.String
	State       field.String
	Country     field.String
	Metadata    field.String
	Devices     siteHasManyDevices

	fieldMap map[string]field.Expr
}

func (s site) Table(newTableName string) *site {
	s.siteDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s site) As(alias string) *site {
	s.siteDo.DO = *(s.siteDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *site) updateTableName(table string) *site {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewUint(table, "id")
	s.CreatedAt = field.NewTime(table, "created_at")
	s.UpdatedAt = field.NewTime(table, "updated_at")
	s.DeletedAt = field.NewField(table, "deleted_at")
	s.Name = field.NewString(table, "name")
	s.Description = field.NewString(table, "description")
	s.Address = field.NewString(table, "address")
	s.City = field.NewString(table, "city")
	s.State = field.NewString(table, "state")
	s.Country = field.NewString(table, "country")
	s.Metadata = field.NewString(table, "metadata")

	s.fillFieldMap()

	return s
}

func (s *site) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *site) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 12)
	s.fieldMap["id"] = s.ID
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
	s.fieldMap["deleted_at"] = s.DeletedAt
	s.fieldMap["name"] = s.Name
	s.fieldMap["description"] = s.Description
	s.fieldMap["address"] = s.Address
	s.fieldMap["city"] = s.City
	s.fieldMap["state"] = s.State
	s.fieldMap["country"] = s.Country
	s.fieldMap["metadata"] = s.Metadata

}

func (s site) clone(db *gorm.DB) site {
	s.siteDo.ReplaceConnPool(db.Statement.ConnPool)
	s.Devices.db = db.Session(&gorm.Session{Initialized: true})
	s.Devices.db.Statement.ConnPool = db.Statement.ConnPool
	return s
}

func (s site) replaceDB(db *gorm.DB) site {
	s.siteDo.ReplaceDB(db)
	s.Devices.db = db.Session(&gorm.Session{})
	return s
}

type siteHasManyDevices struct {
	db *gorm.DB

	field.RelationField

	Site struct {
		field.RelationField
		Devices struct {
			field.RelationField
		}
	}
	ValueStream struct {
		field.RelationField
		Devices struct {
			field.RelationField
		}
	}
	Platforms struct {
		field.RelationField
		Resources struct {
			field.RelationField
		}
		Devices struct {
			field.RelationField
		}
	}
}

func (a siteHasManyDevices) Where(conds ...field.Expr) *siteHasManyDevices {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a siteHasManyDevices) WithContext(ctx context.Context) *siteHasManyDevices {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a siteHasManyDevices) Session(session *gorm.Session) *siteHasManyDevices {
	a.db = a.db.Session(session)
	return &a
}

func (a siteHasManyDevices) Model(m *model.Site) *siteHasManyDevicesTx {
	return &siteHasManyDevicesTx{a.db.Model(m).Association(a.Name())}
}

func (a siteHasManyDevices) Unscoped() *siteHasManyDevices {
	a.db = a.db.Unscoped()
	return &a
}

type siteHasManyDevicesTx struct{ tx *gorm.Association }

func (a siteHasManyDevicesTx) Find() (result []*model.Device, err error) {
	return result, a.tx.Find(&result)
}

func (a siteHasManyDevicesTx) Append(values ...*model.Device) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a siteHasManyDevicesTx) Replace(values ...*model.Device) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a siteHasManyDevicesTx) Delete(values ...*model.Device) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a siteHasManyDevicesTx) Clear() error {
	return a.tx.Clear()
}

func (a siteHasManyDevicesTx) Count() int64 {
	return a.tx.Count()
}

func (a siteHasManyDevicesTx) Unscoped() *siteHasManyDevicesTx {
	a.tx = a.tx.Unscoped()
	return &a
}

type siteDo struct{ gen.DO }

type ISiteDo interface {
	gen.SubQuery
	Debug() ISiteDo
	WithContext(ctx context.Context) ISiteDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ISiteDo
	WriteDB() ISiteDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ISiteDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ISiteDo
	Not(conds ...gen.Condition) ISiteDo
	Or(conds ...gen.Condition) ISiteDo
	Select(conds ...field.Expr) ISiteDo
	Where(conds ...gen.Condition) ISiteDo
	Order(conds ...field.Expr) ISiteDo
	Distinct(cols ...field.Expr) ISiteDo
	Omit(cols ...field.Expr) ISiteDo
	Join(table schema.Tabler, on ...field.Expr) ISiteDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ISiteDo
	RightJoin(table schema.Tabler, on ...field.Expr) ISiteDo
	Group(cols ...field.Expr) ISiteDo
	Having(conds ...gen.Condition) ISiteDo
	Limit(limit int) ISiteDo
	Offset(offset int) ISiteDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ISiteDo
	Unscoped() ISiteDo
	Create(values ...*model.Site) error
	CreateInBatches(values []*model.Site, batchSize int) error
	Save(values ...*model.Site) error
	First() (*model.Site, error)
	Take() (*model.Site, error)
	Last() (*model.Site, error)
	Find() ([]*model.Site, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Site, err error)
	FindInBatches(result *[]*model.Site, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Site) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ISiteDo
	Assign(attrs ...field.AssignExpr) ISiteDo
	Joins(fields ...field.RelationField) ISiteDo
	Preload(fields ...field.RelationField) ISiteDo
	FirstOrInit() (*model.Site, error)
	FirstOrCreate() (*model.Site, error)
	FindByPage(offset int, limit int) (result []*model.Site, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Rows() (*sql.Rows, error)
	Row() *sql.Row
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ISiteDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (s siteDo) Debug() ISiteDo {
	return s.withDO(s.DO.Debug())
}

func (s siteDo) WithContext(ctx context.Context) ISiteDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s siteDo) ReadDB() ISiteDo {
	return s.Clauses(dbresolver.Read)
}

func (s siteDo) WriteDB() ISiteDo {
	return s.Clauses(dbresolver.Write)
}

func (s siteDo) Session(config *gorm.Session) ISiteDo {
	return s.withDO(s.DO.Session(config))
}

func (s siteDo) Clauses(conds ...clause.Expression) ISiteDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s siteDo) Returning(value interface{}, columns ...string) ISiteDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s siteDo) Not(conds ...gen.Condition) ISiteDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s siteDo) Or(conds ...gen.Condition) ISiteDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s siteDo) Select(conds ...field.Expr) ISiteDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s siteDo) Where(conds ...gen.Condition) ISiteDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s siteDo) Order(conds ...field.Expr) ISiteDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s siteDo) Distinct(cols ...field.Expr) ISiteDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s siteDo) Omit(cols ...field.Expr) ISiteDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s siteDo) Join(table schema.Tabler, on ...field.Expr) ISiteDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s siteDo) LeftJoin(table schema.Tabler, on ...field.Expr) ISiteDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s siteDo) RightJoin(table schema.Tabler, on ...field.Expr) ISiteDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s siteDo) Group(cols ...field.Expr) ISiteDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s siteDo) Having(conds ...gen.Condition) ISiteDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s siteDo) Limit(limit int) ISiteDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s siteDo) Offset(offset int) ISiteDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s siteDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ISiteDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s siteDo) Unscoped() ISiteDo {
	return s.withDO(s.DO.Unscoped())
}

func (s siteDo) Create(values ...*model.Site) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s siteDo) CreateInBatches(values []*model.Site, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s siteDo) Save(values ...*model.Site) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s siteDo) First() (*model.Site, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Site), nil
	}
}

func (s siteDo) Take() (*model.Site, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Site), nil
	}
}

func (s siteDo) Last() (*model.Site, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Site), nil
	}
}

func (s siteDo) Find() ([]*model.Site, error) {
	result, err := s.DO.Find()
	return result.([]*model.Site), err
}

func (s siteDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Site, err error) {
	buf := make([]*model.Site, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s siteDo) FindInBatches(result *[]*model.Site, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s siteDo) Attrs(attrs ...field.AssignExpr) ISiteDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s siteDo) Assign(attrs ...field.AssignExpr) ISiteDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s siteDo) Joins(fields ...field.RelationField) ISiteDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s siteDo) Preload(fields ...field.RelationField) ISiteDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s siteDo) FirstOrInit() (*model.Site, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Site), nil
	}
}

func (s siteDo) FirstOrCreate() (*model.Site, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Site), nil
	}
}

func (s siteDo) FindByPage(offset int, limit int) (result []*model.Site, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s siteDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s siteDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s siteDo) Delete(models ...*model.Site) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *siteDo) withDO(do gen.Dao) *siteDo {
	s.DO = *do.(*gen.DO)
	return s
}
