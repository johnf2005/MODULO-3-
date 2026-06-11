package models

type Mantenimiento struct {
	ID                  int    `json:"id_mantenimiento,string"`
	FechaMantenimiento  string `json:"fecha_mantenimiento"`
	Descripcion         string `json:"descripcion"`
	EstadoMantenimiento string `json:"estado_mantenimiento"`
	NumeroCarrito       string `json:"numero_carrito"`
}
