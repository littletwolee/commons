package commons

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	consMysqlHelper *MysqlHelper
)

type MysqlHelper struct {
	MasterDB, SlaveDb *gorm.DB
}

type SQLRule struct {
	Table string
	Where []string
}

func GetMysqlHelper() *MysqlHelper {
	if consMysqlHelper == nil {
		consMysqlHelper = &MysqlHelper{}
	}
	return consMysqlHelper
}

// type mysqlConfig struct {
// 	Master, Slave, DataBase, UserName, PassWord string
// }

func init() {

}

// @Title MysqlInit
// @Description get MasterDB & SlaveDB
// @Parameters
// @Returns
func (m *MysqlHelper) MysqlInit() {
	master := Config.GetString("mysql.master")
	slave := Config.GetString("mysql.slave")
	dataBase := Config.GetString("mysql.database")
	userName := Config.GetString("mysql.username")
	passWord := Config.GetString("mysql.password")
	masterDB, err := getDB(userName, passWord, master, dataBase)
	if err != nil {
		GetLogger().LogPanic(err)
	}
	m.MasterDB = masterDB
	slaveDB, err := getDB(userName, passWord, slave, dataBase)
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
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local", userName, passWord, address, dataBase))
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

func (m *MysqlHelper) FindOne(sqlRule *SQLRule, result interface{}) error {
	m.SlaveDb.Table(sqlRule.Table).Where(sqlRule.Where).First(&result)
	return nil
}

// 	orm := beedb.New(db)
// 	if q.Fields == "" {
// 		q.Fields = "*"
// 	}
// 	return orm.SetTable(q.Table).Where(q.Where).OrderBy(q.OrderBy).Limit(q.Size, q.Offset).Select(q.Fields).Find(result)
// }

// func (q *MysqlQuery) FindAll(result interface{}, newOrm bool) error {
// 	db := SlaveDb
// 	if newOrm {
// 		db = ConnectMysql(false)
// 		defer db.Close()
// 	}

// 	orm := beedb.New(db)
// 	if q.Fields == "" {
// 		q.Fields = "*"
// 	}
// 	return orm.SetTable(q.Table).Where(q.Where).OrderBy(q.OrderBy).Limit(q.Size, q.Offset).Select(q.Fields).FindAll(result)
// }

// func (q *MysqlQuery) Upsert(data interface{}, newOrm bool) error {
// 	db := MasterDB
// 	if newOrm {
// 		db = ConnectMysql(true)
// 		defer db.Close()
// 	}

// 	orm := beedb.New(db)
// 	return orm.SetTable(q.Table).Save(data)
// }

// func (q *MysqlQuery) Delete(newOrm bool) (int64, error) {
// 	db := MasterDB
// 	if newOrm {
// 		db = ConnectMysql(true)
// 		defer db.Close()
// 	}
// 	orm := beedb.New(db)
// 	return orm.SetTable(q.Table).Where(q.Where).DeleteRow()
// }

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
