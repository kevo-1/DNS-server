package models

type Authority struct {
	Username string
	Password string
	Host     string
	Port     string
}

type URL struct {
	Scheme    string
	Authority Authority
	Path      string
	Query     string
	Fragment  string
}