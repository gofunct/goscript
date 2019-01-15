// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package gopher

import (
	"github.com/gofunct/goscript/script"
	"github.com/gofunct/iio"
)

// Injectors from inject.go:

func Inject() script.Executor {
	service := iio.NewStdIO()
	gopher := New(service)
	return gopher
}