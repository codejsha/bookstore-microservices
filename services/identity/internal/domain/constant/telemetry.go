package constant

type TracerName string

const (
	TracerNameGinServer TracerName = "bookstore.identity.gin.server"
	TracerNameGorm      TracerName = "bookstore.identity.gorm"
	TracerNameGrpc      TracerName = "bookstore.identity.grpc"
	TracerNameKafka     TracerName = "bookstore.identity.kafka"
	TracerNameS3        TracerName = "bookstore.identity.s3"
	TracerNameRedis     TracerName = "bookstore.identity.redis"
)
