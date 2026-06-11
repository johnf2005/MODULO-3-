package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"modulo3_go/internal/storage"
)

type Server struct {
	store storage.Storage
}

func NewServer(store storage.Storage) http.Handler {
	s := &Server{store: store}
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.rootHandler)
	mux.HandleFunc("/choferes", s.choferesHandler)
	mux.HandleFunc("/choferes/", s.choferHandler)
	mux.HandleFunc("/mantenimientos", s.mantenimientosHandler)
	mux.HandleFunc("/mantenimientos/", s.mantenimientoHandler)

	return mux
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	choferes, err := s.store.ListChoferes()
	if err != nil {
		Error(w, http.StatusInternalServerError, "error al obtener choferes")
		return
	}

	mantenimientos, err := s.store.ListMantenimientos()
	if err != nil {
		Error(w, http.StatusInternalServerError, "error al obtener mantenimientos")
		return
	}

	choferesOut := make([]map[string]any, 0, len(choferes))
	for _, c := range choferes {
		choferesOut = append(choferesOut, map[string]any{
			"id_chofer":     c.ID,
			"nombre_chofer": c.Nombre,
			"licencia":      c.Licencia,
			"hora_entrada":  c.HoraEntrada,
			"hora_salida":   c.HoraSalida,
			"estado_chofer": c.Estado,
		})
	}

	mantenimientosOut := make([]map[string]any, 0, len(mantenimientos))
	for _, m := range mantenimientos {
		mantenimientosOut = append(mantenimientosOut, map[string]any{
			"id_mantenimiento":     m.ID,
			"fecha_mantenimiento":  m.FechaMantenimiento,
			"descripcion":          m.Descripcion,
			"estado_mantenimiento": m.EstadoMantenimiento,
			"numero_carrito":       m.NumeroCarrito,
		})
	}

	JSON(w, http.StatusOK, map[string]any{
		"choferes":       choferesOut,
		"mantenimientos": mantenimientosOut,
	})
}

func (s *Server) choferesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listChoferes(w, r)
	case http.MethodPost:
		s.createChofer(w, r)
	default:
		Error(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

func (s *Server) choferHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/choferes/")
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getChofer(w, r, id)
	case http.MethodPut:
		s.updateChofer(w, r, id)
	case http.MethodDelete:
		s.deleteChofer(w, r, id)
	default:
		Error(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

func (s *Server) mantenimientosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listMantenimientos(w, r)
	case http.MethodPost:
		s.createMantenimiento(w, r)
	default:
		Error(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

func (s *Server) mantenimientoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/mantenimientos/")
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getMantenimiento(w, r, id)
	case http.MethodPut:
		s.updateMantenimiento(w, r, id)
	case http.MethodDelete:
		s.deleteMantenimiento(w, r, id)
	default:
		Error(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

func parseID(path, prefix string) (int, error) {
	if !strings.HasPrefix(path, prefix) {
		return 0, http.ErrNotSupported
	}

	value := strings.TrimPrefix(path, prefix)
	if value == "" {
		return 0, http.ErrMissingFile
	}

	return strconv.Atoi(value)
}
