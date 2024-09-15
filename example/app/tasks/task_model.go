package tasks

import "github.com/tinh-tinh/mongoose"

type Task struct {
	mongoose.BaseSchema `bson:"inline"`
	Name                string `bson:"name"`
	Status              string `bson:"status"`
	TakeTime            int64  `bson:"takeTime"`
}
