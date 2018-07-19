package commons

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	consMysqlHelper *mysqlHelper
	isDbug          bool
)

type mysqlHelper struct {
	MasterDB, SlaveDb *gorm.DB
}

type SQLRule struct {
	Table  string
	Select []string
	Where  *Conditions
	Limit  interface{}
	OffSet interface{}
	Group  []string
	Order  map[string]string
	Having *Conditions
	Joins  []*Join
}

type Conditions struct {
	Sentence   string
	Parameters []interface{}
}

type Join struct {
	Direction string
	JoinTable string
	On        string
}

func GetMysqlHelper() *mysqlHelper {
	if consMysqlHelper == nil {
		consMysqlHelper = &mysqlHelper{}
		consMysqlHelper.MysqlInit()
	}
	return consMysqlHelper
}

// type mysqlGetConfig() struct {
// 	Master, Slave, DataBase, UserName, PassWord string
// }

// @Title MysqlInit
// @Description get MasterDB & SlaveDB
// @Parameters
// @Returns
func (m *mysqlHelper) MysqlInit() {
	var (
		err               error
		masterDB, slaveDB *gorm.DB
	)
	master := GetConfig().GetString("mysql.master")
	slave := GetConfig().GetString("mysql.slave")
	dataBase := GetConfig().GetString("mysql.database")
	userName := GetConfig().GetString("mysql.username")
	passWord := GetConfig().GetString("mysql.password")
	isDebugStr := GetConfig().GetString("mysql.isdebug")
	isDbug, err = strconv.ParseBool(isDebugStr)
	if err != nil {
		isDbug = true
	}
	masterDB, err = getDB(master, userName, passWord, dataBase)
	if err != nil {
		Console().Error(err.Error())
	}
	m.MasterDB = masterDB
	slaveDB, err = getDB(slave, userName, passWord, dataBase)
	if err != nil {
		Console().Error(err.Error())
	}
	m.SlaveDb = slaveDB
}

// @Title getDB
// @Description get *sql.DB by GetConfig()
// @Parameters
//       address         string          host:ip
//       username        string          db username
//       password        string          db password
// @Returns db:*gorm.DB err:error
func getDB(address, userName, passWord, dataBase string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", userName, passWord, address, dataBase))
	db.DB().SetMaxOpenConns(2000)
	db.DB().SetMaxIdleConns(1000)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// func (q *MysqlQuery) Exec(newOrm bool, query string, args ...interface{}) (sql.Result, error) {
// 	db := SlaveDb
// 	if newOrm {
// 		db = ConnectMysql(false)
// 		defer db.Close()
// 	}
// 	orm := beedb.New(db)
// 	return orm.Exec(query, args...)
// }

// @Title FindOne
// @Description get model from mysql by parameters
// @Parameters
//       sqlRule         *SQLRule        sqlrule
//       result          interface{}     return model
// @Returns err:error
func (m *mysqlHelper) FindOne(sqlRule *SQLRule, result interface{}) error {
	var (
		err    error
		gormdb *gorm.DB
	)
	err = checkWhere(sqlRule.Where)
	if err != nil {
		return err
	}
	gormdb, err = setUpGormDB(sqlRule)
	if err != nil {
		return err
	}
	err = gormdb.Where(sqlRule.Where.Sentence, sqlRule.Where.Parameters...).First(result).Error
	if err != nil {
		return err
	}
	return nil
}

// @Title FindAll
// @Description get model list from mysql by parameters
// @Parameters
//       sqlRule         *SQLRule        sqlrule
//       result          []interface{}     return []model
// @Returns err:error
func (m *mysqlHelper) FindAll(sqlRule *SQLRule, result interface{}) error {
	var (
		err    error
		gormdb *gorm.DB
	)
	gormdb, err = setUpGormDB(sqlRule)
	if err != nil {
		return err
	}
	err = checkWhere(sqlRule.Where)
	if err == nil {
		gormdb = gormdb.Where(sqlRule.Where.Sentence, sqlRule.Where.Parameters...)
	}
	err = gormdb.Find(result).Error
	if err != nil {
		return err
	}
	return nil
}

// @Title FindByPaging
// @Description get model list from mysql by parameters(limit & offset)
// @Parameters
//       sqlRule         *SQLRule        sqlrule
//       result          []interface{}     return []model
// @Returns err:error
func (m *mysqlHelper) FindByPaging(sqlRule *SQLRule, result interface{}) error {
	var (
		err    error
		gormdb *gorm.DB
	)
	gormdb, err = setUpGormDB(sqlRule)
	if err != nil {
		return err
	}
	err = checkWhere(sqlRule.Where)
	if err == nil {
		gormdb = gormdb.Where(sqlRule.Where.Sentence, sqlRule.Where.Parameters...)
	}
	err = gormdb.Limit(sqlRule.Limit).Offset(sqlRule.OffSet).Find(result).Error
	if err != nil {
		return err
	}
	return nil
}

// @Title Insert
// @Description add model into mysql by model
// @Parameters
//       sqlRule         *SQLRule        sqlrule
//       model          interface{}      insert model
// @Returns err:error
func (m *mysqlHelper) Insert(sqlRule *SQLRule, model interface{}) error {
	var (
		err    error
		gormdb *gorm.DB
	)
	gormdb, err = setUpGormDB(sqlRule)
	if err != nil {
		return err
	}
	err = gormdb.Save(model).Error
	if err != nil {
		return err
	}
	return nil
}

// @Title Update
// @Description update model by parameters
// @Parameters
//       sqlRule         *SQLRule                   sqlrule
//       u               map[string]interface{}     modify
// @Returns err:error
func (m *mysqlHelper) Update(sqlRule *SQLRule, u map[string]interface{}) error {
	var (
		err    error
		gormdb *gorm.DB
	)
	gormdb, err = setUpGormDB(sqlRule)
	if err != nil {
		return err
	}
	err = checkWhere(sqlRule.Where)
	if err != nil {
		return err
	}
	err = gormdb.Where(sqlRule.Where.Sentence, sqlRule.Where.Parameters...).Updates(u).Error
	if err != nil {
		return err
	}
	return nil
}

// @Title Update
// @Description update model by parameters
// @Parameters
//       sqlRule         *SQLRule                   sqlrule
//       u               map[string]interface{}     modify
// @Returns err:error
func (m *mysqlHelper) Upsert(sqlRule *SQLRule, model interface{}) error {
	var (
		err    error
		gormdb *gorm.DB
	)
	gormdb, err = setUpGormDB(sqlRule)
	if err != nil {
		return err
	}
	err = checkWhere(sqlRule.Where)
	if err != nil {
		return err
	}
	var result []interface{}
	err = m.FindOne(sqlRule, &result)
	if err != nil {
		return err
	}
	if len(result) <= 0 {
		err = m.Insert(sqlRule, model)
	} else {
		err = gormdb.Where(sqlRule.Where.Sentence, sqlRule.Where.Parameters...).Updates(model).Error
	}
	if err != nil {
		return err
	}
	return nil
}

// @Title Delete
// @Description delete model by model
// @Parameters
//       sqlRule         *SQLRule        sqlrule
//       model          interface{}     model
// @Returns err:error
func (m *mysqlHelper) Delete(sqlRule *SQLRule, model interface{}) error {
	var (
		err    error
		gormdb *gorm.DB
	)
	gormdb, err = setUpGormDB(sqlRule)
	if err != nil {
		return err
	}
	err = checkWhere(sqlRule.Where)
	if err != nil {
		return err
	}
	err = gormdb.Where(sqlRule.Where.Sentence, sqlRule.Where.Parameters...).Delete(model).Error
	if err != nil {
		return err
	}
	return nil
}

func checkWhere(where *Conditions) error {
	if where == nil || where.Sentence == "" || where.Parameters == nil {
		return errors.New("undefined where parameters")
	}
	return nil
}

func setUpGormDB(sqlRule *SQLRule) (*gorm.DB, error) {
	var (
		err    error
		gormdb *gorm.DB
	)
	if sqlRule == nil {
		err = errors.New("sqlrule is nil")
		goto ERR
	}
	if sqlRule.Table == "" {
		err = errors.New("table is empty")
		goto ERR
	}
	if GetMysqlHelper().SlaveDb == nil {
		consMysqlHelper = nil
		return nil, errors.New("connection refused")
	}
	gormdb = GetMysqlHelper().SlaveDb.Table(sqlRule.Table)
	if isDbug {
		gormdb = gormdb.Debug()
	}
	if sqlRule.Order != nil {
		for k, v := range sqlRule.Order {
			gormdb = gormdb.Order(fmt.Sprintf("%s %s", k, v))
		}
	}
	if sqlRule.Group != nil && len(sqlRule.Group) > 0 {
		for _, v := range sqlRule.Group {
			if !strings.Contains(v, ",") {
				gormdb = gormdb.Group(v)
			}

		}
	}
	if sqlRule.Having != nil && sqlRule.Having.Parameters != nil && sqlRule.Having.Sentence != "" {
		gormdb = gormdb.Having(sqlRule.Having.Sentence, sqlRule.Having.Parameters...)
	}
	if sqlRule.Joins != nil && len(sqlRule.Joins) > 0 {
		for _, v := range sqlRule.Joins {
			if v.Direction != "" && v.JoinTable != "" && v.On != "" {
				gormdb = gormdb.Joins(fmt.Sprintf("%s join %s on %s", v.Direction, v.JoinTable, v.On))
			}
		}
	}
	if sqlRule.Select != nil && len(sqlRule.Select) > 0 {
		gormdb = gormdb.Select(sqlRule.Select)
	}
	goto RETURN
ERR:
	return gormdb, err
RETURN:
	return gormdb, err
}
