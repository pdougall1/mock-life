package main

import (
	"encoding/json"
	"fmt"
)

func CreateCitizen(data citizenData) (citizen Citizen) {
	citizen.Id = data.Id
	citizen.Excitation = data.Excitation
	citizen.Momentum = data.Momentum
	citizen.Leak = 0.2
	return citizen
}

type Citizen struct {
	Id         int
	Excitation int
	Momentum   int
	Leak       float64
}

func (c Citizen) toJson() ([]byte, error) {
	json, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("Citizen json error : %s\n", err)
	}
	return json, err
}

func (c Citizen) update(data citizenData) Citizen {
	c.Excitation += data.Excitation
	c.Momentum += data.Momentum
	fmt.Printf("updating instance? : %d\n", c.Excitation)
	return c
}
