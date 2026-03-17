package dto

type CreateChildRequest struct {
	PhotoURL  string `json:"photo_url"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Hobby     string `json:"hobby"`
	Allergy   string `json:"allergy"`
}

type UpdateChildRequest struct {
	PhotoURL  string `json:"photo_url"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Hobby     string `json:"hobby"`
	Allergy   string `json:"allergy"`
}

type ChildResponse struct {
	ID        int64  `json:"id"`
	PhotoURL  string `json:"photo_url"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Hobby     string `json:"hobby"`
	Allergy   string `json:"allergy"`
}