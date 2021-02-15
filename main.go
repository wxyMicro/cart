package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"

	"github.com/wxyMicro/cart/domain/repository"
	service2 "github.com/wxyMicro/cart/domain/service"
	"github.com/wxyMicro/cart/handler"
	cart "github.com/wxyMicro/cart/proto/cart"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/wxyMicro/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var QPS = 100

func main() {
	// 配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	//注册中心
	consul := consul2.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	//链路追踪
	tracer, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(tracer)
	//数据库连接
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlInfo.User,
		mysqlInfo.Pwd,
		mysqlInfo.Host,
		mysqlInfo.Port,
		mysqlInfo.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		log.Error(err)
	}
	//初始化表
	//err = repository.NewCartRepository(db).InitTable()
	//if err != nil {
	//	log.Error(err)
	//}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		//暴露的服务地址
		micro.Address("127.0.0.1:8087"),
		//注册中心
		micro.Registry(consul),
		//链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)

	// Initialise service
	service.Init()

	cartRepository := service2.NewCartDataService(repository.NewCartRepository(db))
	// Register Handler
	err = cart.RegisterCartHandler(service.Server(), &handler.Cart{CartDataService: cartRepository})
	if err != nil {
		log.Error(err)
	}
	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
