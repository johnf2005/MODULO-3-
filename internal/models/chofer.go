package models

type Chofer struct {
	ID          int    `json:"id_chofer,string"`
	Nombre      string `json:"nombre_chofer"`
	Licencia    string `json:"licencia"`
	Telefono    string `json:"telefono"`
	HoraEntrada string `json:"hora_entrada"`
	HoraSalida  string `json:"hora_salida"`
	Estado      string `json:"estado_chofer"`
}
