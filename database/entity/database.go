package entity

type DatabaseEntity struct {
	ID             int
	Name           string
	Filepath       string
	Is_initialized int
}

func (e DatabaseEntity) TableName() string {
	return "databases"
}
