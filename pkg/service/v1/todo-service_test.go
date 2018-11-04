package v1

import (
	"context"
	"github.com/demomon/go-grpc-server/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_toDoServiceServer_Create(t *testing.T) {
	context := context.Background()
	tm := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(tm)
	createRequest := &v1.CreateRequest{
		Api: "v1",
		ToDo: &v1.ToDo{
			Title:       "title",
			Description: "description",
			Reminder:    reminder,
		},
	}

	server := NewToDoServiceServer()
	response, err := server.Create(context, createRequest)
	assert.Nil(t, err)
	if assert.NotNil(t, response) {
		// now we know that object isn't nil, we are safe to make
		// further assertions without causing any errors
		var expectedId int64 = 0
		assert.Equal(t, expectedId, response.Id)
	}
}
