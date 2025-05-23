// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dal

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q               = new(Query)
	ApiKey          *apiKey
	Device          *device
	DevicePlatform  *devicePlatform
	Platform        *platform
	Resource        *resource
	Site            *site
	User            *user
	UserInteraction *userInteraction
	ValueStream     *valueStream
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	ApiKey = &Q.ApiKey
	Device = &Q.Device
	DevicePlatform = &Q.DevicePlatform
	Platform = &Q.Platform
	Resource = &Q.Resource
	Site = &Q.Site
	User = &Q.User
	UserInteraction = &Q.UserInteraction
	ValueStream = &Q.ValueStream
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:              db,
		ApiKey:          newApiKey(db, opts...),
		Device:          newDevice(db, opts...),
		DevicePlatform:  newDevicePlatform(db, opts...),
		Platform:        newPlatform(db, opts...),
		Resource:        newResource(db, opts...),
		Site:            newSite(db, opts...),
		User:            newUser(db, opts...),
		UserInteraction: newUserInteraction(db, opts...),
		ValueStream:     newValueStream(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	ApiKey          apiKey
	Device          device
	DevicePlatform  devicePlatform
	Platform        platform
	Resource        resource
	Site            site
	User            user
	UserInteraction userInteraction
	ValueStream     valueStream
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:              db,
		ApiKey:          q.ApiKey.clone(db),
		Device:          q.Device.clone(db),
		DevicePlatform:  q.DevicePlatform.clone(db),
		Platform:        q.Platform.clone(db),
		Resource:        q.Resource.clone(db),
		Site:            q.Site.clone(db),
		User:            q.User.clone(db),
		UserInteraction: q.UserInteraction.clone(db),
		ValueStream:     q.ValueStream.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:              db,
		ApiKey:          q.ApiKey.replaceDB(db),
		Device:          q.Device.replaceDB(db),
		DevicePlatform:  q.DevicePlatform.replaceDB(db),
		Platform:        q.Platform.replaceDB(db),
		Resource:        q.Resource.replaceDB(db),
		Site:            q.Site.replaceDB(db),
		User:            q.User.replaceDB(db),
		UserInteraction: q.UserInteraction.replaceDB(db),
		ValueStream:     q.ValueStream.replaceDB(db),
	}
}

type queryCtx struct {
	ApiKey          IApiKeyDo
	Device          IDeviceDo
	DevicePlatform  IDevicePlatformDo
	Platform        IPlatformDo
	Resource        IResourceDo
	Site            ISiteDo
	User            IUserDo
	UserInteraction IUserInteractionDo
	ValueStream     IValueStreamDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		ApiKey:          q.ApiKey.WithContext(ctx),
		Device:          q.Device.WithContext(ctx),
		DevicePlatform:  q.DevicePlatform.WithContext(ctx),
		Platform:        q.Platform.WithContext(ctx),
		Resource:        q.Resource.WithContext(ctx),
		Site:            q.Site.WithContext(ctx),
		User:            q.User.WithContext(ctx),
		UserInteraction: q.UserInteraction.WithContext(ctx),
		ValueStream:     q.ValueStream.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
