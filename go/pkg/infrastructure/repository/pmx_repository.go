package repository

import (
	"github.com/miu200521358/mmd-auto-trace-5/pkg/domain/pmx"
)

type PmxRepository struct {
	*baseRepository[*pmx.PmxModel]
	isLog bool
}

func NewPmxRepository(isLog bool) *PmxRepository {
	return &PmxRepository{
		baseRepository: &baseRepository[*pmx.PmxModel]{
			newFunc: func(path string) *pmx.PmxModel {
				return pmx.NewPmxModel(path)
			},
		},
		isLog: isLog,
	}
}
