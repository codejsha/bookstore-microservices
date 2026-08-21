package object

type Env string

const (
	Local       Env = "local"
	DevCommand  Env = "dev-command"
	DevQuery    Env = "dev-query"
	QaCommand   Env = "qa-command"
	QaQuery     Env = "qa-query"
	StgCommand  Env = "stg-command"
	StgQuery    Env = "stg-query"
	ProdCommand Env = "prod-command"
	ProdQuery   Env = "prod-query"
)
