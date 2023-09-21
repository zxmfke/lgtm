package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zxmfke/lgtm/example/gin-lgtm/tracer"
	"net/http"
	"time"
)

func main() {

	if err := tracer.InitTracer(tracer.Config{
		ServiceName: "tracer-demo",
		Endpoint:    "127.0.0.1:4318",
		Sampler:     1.0,
		Batcher:     "otlphttp",
	}); err != nil {
		fmt.Println(fmt.Errorf("%s", err.Error()))
		return
	}

	r := gin.Default()

	r.Use(tracer.Trace)
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/test", func(ctx *gin.Context) {

		Sub(ctx)

		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
		})

	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

func Sub(ctx *gin.Context) {

	newSpan := tracer.NewGinSpan(ctx, "tracer-sub-func")
	defer newSpan.End()

	time.Sleep(time.Second)

	newSpan.SetStringTag("Date", fmt.Sprintf("%s", time.Now().Format(time.DateTime)))
	newSpan.SetInt64Tag("DateTS", time.Now().UnixMilli())

	fmt.Println("go to sub function")
}
