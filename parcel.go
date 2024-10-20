package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	result, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)", p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = ?", number)

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt) // Обратите внимание на '&'

	if err != nil {
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcels WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Parcel

	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcels SET ststus = ? WHERE number = ?", status, number)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	status := ParcelStatusRegistered
	_, err := s.db.Exec("UPDATE parcel SET address = ? WHERE number = ? AND status = ?", address, number, status)
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	status := ParcelStatusRegistered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = ? AND status = ?", number, status)
	if err != nil {
		return err
	}
	return nil
}
