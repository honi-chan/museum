package handler

import "museum/api/generated"

// Composition delegates generated API methods to their feature handlers.
// Business rules remain in the use cases.
type APIHandler struct {
	*ExhibitionHandler
	*ExhibitHandler
}

func NewAPIHandler(exhibitions *ExhibitionHandler, exhibits *ExhibitHandler) *APIHandler {
	return &APIHandler{ExhibitionHandler: exhibitions, ExhibitHandler: exhibits}
}

var _ generated.StrictServerInterface = (*APIHandler)(nil)
