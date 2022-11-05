package main

import (
	"UsersBalanceWorker/entities"
	"UsersBalanceWorker/internal/balanceWorker"
	"UsersBalanceWorker/internal/conn"
	"UsersBalanceWorker/internal/db"
	"UsersBalanceWorker/internal/helpers/bind"
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

	Db := &balanceWorker.DB{
		Users:    db.UsersInstance{Db: users},
		Services: db.ServicesInstance{Db: services},
		Record:   db.RecordInstance{Db: record},
	}

	r := gin.Default()

	r.Static("/records", "records")

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

		if err := bind.RightBindedCredit(cr); err != nil {
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
			"msg": "ok",
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

		if err := bind.RightBindedBalance(b); err != nil {
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

		if err := bind.RightBindedService(s); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := balanceWorker.Service(s, Db); err != nil {
			switch {
			case errors.Is(err, db.NotEnoughMoney):
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": db.NotEnoughMoney.Error(),
				})

			case errors.Is(err, db.DuplicateError):
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": db.DuplicateError.Error(),
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": err.Error(),
				})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})

	})

	r.POST("/transfer", func(c *gin.Context) {
		var t entities.Transfer
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := bind.RightBindedTransfer(t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := balanceWorker.Transfer(t, Db); err != nil {
			switch {
			case errors.Is(err, db.NotEnoughMoney):
				c.JSON(http.StatusOK, gin.H{
					"msg": db.NotEnoughMoney.Error(),
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": err.Error(),
				})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	})

	r.POST("/orderStatus", func(c *gin.Context) {
		var ord entities.OrderStatus
		if err := c.ShouldBindJSON(&ord); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := bind.RightBindedStatus(ord); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := balanceWorker.OrderStatus(ord, Db); err != nil {
			switch {
			case errors.Is(err, db.DoesNotExistError):
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": db.DoesNotExistError.Error(),
				})

			case errors.Is(err, db.DoneOrderError):
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": db.DoneOrderError.Error(),
				})

			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": err.Error(),
				})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})

	})

	r.POST("/record", func(c *gin.Context) {
		var rec entities.Record
		if err := c.ShouldBindJSON(&rec); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := bind.RightBindedRecord(rec); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		link, err := balanceWorker.Record(rec, Db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":    "ok",
			"record": link,
		})
	})

	r.POST("/history", func(c *gin.Context) {
		var h entities.History
		if err := c.ShouldBindJSON(&h); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if err := bind.RightBindedHistory(h); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		userHistory, err := balanceWorker.History(h, Db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":          "ok",
			"user_history": userHistory,
		})
	})

	log.Fatal(r.Run())
}
