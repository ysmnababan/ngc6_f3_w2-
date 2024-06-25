deps:
	go get -u golang.org/x/crypto/bcrypt
	go get -u github.com/joho/godotenv
	go get -u github.com/golang-jwt/jwt/v4
	go get -u gorm.io/gorm
	go get github.com/labstack/echo/v4
	go get -u github.com/labstack/echo/v4/middleware
	go get -u github.com/sirupsen/logrus
	go get go.mongodb.org/mongo-driver/mongo
	go get go.mongodb.org/mongo-driver/mongo/options
	go get -u github.com/robfig/cron
	go get -u google.golang.org/grpc

.PHONY: all
all: deps