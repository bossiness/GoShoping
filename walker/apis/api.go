package main

// API general
type API interface {
	Path() string
	Register()
}
