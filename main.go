package main

import (
	"fmt"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"user/common"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"user/domain/repository"
	service2 "user/domain/service"
	"user/handler"
	user "user/proto/user"
)

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("192.168.119.128", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
		return
	}
	// 注册中心
	consulRegister := consul.NewRegistry(
		//
		func(options *registry.Options) {
			options.Addrs = []string{
				"192.168.119.128:8500",
			}
		},
	)
	// Create service
	srv := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8082"),
		// 添加注册中心
		micro.Registry(consulRegister),
	)
	// 初始化服务
	srv.Init()

	// 获取mysql配置
	mysqlConfig := common.GetMysqlConfigFromConsul(consulConfig, "mysql")
	fmt.Println(mysqlConfig)

	// 创建数据库链接
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:123456@tcp(192.168.119.128:3306)/go-micro-demo?charset=utf8&parseTime=True&loc=Local", // DSN data source name
		DefaultStringSize:         256,                                                                                         // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                                                                                        // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                                                                                        // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                                                                                        // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                                                                                       // 根据版本自动配置
	}), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	idb, err := db.DB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer idb.Close()

	rp := repository.NewUserRepository(db)
	_ = rp.InitTable()

	userDataService := service2.NewUserDataService(rp)
	err = user.RegisterUserHandler(srv.Server(), &handler.User{UserDataService: userDataService})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
