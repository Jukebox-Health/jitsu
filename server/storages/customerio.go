package storages

import (
	_ "embed"
	"fmt"

	"github.com/jitsucom/jitsu/server/adapters"
)

//go:embed transform/customerio.js
var customerIOTransform string

//CustomerIO is a destination that can send data into CustomerIO
type CustomerIO struct {
	HTTPStorage
}

func init() {
	RegisterStorage(StorageType{typeName: CustomerIOType, createFunc: NewCustomerIO, isSQL: false})
}

//NewCustomerIO returns configured CustomerIO destination
func NewCustomerIO(config *Config) (storage Storage, err error) {
	defer func() {
		if err != nil && storage != nil {
			storage.Close()
			storage = nil
		}
	}()
	if !config.streamMode {
		return nil, fmt.Errorf("NewCustomerIO destination doesn't support %s mode", BatchMode)
	}
	customerIOConfig := &adapters.CustomerIOConfig{}
	if err = config.destination.GetDestConfig(config.destination.CustomerIO, customerIOConfig); err != nil {
		return
	}

	cio := &CustomerIO{}
	err = cio.Init(config, cio, customerIOTransform, `return toCustomerIO($)`)
	if err != nil {
		return
	}
	storage = cio

	requestDebugLogger := config.loggerFactory.CreateSQLQueryLogger(config.destinationID)
	aAdapter, err := adapters.NewCustomerIO(customerIOConfig, &adapters.HTTPAdapterConfiguration{
		DestinationID:  config.destinationID,
		Dir:            config.logEventPath,
		HTTPConfig:     DefaultHTTPConfiguration,
		QueueFactory:   config.queueFactory,
		PoolWorkers:    defaultWorkersPoolSize,
		DebugLogger:    requestDebugLogger,
		ErrorHandler:   cio.ErrorEvent,
		SuccessHandler: cio.SuccessEvent,
	})
	if err != nil {
		return
	}
	//HTTPStorage
	cio.adapter = aAdapter

	//streaming worker (queue reading)
	cio.streamingWorker = newStreamingWorker(config.eventQueue, cio)
	return
}

//Type returns CustomerIO type
func (a *CustomerIO) Type() string {
	return CustomerIOType
}
