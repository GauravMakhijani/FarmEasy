package services

import (
	"FarmEasy/db"
)

type dependencies struct {
	FarmService Service
}

func InitDependencies() (deps dependencies, err error) {
	store, err := db.Init()
	if err != nil {
		return
	}
	farmService := NewFarmService(store)

	deps = dependencies{
		FarmService: farmService,
	}
	return deps, nil
}
