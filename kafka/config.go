package kafka

type Topic string

func (t Topic) String() string {
	return string(t)
}

const (
	CreateToDo Topic = "create-todo"
	DeleteToDo Topic = "delete-todo"
	UpdateToDo Topic = "update-todo"
)
