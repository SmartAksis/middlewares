package broker_consumer

type BrokerJob interface {
	Cron()
	Execute()
}

