package go_mysql_handle

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"regexp"
	"strings"
	"time"
)

type Mysql struct {
	Conn *sql.DB
	// config
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func (m *Mysql) GetDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.DBName)
	log.Println(dsn)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		m.HandleMysqlError(err)
		return err
	}
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)
	m.Conn = conn
	return nil
}

func (m *Mysql) HandleMysqlError(err error) {
	log.Println(err)
}
func (m *Mysql) ParseFromJDBCURL(jdbcUrl string) {
	driverName := getDriverName(jdbcUrl)
	log.Println("driver name is :", driverName)
	dbDriverName := getDBDriverName(jdbcUrl)
	log.Println("db driver name is :", dbDriverName)
	dbHost := getDBHost(jdbcUrl)
	log.Println("db host is :", dbHost)
	dbPort := getDBPort(jdbcUrl)
	log.Println("db host is :", dbPort)
	dbName := getDBName(jdbcUrl)
	log.Println("db name is :", dbName)
	m.Port = dbPort
	m.Host = dbHost
	m.DBName = dbName
}

func getDriverName(jdbcUrl string) string {
	driverNameRegexp, _ := regexp.Compile("(.*?):")
	s := driverNameRegexp.FindStringSubmatch(jdbcUrl)
	return checkRegexpResult(s)
}
func getDBDriverName(jdbcUrl string) string {
	driverNameRegexp, _ := regexp.Compile(":(.*?):")
	s := driverNameRegexp.FindStringSubmatch(jdbcUrl)
	return checkRegexpResult(s)
}
func getDBHost(jdbcUrl string) string {
	dbHostRegexp, _ := regexp.Compile("://(.*?):")
	s := dbHostRegexp.FindStringSubmatch(jdbcUrl)
	return checkRegexpResult(s)
}
func getDBPort(jdbcUrl string) string {
	dbHostRegexp, _ := regexp.Compile(":([0-9].*?)/")
	s := dbHostRegexp.FindStringSubmatch(jdbcUrl)
	return checkRegexpResult(s)
}

func getDBName(jdbcUrl string) string {
	r, _ := regexp.Compile("[0-9].*/(.*?)\\?")
	strArray := r.FindStringSubmatch(jdbcUrl)
	return checkRegexpResult(strArray)
}

func checkRegexpResult(array []string) string {
	if len(array) == 0 {
		//	get no result
		return ""
	}
	return array[1]
}
func (m *Mysql) ParseDBNameFromJdbcURL(jdbcUrl string) string {
	r, _ := regexp.Compile(":[0-9](.*)\\?")
	strArray := r.FindStringSubmatch(jdbcUrl)
	if len(strArray) == 0 {
		// not find
		return ""
	}
	dbName := strArray[1]
	if index := strings.Index(dbName, "/"); index >= 0 {
		return string([]rune(dbName)[(index + 1):])
	}
	return ""
}
