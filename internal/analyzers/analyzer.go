package analyzers

import (
	"context"
)

type Analyzer interface {
	Run(context.Context)
}
