package constant

type TracerName string

const (
	TracerNameGinServer TracerName = "bookstore.catalog.gin.server"
	TracerNameGorm      TracerName = "bookstore.catalog.gorm"
	TracerNameGrpc      TracerName = "bookstore.catalog.grpc"
	TracerNameConductor TracerName = "bookstore.catalog.conductor"
	TraceNameKafka      TracerName = "bookstore.catalog.kafka"
	TraceNameS3         TracerName = "bookstore.catalog.s3"
	TraceNameRedis      TracerName = "bookstore.catalog.redis"
)
