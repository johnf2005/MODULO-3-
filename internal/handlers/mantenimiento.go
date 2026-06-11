package handlers

import (
	"database/sql"
	"net/http"

	"modulo3_go/internal/models"
)

func (s *Server) listMantenimientos(w http.ResponseWriter, r *http.Request) {
	items, err := s.store.ListMantenimientos()
	if err != nil {
		Error(w, http.StatusInternalServerError, "error al listar mantenimientos")
		return
	}

	JSON(w, http.StatusOK, items)
}

func (s *Server) createMantenimiento(w http.ResponseWriter, r *http.Request) {
	var mantenimiento models.Mantenimiento
	if err := DecodeJSON(r, &mantenimiento); err != nil {
		Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	id, err := s.store.CreateMantenimiento(mantenimiento)
	if err != nil {
		Error(w, http.StatusInternalServerError, "error al crear mantenimiento")
		return
	}

	mantenimiento.ID = id
	JSON(w, http.StatusCreated, mantenimiento)
}

func (s *Server) getMantenimiento(w http.ResponseWriter, r *http.Request, id int) {
	item, err := s.store.GetMantenimiento(id)
	if err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, "mantenimiento no encontrado")
			return
		}
		Error(w, http.StatusInternalServerError, "error al obtener mantenimiento")
		return
	}

	JSON(w, http.StatusOK, item)
}

func (s *Server) updateMantenimiento(w http.ResponseWriter, r *http.Request, id int) {
	var mantenimiento models.Mantenimiento
	if err := DecodeJSON(r, &mantenimiento); err != nil {
		Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	mantenimiento.ID = id
	if err := s.store.UpdateMantenimiento(id, mantenimiento); err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, "mantenimiento no encontrado")
			return
		}
		Error(w, http.StatusInternalServerError, "error al actualizar mantenimiento")
		return
	}

	JSON(w, http.StatusOK, mantenimiento)
}

func (s *Server) deleteMantenimiento(w http.ResponseWriter, r *http.Request, id int) {
	if err := s.store.DeleteMantenimiento(id); err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, "mantenimiento no encontrado")
			return
		}
		Error(w, http.StatusInternalServerError, "error al eliminar mantenimiento")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
