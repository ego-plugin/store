package eqmgo

import (
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"
)

func ScanToM(v struct{}) (m bson.M) {
	s := structs.New(v)
	s.TagName = "bson"
	m = s.Map()
	return m
}
