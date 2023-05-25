package mysqlConnect

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/goalkeeper1983/seakoi/tools"
	"os"
	"time"

	mysql2 "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// CreateMysqlConnect 初始链接的时候需要传入参数 mysqlOption: 0user, 1pass, 2host, 3port, 4dbName, 5charset
func CreateMysqlConnect(mysqlOption ...string) *gorm.DB {
	dns := formatDNS(mysqlOption...)
	gormConfig := &gorm.Config{Logger: &tools.GormLoggerV2{}, NamingStrategy: schema.NamingStrategy{SingularTable: true}}
	return connect(dns, gormConfig)
}

// CreateMysqlConnectByTLS 初始链接的时候需要传入参数 mysqlOption: 0user, 1pass, 2host, 3port, 4dbName, 5charset, 6cert
func CreateMysqlConnectByTLS(mysqlOption ...string) *gorm.DB {
	caCert, err := os.ReadFile(mysqlOption[6])
	if err != nil {
		tools.Log.Panic(tools.RunFuncName(), zap.Any("err", err.Error()))
	}
	tlsConfig := &tls.Config{
		RootCAs:            x509.NewCertPool(),
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}
	if ok := tlsConfig.RootCAs.AppendCertsFromPEM(caCert); !ok {
		tools.Log.Panic(tools.RunFuncName(), zap.Any("err", "ErrorAppendCertsFromPEM"))
	}
	if err = mysql2.RegisterTLSConfig("custom", tlsConfig); err != nil {
		tools.Log.Panic(tools.RunFuncName(), zap.Any("err", err.Error()))
	}
	dns := formatDNS(mysqlOption...)
	dns += "&tls=custom"
	gormConfig := &gorm.Config{Logger: &tools.GormLoggerV2{}, NamingStrategy: schema.NamingStrategy{SingularTable: true}}
	return connect(dns, gormConfig)
}

func connect(dns string, gormConfig *gorm.Config) *gorm.DB {
	var db *gorm.DB
	var err error
	if db, err = gorm.Open(mysql.Open(dns), gormConfig); err != nil {
		tools.Log.Panic(tools.RunFuncName(), zap.Any("err", err.Error()))
	}
	if err = db.Exec("select 1").Error; err != nil {
		tools.Log.Panic(tools.RunFuncName(), zap.Any("err", err.Error()))
	}
	return db
}

func formatDNS(mysqlOption ...string) string {
	mysqlConfig := mysql2.Config{
		User:                 mysqlOption[0],
		Passwd:               mysqlOption[1],
		Addr:                 mysqlOption[2] + ":" + mysqlOption[3],
		DBName:               mysqlOption[4],
		Params:               map[string]string{"charset": mysqlOption[5]},
		Net:                  "tcp",
		Collation:            "utf8mb4_general_ci",
		AllowNativePasswords: true,
		Loc:                  time.Local,
		ParseTime:            true,
	}
	return mysqlConfig.FormatDSN()
}
