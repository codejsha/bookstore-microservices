package constant

type TracerName string

const (
	TracerNameGinServer TracerName = "bookstore.customer.gin.server"
	TracerNameGorm      TracerName = "bookstore.customer.gorm"
	TracerNameGrpc      TracerName = "bookstore.customer.grpc"
	TracerNameKafka     TracerName = "bookstore.customer.kafka"
	TracerNameS3        TracerName = "bookstore.customer.s3"
	TracerNameRedis     TracerName = "bookstore.customer.redis"
)
