package twigo

type MetaEntity struct {
	ResultCount int
	NewestID string
	OldestID string
	NextToken string
}

type Tweet struct {
	ID string
	Text string
}

type User struct {
	ID string
}

// type Space struct {}
// type List struct {}

type ErrorEntity struct {
	Parameters map[string]interface{}
	Message string
}

