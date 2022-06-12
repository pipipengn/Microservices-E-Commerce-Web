package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

var resName = "goods-list"

func InitSentinel() {
	if err := sentinel.InitDefault(); err != nil {
		zap.S().Error()
	}

	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               resName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              3,
			StatIntervalInMs:       6000,
		},
	})
	if err != nil {
		zap.S().Errorf("Unexpected error: %+v", err)
		return
	}
}
