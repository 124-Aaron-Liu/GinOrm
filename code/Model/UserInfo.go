package Model

import (
	"database/sql/driver"
	"fmt"
	"time"

	//V2需要引用這package
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type LocalTime time.Time

type UserInfo struct {
	UID        int       `gorm:"type:bigint(20) NOT NULL auto_increment;primary_key;" json:"uid,omitempty"`
	Username   string    `gorm:"type:varchar(40) NOT NULL;" json:"username,omitempty"`
	Department string    `gorm:"type:varchar(40) NOT NULL;" json:"department,omitempty"`
	Created    LocalTime `gorm:"type:timestamp" json:"created,omitempty" time_format:"2006-01-02 15:04:05"`
	// Username   string
	// Department string
	// Created    LocalTime
}

// type User struct {
// 	ID         int64     `gorm:"type:bigint(20) NOT NULL auto_increment;primary_key;" json:"id,omitempty"`
// 	Username   string    `gorm:"type:varchar(40) NOT NULL;" json:"username,omitempty"`
// 	Department string    `gorm:"type:varchar(40) NOT NULL;" json:"department,omitempty"`
// 	Created    LocalTime `gorm:"type:timestamp" json:"created,omitempty" time_format:"2006-01-02 15:04:05"`
// }

const TimeFormat = "2006-01-02 15:04:05"

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return
	}

	// 指定解析的格式
	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = LocalTime(now)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

// 写入 mysql 时调用
func (t LocalTime) Value() (driver.Value, error) {
	// 0001-01-01 00:00:00 属于空值，遇到空值解析成 null 即可
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

// 检出 mysql 时调用
func (t *LocalTime) Scan(v interface{}) error {
	// mysql 内部日期的格式可能是 2006-01-02 15:04:05 +0800 CST 格式，所以检出的时候还需要进行一次格式化
	tTime, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", v.(time.Time).String())
	*t = LocalTime(tTime)
	return nil
}

// 用于 fmt.Println 和后续验证场景
func (t LocalTime) String() string {
	return time.Time(t).Format(TimeFormat)
}

const (
	UserName     string = "root"
	Password     string = "123456"
	Addr         string = "ginorm-mysql-1"
	Port         int    = 3306
	Database     string = "test"
	MaxLifetime  int    = 10
	MaxOpenConns int    = 10
	MaxIdleConns int    = 10
)

func getDB() *gorm.DB {
	//組合sql連線字串
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True", UserName, Password, Addr, Port, Database)
	//連接MySQL
	conn, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	checkErr(err)

	db, err := conn.DB()
	checkErr(err)

	db.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)

	return conn
}

func Create(name string, department string, created LocalTime) error {

	conn := getDB()

	user := UserInfo{Username: name, Department: department, Created: created}
	//
	result := conn.Debug().Create(&user)

	if result.Error != nil {
		fmt.Println("Create failt")
	}
	if result.RowsAffected != 1 {
		fmt.Println("RowsAffected Number failt")
	}

	return nil
}

func GetUserInfo(id int) (UserInfo, error) {
	conn := getDB()

	// var users []UserInfo
	var user UserInfo
	conn.Debug().Where("uid = ?", id).Find(&user)

	// res := conn.Debug().Find(&users)
	// fmt.Println(res.RowsAffected)
	// checkErr(err)
	// userInfo := UserInfo{}
	// row := db.QueryRow("select username, department, created FROM userinfo WHERE uid=?", id)

	// //Scan對應的欄位與select語法的欄位順序一致
	// if err := row.Scan(&userInfo.Username, &userInfo.Department, &userInfo.Created); err != nil {
	// 	fmt.Printf("scan failed, err:%v\n", err)
	// 	return UserInfo{}, err
	// }

	return user, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
