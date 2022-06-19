package model

type DomainError struct {
	Message string `json:"Message"`
}

func (e *DomainError) Error() string {
	return e.Message
}

type MalformedRequest struct {
	Status      int
	Diagnostics DomainError
}

func (mr *MalformedRequest) Error() string {
	return mr.Diagnostics.Message
}
