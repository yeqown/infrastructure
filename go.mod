module github.com/yeqown/infrastructure

require (
	github.com/AlecAivazis/survey/v2 v2.0.5
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/jinzhu/gorm v1.9.9
	github.com/kr/pretty v0.3.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.1
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	github.com/yeqown/log v0.0.0-20200108034421-d68941cd8fd3
	go.etcd.io/etcd v0.0.0-20181022230727-86b933311d23
	golang.org/x/crypto v0.17.0
	golang.org/x/image v0.4.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)

go 1.13

// replace github.com/yeqown/log => ../log
