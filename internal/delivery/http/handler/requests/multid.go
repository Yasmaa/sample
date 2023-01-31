package requests

type MultiID struct {
	Ids []int `json:"ids" validate:"required"`
}