package pmx

import (
	"github.com/miu200521358/mmd-auto-trace-5/pkg/domain/core"
)

// 面リスト
type Faces struct {
	*core.IndexModels[*Face]
}

func NewFaces(capacity int) *Faces {
	return &Faces{
		IndexModels: core.NewIndexModels[*Face](capacity),
	}
}
