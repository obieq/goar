package goar_test_models

type Vehicle struct {
	Year  int    `json:"year,omitempty"`
	Make  string `json:"make,omitempty"`
	Model string `json:"model,omitempty"`
}

type Automobile struct {
	Vehicle
}

type Motorcycle struct {
	Vehicle
}
