package handlers

import (
	"database/sql"
	"net/http"

	"modulo3_go/internal/models"
)

func (s *Server) listChoferes(w http.ResponseWriter, r *http.Request) {
	choferes, err := s.store.ListChoferes()
	if err != nil {
		Error(w, http.StatusInternalServerError, "error al listar choferes")
		return
	}

	JSON(w, http.StatusOK, choferes)
}

func (s *Server) createChofer(w http.ResponseWriter, r *http.Request) {
	var chofer models.Chofer
	if err := DecodeJSON(r, &chofer); err != nil {
		Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	id, err := s.store.CreateChofer(chofer)
	if err != nil {
		Error(w, http.StatusInternalServerError, "error al crear chofer")
		return
	}

	chofer.ID = id
	JSON(w, http.StatusCreated, chofer)
}

func (s *Server) getChofer(w http.ResponseWriter, r *http.Request, id int) {
	chofer, err := s.store.GetChofer(id)
	if err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, "chofer no encontrado")
			return
		}
		Error(w, http.StatusInternalServerError, "error al obtener chofer")
		return
	}

	JSON(w, http.StatusOK, chofer)
}

func (s *Server) updateChofer(w http.ResponseWriter, r *http.Request, id int) {
	var chofer models.Chofer
	if err := DecodeJSON(r, &chofer); err != nil {
		Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	chofer.ID = id
	if err := s.store.UpdateChofer(id, chofer); err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, "chofer no encontrado")
			return
		}
		Error(w, http.StatusInternalServerError, "error al actualizar chofer")
		return
	}

	JSON(w, http.StatusOK, chofer)
}

func (s *Server) deleteChofer(w http.ResponseWriter, r *http.Request, id int) {
	if err := s.store.DeleteChofer(id); err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, "chofer no encontrado")
			return
		}
		Error(w, http.StatusInternalServerError, "error al eliminar chofer")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
