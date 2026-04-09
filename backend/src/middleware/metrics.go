package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var requestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "nexus_http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "route", "status"},
)

func init() {
	prometheus.MustRegister(requestsTotal)
}

func Metrics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		route := ""
		if c.Route() != nil {
			route = c.Route().Path
		}
		if route == "" {
			route = c.Path()
		}
		requestsTotal.WithLabelValues(c.Method(), route, strconv.Itoa(c.Response().StatusCode())).Inc()
		return err
	}
}
