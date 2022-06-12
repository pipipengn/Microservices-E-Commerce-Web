package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"order-web/global"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// tracer config
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: global.ServerConfig.JaegerInfo.Address,
			},
			ServiceName: "shop",
		}

		// new tracer
		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}
		defer func(closer io.Closer) {
			err := closer.Close()
			if err != nil {

			}
		}(closer)

		startSpan := tracer.StartSpan(c.Request.URL.Path)
		defer startSpan.Finish()

		c.Set("tracer", tracer)
		c.Set("parentSapn", startSpan)
		c.Next()
	}
}
