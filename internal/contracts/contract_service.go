package contracts

import (
	store "github.com/Barms1218/nagl/internal"
)

type ContractService struct {
	store.Store
}

func NewContractService(s *store.Store) *ContractService {
	return &ContractService{store: s}
}
