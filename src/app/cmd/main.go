package main

import (
	"log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kuzkuss/VK_DB_Project/app/cmd/server"

	forumDelivery "github.com/kuzkuss/VK_DB_Project/app/internal/forum/delivery"
	forumRep "github.com/kuzkuss/VK_DB_Project/app/internal/forum/repository"
	forumUsecase "github.com/kuzkuss/VK_DB_Project/app/internal/forum/usecase"
	postDelivery "github.com/kuzkuss/VK_DB_Project/app/internal/post/delivery"
	postRep "github.com/kuzkuss/VK_DB_Project/app/internal/post/repository"
	postUsecase "github.com/kuzkuss/VK_DB_Project/app/internal/post/usecase"
	serviceDelivery "github.com/kuzkuss/VK_DB_Project/app/internal/service/delivery"
	serviceRep "github.com/kuzkuss/VK_DB_Project/app/internal/service/repository"
	serviceUsecase "github.com/kuzkuss/VK_DB_Project/app/internal/service/usecase"
	threadDelivery "github.com/kuzkuss/VK_DB_Project/app/internal/thread/delivery"
	threadRep "github.com/kuzkuss/VK_DB_Project/app/internal/thread/repository"
	threadUsecase "github.com/kuzkuss/VK_DB_Project/app/internal/thread/usecase"
	userDelivery "github.com/kuzkuss/VK_DB_Project/app/internal/user/delivery"
	userRep "github.com/kuzkuss/VK_DB_Project/app/internal/user/repository"
	userUsecase "github.com/kuzkuss/VK_DB_Project/app/internal/user/usecase"
)

var cfgPg = postgres.Config{DSN: "host=localhost user=db_pg password=db_postgres database=db_forum port=5432"}

func main() {
	db, err := gorm.Open(postgres.New(cfgPg),
		&gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	forumDB := forumRep.New(db)
	userDB := userRep.New(db)
	postDB := postRep.New(db)
	threadDB := threadRep.New(db)
	serviceDB := serviceRep.New(db)

	forumUC := forumUsecase.New(forumDB, userDB)
	userUC := userUsecase.New(userDB)
	postUC := postUsecase.New(postDB, userDB, threadDB, forumDB)
	threadUC := threadUsecase.New(threadDB, userDB, forumDB)
	serviceUC := serviceUsecase.New(serviceDB)

	e := echo.New()

	e.Logger.SetHeader(`time=${time_rfc3339} level=${level} prefix=${prefix} ` +
		`file=${short_file} line=${line} message:`)
	e.Logger.SetLevel(elog.INFO)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `time=${time_custom} remote_ip=${remote_ip} ` +
			`host=${host} method=${method} uri=${uri} user_agent=${user_agent} ` +
			`status=${status} error="${error}" ` +
			`bytes_in=${bytes_in} bytes_out=${bytes_out}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	e.Use(middleware.Recover())

	forumDelivery.NewDelivery(e, forumUC)
	userDelivery.NewDelivery(e, userUC)
	postDelivery.NewDelivery(e, postUC)
	threadDelivery.NewDelivery(e, threadUC)
	serviceDelivery.NewDelivery(e, serviceUC)

	s := server.NewServer(e)
	if err := s.Start(); err != nil {
		e.Logger.Fatal(err)
	}
}
