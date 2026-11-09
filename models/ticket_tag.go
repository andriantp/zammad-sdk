package models

type (
	TagList struct {
		Total int   `json:"total"`
		Tags  []Tag `json:"tags"`
	}

	Tag struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
)
