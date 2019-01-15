//+build wireinject

package protoc

import (
	"github.com/gofunct/goscript/script"
	"github.com/google/wire"
)

func Inject() script.Executor {
	wire.Build(Set)
	return nil
}
