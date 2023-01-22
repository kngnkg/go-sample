package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func StatusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}

// DBのトランザクションをセットアップする
func DBTransactionMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		txHandle, err := db.BeginTxx(c.Request.Context(), nil)
		if err != nil {
			log.Printf("failed to begin transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			c.Abort()
			return
		}
		log.Print("beginning database transaction")

		defer func() {
			if r := recover(); r != nil {
				txHandle.Rollback()
			}
		}()

		c.Set("db_trx", txHandle)
		c.Next()

		if StatusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			log.Print("committing transactions")
			if err := txHandle.Commit(); err != nil {
				log.Print("trx commit error: ", err)
			}
		} else {
			log.Print("rolling back transaction due to status code: ", c.Writer.Status())
			txHandle.Rollback()
		}
	}
}
