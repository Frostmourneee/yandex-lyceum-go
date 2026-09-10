package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	ex2()
}

func ex1() {
	log.SetOutput(os.Stderr)
	log.Println("[ERROR] Something went wrong!")
}

func ex2() {
	log := logrus.New()
	log.Info("Hello, World!")
}

func WriteToLogFile(message, fileName string) error {
	if err := os.WriteFile(fileName, []byte(message), 0600); err != nil {
		return err
	}
	return nil
}

type Order struct {
	OrderNumber  int
	CustomerName string
	OrderAmount  float64
}

type OrderLogger struct{}

func (logger *OrderLogger) AddOrder(order Order) {
	fmt.Printf("Добавлен заказ #%d, Имя клиента: %s, Сумма заказа: $%.2f\n", order.OrderNumber, order.CustomerName, order.OrderAmount)
}

func NewOrderLogger() *OrderLogger {
	return &OrderLogger{}
}

func LogUserAction(logger *slog.Logger, user string, action string) {
	logger.Info("user action", slog.String("User", user), slog.String("Action", action))
}

func LogHTTPRequest(logger *slog.Logger, method, path string, status int, durationMs int64) {
	logger.Info("http request", slog.String("method", method), slog.String("path", path), slog.Int("status", status), slog.Int64("duration_ms", durationMs))
}
