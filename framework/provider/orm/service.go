package orm

import (
	"context"
	"github.com/SKYBroGardenLush/skyscraper/framework"
	"github.com/SKYBroGardenLush/skyscraper/framework/contract"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"sync"
	"time"
)

func GetBaseConfig(container framework.Container) *contract.DBConfig {
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	logService := container.MustMake(contract.LogKey).(contract.Log)

	config := &contract.DBConfig{}
	//直接使用配置服务的load方法读取yaml文件
	err := configService.Load("database", config)
	if err != nil {
		//直接使用logService 来打印错误信息
		logService.Error(context.Background(), "parse datatbase config error", nil)
		return nil
	}
	return config

}

// HadeGorm 代表hade框架的orm实现
type HadeGorm struct {
	container framework.Container //服务容器
	dbs       map[string]*gorm.DB //key为dsn,value 为gorm.DB(连接池)

	lock *sync.RWMutex
}

func NewHadeGorm(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	dbs := make(map[string]*gorm.DB)
	lock := &sync.RWMutex{}
	return &HadeGorm{
		container: container,
		dbs:       dbs,
		lock:      lock,
	}, nil
}

func (orm *HadeGorm) GetDB(option ...contract.DBOption) (*gorm.DB, error) {
	logger := orm.container.MustMake(contract.LogKey).(contract.Log)

	//读取默认配置
	config := GetBaseConfig(orm.container)

	logService := orm.container.MustMake(contract.LogKey).(contract.Log)

	//设置logger
	ormLogger := NewOrmLogger(logService)
	config.Config = &gorm.Config{
		Logger: ormLogger,
	}

	//option 对opt进行修改
	for _, opt := range option {
		if err := opt(orm.container, config); err != nil {
			return nil, err
		}
	}

	//如果最终的dsn没有生成dsn
	if config.Dsn == "" {
		dsn, err := config.FormatDsn()
		if err != nil {
			return nil, err
		}
		config.Dsn = dsn
	}

	//判断是否已经实例了gorm.DB
	orm.lock.RLock()
	if db, ok := orm.dbs[config.Dsn]; ok {
		orm.lock.RUnlock()
		return db, nil
	}
	orm.lock.RUnlock()

	//没有实例化gorm.DB ,那么就要进行实例化操作
	orm.lock.Lock()
	defer orm.lock.Unlock()

	//实例化gorm.DB
	var db *gorm.DB
	var err error
	switch config.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.Dsn), config)
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.Dsn), config)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Dsn), config)
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(config.Dsn), config)
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(config.Dsn), config)

	}

	//设置对应连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}

	if config.ConnMaxIdle > 0 {
		sqlDB.SetMaxIdleConns(config.ConnMaxIdle)
	}

	if config.ConnMaxOpen > 0 {
		sqlDB.SetMaxOpenConns(config.ConnMaxOpen)
	}
	if config.ConnMaxLifetime != "" {
		lifeTime, err := time.ParseDuration(config.ConnMaxLifetime)
		if err != nil {
			logger.Error(context.Background(), "conn max lift time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxLifetime(lifeTime)
		}
	}
	if config.ConnMaxIdletime != "" {
		idleTime, err := time.ParseDuration(config.ConnMaxIdletime)
		if err != nil {
			logger.Error(context.Background(), "conn max idle time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxIdleTime(idleTime)
		}
	}
	// 挂载到map中，结束配置
	if err != nil {
		orm.dbs[config.Dsn] = db
	}

	return db, err

}
