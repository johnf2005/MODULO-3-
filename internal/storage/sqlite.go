package storage

import (
	"database/sql"
	"encoding/json"
	"os"

	"modulo3_go/internal/models"

	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dsn string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	if err := seedData(db); err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteStorage{db: db}, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS choferes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre TEXT NOT NULL,
			licencia TEXT NOT NULL,
			telefono TEXT NOT NULL,
			hora_entrada TEXT NOT NULL,
			hora_salida TEXT NOT NULL,
			estado TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS mantenimientos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			numero_carrito TEXT NOT NULL,
			fecha_mantenimiento TEXT NOT NULL,
			descripcion TEXT NOT NULL,
			estado_mantenimiento TEXT NOT NULL
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func seedData(db *sql.DB) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM choferes`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	content, err := os.ReadFile("sample_data.json")
	if err != nil {
		return err
	}

	type seedRecord struct {
		ID                  string `json:"id_chofer"`
		Nombre              string `json:"nombre_chofer"`
		Licencia            string `json:"licencia"`
		HoraEntrada         string `json:"hora_entrada"`
		HoraSalida          string `json:"hora_salida"`
		EstadoChofer        string `json:"estado_chofer"`
		IDMantenimiento     string `json:"id_mantenimiento"`
		FechaMantenimiento  string `json:"fecha_mantenimiento"`
		Descripcion         string `json:"descripcion"`
		EstadoMantenimiento string `json:"estado_mantenimiento"`
		NumeroCarrito       string `json:"numero_carrito"`
	}

	var records []seedRecord
	if err := json.Unmarshal(content, &records); err != nil {
		return err
	}

	for _, record := range records {
		if _, err := db.Exec(
			`INSERT OR REPLACE INTO choferes (id, nombre, licencia, telefono, hora_entrada, hora_salida, estado) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			record.ID, record.Nombre, record.Licencia, "", record.HoraEntrada, record.HoraSalida, record.EstadoChofer,
		); err != nil {
			return err
		}

		if _, err := db.Exec(
			`INSERT OR REPLACE INTO mantenimientos (id, numero_carrito, fecha_mantenimiento, descripcion, estado_mantenimiento) VALUES (?, ?, ?, ?, ?)`,
			record.IDMantenimiento, record.NumeroCarrito, record.FechaMantenimiento, record.Descripcion, record.EstadoMantenimiento,
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

func (s *SQLiteStorage) CreateChofer(chofer models.Chofer) (int, error) {
	result, err := s.db.Exec(
		`INSERT INTO choferes (nombre, licencia, telefono, hora_entrada, hora_salida, estado)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		chofer.Nombre, chofer.Licencia, chofer.Telefono, chofer.HoraEntrada, chofer.HoraSalida, chofer.Estado,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SQLiteStorage) ListChoferes() ([]models.Chofer, error) {
	rows, err := s.db.Query(`SELECT id, nombre, licencia, telefono, hora_entrada, hora_salida, estado FROM choferes`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Chofer
	for rows.Next() {
		var c models.Chofer
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Licencia, &c.Telefono, &c.HoraEntrada, &c.HoraSalida, &c.Estado); err != nil {
			return nil, err
		}
		list = append(list, c)
	}

	return list, rows.Err()
}

func (s *SQLiteStorage) GetChofer(id int) (models.Chofer, error) {
	var c models.Chofer
	row := s.db.QueryRow(`SELECT id, nombre, licencia, telefono, hora_entrada, hora_salida, estado FROM choferes WHERE id = ?`, id)
	if err := row.Scan(&c.ID, &c.Nombre, &c.Licencia, &c.Telefono, &c.HoraEntrada, &c.HoraSalida, &c.Estado); err != nil {
		return models.Chofer{}, err
	}
	return c, nil
}

func (s *SQLiteStorage) UpdateChofer(id int, chofer models.Chofer) error {
	result, err := s.db.Exec(
		`UPDATE choferes SET nombre = ?, licencia = ?, telefono = ?, hora_entrada = ?, hora_salida = ?, estado = ? WHERE id = ?`,
		chofer.Nombre, chofer.Licencia, chofer.Telefono, chofer.HoraEntrada, chofer.HoraSalida, chofer.Estado, id,
	)
	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *SQLiteStorage) DeleteChofer(id int) error {
	result, err := s.db.Exec(`DELETE FROM choferes WHERE id = ?`, id)
	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *SQLiteStorage) CreateMantenimiento(m models.Mantenimiento) (int, error) {
	result, err := s.db.Exec(
		`INSERT INTO mantenimientos (numero_carrito, fecha_mantenimiento, descripcion, estado_mantenimiento)
		 VALUES (?, ?, ?, ?)`,
		m.NumeroCarrito, m.FechaMantenimiento, m.Descripcion, m.EstadoMantenimiento,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SQLiteStorage) ListMantenimientos() ([]models.Mantenimiento, error) {
	rows, err := s.db.Query(`SELECT id, numero_carrito, fecha_mantenimiento, descripcion, estado_mantenimiento FROM mantenimientos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Mantenimiento
	for rows.Next() {
		var m models.Mantenimiento
		if err := rows.Scan(&m.ID, &m.NumeroCarrito, &m.FechaMantenimiento, &m.Descripcion, &m.EstadoMantenimiento); err != nil {
			return nil, err
		}
		list = append(list, m)
	}

	return list, rows.Err()
}

func (s *SQLiteStorage) GetMantenimiento(id int) (models.Mantenimiento, error) {
	var m models.Mantenimiento
	row := s.db.QueryRow(`SELECT id, numero_carrito, fecha_mantenimiento, descripcion, estado_mantenimiento FROM mantenimientos WHERE id = ?`, id)
	if err := row.Scan(&m.ID, &m.NumeroCarrito, &m.FechaMantenimiento, &m.Descripcion, &m.EstadoMantenimiento); err != nil {
		return models.Mantenimiento{}, err
	}
	return m, nil
}

func (s *SQLiteStorage) UpdateMantenimiento(id int, m models.Mantenimiento) error {
	result, err := s.db.Exec(
		`UPDATE mantenimientos SET numero_carrito = ?, fecha_mantenimiento = ?, descripcion = ?, estado_mantenimiento = ? WHERE id = ?`,
		m.NumeroCarrito, m.FechaMantenimiento, m.Descripcion, m.EstadoMantenimiento, id,
	)
	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *SQLiteStorage) DeleteMantenimiento(id int) error {
	result, err := s.db.Exec(`DELETE FROM mantenimientos WHERE id = ?`, id)
	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
