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
	ip, err := s.db.Exec("insert into parcel (client, status, address, created_at) values (?, ?, ?, ?)", p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	currentNumberParcel, _ := ip.LastInsertId()
	return int(currentNumberParcel), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	err := s.db.QueryRow("select number, client, status, address, created_at from parcel where number = ?", number).Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("select number, client, status, address, created_at from parcel where client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("update parcel set status = ? where number = ?", status, number)
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("update parcel set address = ? where number = ? and status = ?", address, number, ParcelStatusRegistered)
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("delete from parcel where number = ? and status = ?", number, ParcelStatusRegistered)
	if err != nil {
		return err
	}
	return nil
}
