package broker_consumer

import "time"

const (
	BrokerSqsType = "SQS"
	BrokerRabbitMQType = "RABBIT"
)

type AssyncProcess struct {
	Hashcode 		string
	BrokerType		string
	Data			interface{}
	Time			time.Time
	Items			[]AssyncProcessItem
}

type AssyncProcessItem struct {
	Status 			string
	Time			time.Time
	Application		string
}

func NewSqsAssyncProcess(hashcode string, application string, data interface{}) * AssyncProcess {
	return &AssyncProcess{
		Hashcode:   hashcode,
		BrokerType: BrokerSqsType,
		Data:       data,
		Time:       time.Time{},
		Items:      []AssyncProcessItem{
			*NewSent(application),
		},
	}
}


func NewSent(application string) *AssyncProcessItem {
	return &AssyncProcessItem{
		Status: "SENT",
		Time: time.Now(),
		Application: application,
	}
}

func Done(application string) *AssyncProcessItem {
	return &AssyncProcessItem{
		Status: "DONE",
		Time: time.Now(),
		Application: application,
	}
}
