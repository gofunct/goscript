// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package terraform

import (
	"github.com/gofunct/goscript/script"
	"github.com/gofunct/iio"
)

// Injectors from inject.go:

func Inject() script.Executor {
	service := iio.NewStdIO()
	terraform := New(service)
	return terraform
}