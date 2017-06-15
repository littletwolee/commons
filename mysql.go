package commons

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	consMysqlHelper *mysqlHelper
)

type mysqlHelper struct {
	MasterDB, SlaveDb *gorm.DB
}

type SQLRule struct {
	Table  string
	Where  *Where
	Limit  interface{}
	OffSet interface{}
}

type Where struct {
	Sentence   string
	Parameters []interface{}
}

func GetMysqlHelper() *mysqlHelper {
	if consMysqlHelper == nil {
		consMysqlHelper = &mysqlHelper{}
		consMysqlHelper.MysqlInit()
	}
	return consMysqlHelper
}

// type mysqlConfig struct {
// 	Master, Slave, DataBase, UserName, PassWord string
// }

// @Title MysqlInit
// @Description get MasterDB & SlaveDB
// @Parameters
// @Returns
func (m *mysqlHelper) MysqlInit() {
	master := Config.GetString("mysql.master")
	slave := Config.GetString("mysql.slave")
	dataBase := Config.GetString("mysql.database")
	userName := Config.GetString("mysql.username")
	passWord := Config.GetString("mysql.password")
	masterDB, err := getDB(master, userName, passWord, dataBase)
	if err != nil {
		GetLogger().LogPanic(err)
	}
	m.MasterDB = masterDB
	slaveDB, err := getDB(slave, userName, passWord, dataBase)
	if err != nil {
		GetLogger().LogPanic(err)
	}
	m.SlaveDb = slaveDB
}

// @Title getDB
// @Description get *sql.DB by config
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
	err := m.SlaveDb.Table(sqlRule.Table).Where(sqlRule.Where.Sentence, sqlRule.Where.Parameters...).First(result).Error
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
	err := m.SlaveDb.Table(sqlRule.Table).Where(sqlRule.Where.Sentence, sqlRule.Where.Parameters...).Find(result).Error
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
	err := m.SlaveDb.Table(sqlRule.Table).Where(sqlRule.Where).Limit(sqlRule.Limit).Offset(sqlRule.OffSet).Find(result).Error
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
	err := m.SlaveDb.Table(sqlRule.Table).Save(model).Error
	if err != nil {
		return err
	}
	return nil
}

// @Title Update
// @Description update model by parameters
// @Parameters
//       sqlRule         *SQLRule        sqlrule
//       result          []interface{}     return []model
// @Returns
func (m *mysqlHelper) Update(sqlRule *SQLRule, v interface{}, u map[string]interface{}) error {
	err := m.MasterDB.Table(sqlRule.Table).Model(v).Where(sqlRule.Where.Sentence, sqlRule.Where.Parameters...).Update(u).Error
	if err != nil {
		return err
	}
	return nil
}

// @Title Delete
// @Description delete model by model
// @Parameters
//       sqlRule         *SQLRule        sqlrule
//       model          []interface{}     return []model
// @Returns
func (m *mysqlHelper) Delete(sqlRule *SQLRule, model interface{}) error {
	err := m.MasterDB.Table(sqlRule.Table).Delete(model).Error
	if err != nil {
		return err
	}
	return nil
}

// // @Title name
// // @Description
// // @Parameters
// //
// // @Returns err:error
// func (m *MysqlHelper) SQLInjectionAttack(value string) (string, error) {
// 	var (
// 		re  *regexp.Regexp
// 		err error
// 	)
// 	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
// 	re, err = regexp.Compile(str)
// 	if err != nil {
// 		value = ""
// 		goto RESULT
// 	}
// 	if value == "" {
// 		err = errors.New(PARAMETER_EMPTY)
// 		goto RESULT
// 	}
// 	if re.MatchString(value) {
// 		err = errors.New(fmt.Sprintf("%s:%s", PARAMETER_SQL_ATTACK, value))
// 		value = ""
// 		goto RESULT
// 	}
// 	goto RESULT
// RESULT:
// 	return value, err
// }
