package service

//go:generate go run github.com/matryer/moq -out moq_test.go . HealthRepository UserRepository

//go:generate go run github.com/matryer/moq -out store_moq_test.go -skip-ensure -pkg service ../store DBConnection
