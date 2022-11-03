package main

import (
	"UsersBalanceWorker/entities"
	"UsersBalanceWorker/internal/balanceWorker"
	"UsersBalanceWorker/internal/conn"
	"UsersBalanceWorker/internal/db"
	"UsersBalanceWorker/pkg/helpers/bind"
	"UsersBalanceWorker/pkg/helpers/pg"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	users, err := conn.Connection(&pg.Config{
		Host:     "localhost",
		Username: "db_user",
		Password: "pwd123",
		Port:     "54320",
		DbName:   "users",
		Timeout:  5,
	})

	if err != nil {
		log.Fatal(err)
	}

	services, err := conn.Connection(&pg.Config{
		Host:     "localhost",
		Username: "db_user",
		Password: "pwd123",
		Port:     "54320",
		DbName:   "services",
		Timeout:  5,
	})

	if err != nil {
		log.Fatal(err)
	}

	record, err := conn.Connection(&pg.Config{
		Host:     "localhost",
		Username: "db_user",
		Password: "pwd123",
		Port:     "54320",
		DbName:   "record",
		Timeout:  5,
	})

	if err != nil {
		log.Fatal(err)
	}

	Db := &entities.Db{
		Users:    db.UsersInstance{Db: users},
		Services: db.ServicesInstance{Db: services},
		Record:   db.RecordInstance{Db: record},
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/credit", func(c *gin.Context) {
		cr := entities.Credit{}
		if err := c.ShouldBindJSON(&cr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := bind.IsBindedCredit(cr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := balanceWorker.Credit(cr, Db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "OK",
		})
	})

	r.POST("/balance", func(c *gin.Context) {
		var b entities.Balance
		if err := c.ShouldBindJSON(&b); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := bind.IsBindedBalance(b); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		balance, err := balanceWorker.Balance(b, Db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":     "ok",
			"balance": balance,
		})
	})

	r.POST("/service", func(c *gin.Context) {
		var s entities.Service
		if err := c.ShouldBindJSON(&s); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := bind.IsBindedService(s); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := balanceWorker.Service(s, Db); err != nil {
			if errors.Is(err, db.NotEnoughMoney) {
				c.JSON(http.StatusOK, gin.H{
					"msg": db.NotEnoughMoney.Error(),
				})
				return
			}

			if errors.Is(err, db.DuplicateError) {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "duplicates: order with this id already exists",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})

	})

	r.POST("/transfer", func(c *gin.Context) {
		// TODO: transfer user to user
	})

	r.POST("/orderStatus", func(c *gin.Context) {
		// TODO order status failed
	})

	r.POST("/record", func(c *gin.Context) {
		// TODO get record for needed period
	})

	r.POST("/history", func(c *gin.Context) {
		// TODO get user payments history for needed period
	})

	log.Fatal(r.Run())
}
