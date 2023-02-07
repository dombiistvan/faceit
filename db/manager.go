package db

import (
	"faceit/common"
	"faceit/db/models"
	"faceit/requests"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type Manager struct {
	gorm  *gorm.DB
	tries uint8
}

// Conn returns *gorm.DB pointer or use passed transaction pointer
func (m Manager) Conn(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return m.gorm
	}
	return tx
}

// return connection based on config
func getGormConnection(config Config, connectDB bool) (db *gorm.DB, err error) {
	return gorm.Open(mysql.New(mysql.Config{
		DSN:                      BuildMysqlDSN(config, connectDB),
		DisableDatetimePrecision: true,
	}), &gorm.Config{})
}

// BuildMysqlDSN build and return connection string based on config
func BuildMysqlDSN(config Config, connectDB bool) string {
	var host, db, param string
	var params []string

	db = "/"
	if connectDB && !common.IsEmptyString(config.DBName) {
		db = fmt.Sprintf("/%s", config.DBName)
	}

	host = fmt.Sprintf("%s:%s@%s(%s)", config.DBUser, config.DBPassword, "tcp", config.DBHost)
	if !common.IsEmptyString(config.Charset) {
		params = append(params, fmt.Sprintf("charset=%s", config.Charset))
	}
	if !common.IsEmptyString(config.Loc) {
		params = append(params, fmt.Sprintf("loc=%s", config.Loc))
	}
	if config.ParseTime {
		params = append(params, fmt.Sprintf("parseTime=%v", config.ParseTime))
	}

	if len(params) > 0 {
		param = fmt.Sprintf("?%s", strings.Join(params, "&"))
	}

	return fmt.Sprintf("%s%s%s", host, db, param)
}

// GetUserToEmail returns user to email address and error
func (m Manager) GetUserToEmail(email string, tx *gorm.DB) (user models.User, err error) {
	tx = m.Conn(tx)

	err = tx.Where("email = ?", email).First(&user).Error

	return user, err
}

// GetUserToUUID returns user to uuid address and error
func (m Manager) GetUserToUUID(uuid string, tx *gorm.DB) (user models.User, err error) {
	tx = m.Conn(tx)

	err = tx.Where("uuid = ?", uuid).First(&user).Error

	return user, err
}

// ListUsersTo returns user list ListUsersRequest
func (m Manager) ListUsersTo(req requests.ListUsersRequest, tx *gorm.DB) (users []models.User, err error) {
	tx = m.Conn(tx)

	var queryParts = []string{}
	var params = []interface{}{}

	if req.ID != nil {
		queryParts = append(queryParts, "uuid LIKE ?")
		params = append(params, *req.ID)
	}

	if req.Email != nil {
		queryParts = append(queryParts, "email LIKE ?")
		params = append(params, *req.Email)
	}

	if req.FirstName != nil {
		queryParts = append(queryParts, "first_name LIKE ?")
		params = append(params, fmt.Sprintf("%%%s%%", *req.FirstName))
	}

	if req.LastName != nil {
		queryParts = append(queryParts, "last_name LIKE ?")
		params = append(params, fmt.Sprintf("%%%s%%", *req.LastName))
	}

	if req.NickName != nil {
		queryParts = append(queryParts, "nickname LIKE ?")
		params = append(params, fmt.Sprintf("%%%s%%", *req.NickName))
	}

	if req.Country != nil {
		queryParts = append(queryParts, "country LIKE ?")
		params = append(params, fmt.Sprintf("%%%s%%", *req.Country))
	}

	if len(queryParts) > 0 {
		err = tx.Where(strings.Join(queryParts, " AND "), params...).Limit(int(req.GetLimit())).Offset(int(req.GetOffset())).Find(&users).Error
		return users, err
	}

	err = tx.Limit(int(req.GetLimit())).Offset(int(req.GetOffset())).Find(&users).Error
	return users, err
}

// Save saves entity to db and returns nil, or error on failure
func (m Manager) Save(model interface{}, tx *gorm.DB) (err error) {
	tx = m.Conn(tx)

	if err = tx.Save(model).Error; err != nil {
		return err
	}

	tx.Preload(clause.Associations).Find(model)

	return tx.Error
}

// Delete tries to delete struct record from DB and/or returns error
func (m Manager) Delete(model interface{}, tx *gorm.DB) error {
	tx = m.Conn(tx)
	return tx.Unscoped().Delete(model).Error
}
