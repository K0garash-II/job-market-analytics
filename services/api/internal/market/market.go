package market

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("vacancy not found")

type Vacancy struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Location    string     `json:"location"`
	SalaryMin   *int64     `json:"salary_min"`
	SalaryMax   *int64     `json:"salary_max"`
	Currency    string     `json:"currency"`
	Remote      bool       `json:"remote"`
	URL         string     `json:"url"`
	PublishedAt *time.Time `json:"published_at"`
	Company     string     `json:"company"`
	Source      string     `json:"source"`
	Description string     `json:"description,omitempty"`
}
type SearchResult struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Location string  `json:"location"`
	Remote   bool    `json:"remote"`
	URL      string  `json:"url"`
	Company  string  `json:"company"`
	Rank     float64 `json:"rank"`
}
type Filter struct {
	Query, Location, Source string
	Remote                  *bool
	Page, Limit             int
}
type Skill struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Count    int    `json:"count"`
}
type Count struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
type Overview struct {
	Total      int     `json:"total"`
	Remote     int     `json:"remote"`
	Companies  int     `json:"companies"`
	WithSalary int     `json:"with_salary"`
	Days       []Count `json:"days"`
	Sources    []Count `json:"sources"`
}
type Salary struct {
	Total      int     `json:"total_vacancies"`
	WithSalary int     `json:"with_salary"`
	AvgMin     float64 `json:"avg_salary_min"`
	AvgMax     float64 `json:"avg_salary_max"`
}
type Reader interface {
	Vacancies(context.Context, Filter) ([]Vacancy, error)
	Vacancy(context.Context, int) (Vacancy, error)
	Search(context.Context, string) ([]SearchResult, error)
	Skills(context.Context) ([]Skill, error)
	Companies(context.Context) ([]Count, error)
	Overview(context.Context) (Overview, error)
	Salary(context.Context) (Salary, error)
}
