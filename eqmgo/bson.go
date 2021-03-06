package eqmgo

import (
	"github.com/ego-plugin/structs"
	"go.mongodb.org/mongo-driver/bson"
)

func ScanToM(v interface{}) (m bson.M) {
	s := structs.New(v)
	s.TagName = "bson"
	m = s.Map()
	return m
}
