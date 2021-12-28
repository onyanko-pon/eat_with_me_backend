package entity

type Friend struct {
	User     User `json:"user"`
	Blinding bool `json:"blinding"`
}
