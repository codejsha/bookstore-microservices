package constant

type TracerName string

const (
	TracerNameGinServer TracerName = "bookstore.catalog.gin.server"
	TracerNameGorm      TracerName = "bookstore.catalog.gorm"
	TracerNameKafka     TracerName = "bookstore.catalog.kafka"
	TracerNameS3        TracerName = "bookstore.catalog.s3"
	TracerNameRedis     TracerName = "bookstore.catalog.redis"
)
