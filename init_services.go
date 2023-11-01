package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	var (
		user = os.Getenv("db_user")
		pass = os.Getenv("db_password")
		host = os.Getenv("db_host")
		name = os.Getenv("db_name")
	)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, name,
	)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}

	if os.Getenv("db_debug") == "true" {
		db = db.Debug()
	}

	log.Printf("Connected to database '%s' at %s (user %s)", name, host, user)
	return db
}

func initRedis() *redis.Client {
	var (
		addr = os.Getenv("redis_addr")
		pass = os.Getenv("redis_password")
	)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to redis at " + addr)
	return rdb
}
