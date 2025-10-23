package bootstrap

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	_ "modernc.org/sqlite"

	"go.xoder.cn/shortener-server/internal/dal/db/model"
	"go.xoder.cn/shortener-server/internal/shared"
	"go.xoder.cn/shortener-server/internal/utils"
)

// initDB 初始化数据库
func initDB() {
	var err error
	var dialector gorm.Dialector

	dbType := os.Getenv("DATABASE_TYPE")
	if dbType == "" {
		dbType = viper.GetString("database.type")
	}

	switch dbType {
	case "sqlite":
		dialector = connectSqlite()
	case "postgres":
		dialector = connectPostgres()
	case "mysql":
		dialector = connectMysql()
	default:
		panic("database type not support: " + dbType)
	}

	level := viper.GetInt("database.log_level")
	// if level < 0 || level > 4 {
	// 	level = 1
	// }
	gormCfg := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.LogLevel(level)),
	}
	shared.GlobalDB, err = gorm.Open(dialector, gormCfg)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	migrate()
}

// connectSqlite 连接 sqlite
func connectSqlite() gorm.Dialector {
	dsn := viper.GetString("database.sqlite.path")
	// log.Printf("database.sqlite: %s\n", dsn)
	if dsn == "" {
		panic("database.sqlite.path is empty")
	}
	// 如果目录不存在，则创建目录
	if err := utils.MkdirIfNotExist(dsn); err != nil {
		panic("failed to create database file: " + err.Error())
	}

	// 使用 gorm 内置 SQLite 驱动
	// dialector := sqlite.Open(dsn)

	// 使用 modernc.org/sqlite 作为驱动
	dialector := sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dsn, // 数据库文件
	}

	return dialector
}

// connectPostgres 连接 postgres
func connectPostgres() gorm.Dialector {
	host := viper.GetString("database.postgres.host")
	if host == "" {
		host = "localhost"
	}
	port := viper.GetInt("database.postgres.port")
	if port == 0 {
		port = 5432
	}
	sslmode := viper.GetString("database.postgres.sslmode")
	if sslmode == "" {
		sslmode = "disable"
	}
	loc := viper.GetString("database.postgres.timezone")
	if loc == "" {
		loc = "Asia/Shanghai"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		host,
		viper.GetString("database.postgres.user"),
		viper.GetString("database.postgres.password"),
		viper.GetString("database.postgres.database"),
		port,
		sslmode,
		loc,
	)
	// log.Printf("database.postgres: %s\n", dsn)
	return postgres.Open(dsn)
}

// connectMysql 连接 mysql
func connectMysql() gorm.Dialector {
	host := viper.GetString("database.mysql.host")
	if host == "" {
		host = "localhost"
	}
	port := viper.GetInt("database.mysql.port")
	if port == 0 {
		port = 3306
	}
	charset := viper.GetString("database.mysql.charset")
	if charset == "" {
		charset = "utf8mb4"
	}
	parseTime := viper.GetBool("database.mysql.parse_time")
	loc := viper.GetString("database.mysql.loc")
	if loc == "" {
		loc = "Local"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		viper.GetString("database.mysql.user"),
		viper.GetString("database.mysql.password"),
		host,
		port,
		viper.GetString("database.mysql.database"),
		charset,
		parseTime,
		loc,
	)
	// log.Printf("database.mysql: %s\n", dsn)
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	return mysql.Open(dsn)
}

// migrate 数据库迁移 schema
func migrate() {
	// log.Println("migrate")
	err := shared.GlobalDB.AutoMigrate(&model.Url{}, &model.History{})
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	// shared.GlobalDB.Migrator().CurrentDatabase()              // 查看数据库类型
	// shared.GlobalDB.Migrator().GetTables()                    // 查看所有表
	// shared.GlobalDB.Migrator().HasColumn(&model.Urls{}, "id") // 检查字段
}
