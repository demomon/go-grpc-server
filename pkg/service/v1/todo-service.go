package v1

import (
	"context"
	"fmt"
	"github.com/demomon/go-grpc-server/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type toDoServiceServer struct {
	toDos map[int64]v1.ToDo
	index int64
}

func (t toDoServiceServer) AddToDo(toDo v1.ToDo) int64 {
	var newId int64
	newId += t.index
	toDo.Id = newId
	t.toDos[newId] = toDo
	return newId
}

func (t toDoServiceServer) GetToDo(index int64) v1.ToDo {
	return t.toDos[index]
}

func (t toDoServiceServer) UpdateToDo(index int64, toDo v1.ToDo) int64 {
	_, ok := t.toDos[index]
	if ok {
		t.toDos[index] = toDo
		return 1
	}
	return 0
}

func (t toDoServiceServer) DeleteToDo(index int64) int64 {
	_, ok := t.toDos[index]
	if ok {
		delete(t.toDos, index)
		return 1
	}
	return 0
}

func (t toDoServiceServer) GetAllToDos() []*v1.ToDo {
	list := []*v1.ToDo{}
	for _, toDo := range t.toDos {
		list = append(list, &toDo)
	}
	return list
}

// NewToDoServiceServer creates ToDo service
func NewToDoServiceServer() v1.ToDoServiceServer {
	toDos := make(map[int64]v1.ToDo)
	return &toDoServiceServer{
		toDos: toDos,
		index: 0,
	}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *toDoServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// Create new todo task
func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	_, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	// insert ToDo entity data
	id := s.AddToDo(*req.ToDo)
	log.Printf("Added TODO with id: %v\n", id)

	return &v1.CreateResponse{
		Api: apiVersion,
		Id:  id,
	}, nil
}

// Read todo task
func (s *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get ToDo data
	toDo := s.GetToDo(req.Id)

	return &v1.ReadResponse{
		Api:  apiVersion,
		ToDo: &toDo,
	}, nil

}

// Update todo task
func (s *toDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	if req.ToDo.Id < 0 {
		return nil, fmt.Errorf("Not an existing ToDo")
	}

	updateCount := s.UpdateToDo(req.ToDo.Id, *req.ToDo)

	return &v1.UpdateResponse{
		Api:     apiVersion,
		Updated: updateCount,
	}, nil
}

// Delete todo task
func (s *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	deleteCount := s.DeleteToDo(req.Id)

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: deleteCount,
	}, nil
}

// Read all todo tasks
func (s *toDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	list := s.GetAllToDos()

	return &v1.ReadAllResponse{
		Api:   apiVersion,
		ToDos: list,
	}, nil
}
