package constant

type TracerName string

const (
	TracerNameGinServer TracerName = "bookstore.inventory.gin.server"
	TracerNameGorm      TracerName = "bookstore.inventory.gorm"
	TraceNameKafka      TracerName = "bookstore.inventory.kafka"
	TraceNameS3         TracerName = "bookstore.inventory.s3"
	TraceNameRedis      TracerName = "bookstore.inventory.redis"
)
