package main


type Store interface {
	CreateBird(bird *Bird) error
	GetBirds() ([]*Bird, error)
}

func (store *dbStore) CreateBird(bird *Bird) error {
	_, err := store.db.Query("INSERT INTO birds(bird, description) VALUES ($1,$2)", bird.Species, bird.Descritpion)
	return err
}

func (store *dbStore) GetBirds() ([]*Bird, error) {
	rows, err := store.db.Query("SELECT bird, description from birds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	birds := []*Bird{}
	for rows.Next() {
		bird := &Bird{}
		if err := rows.Scan(&bird.Species, &bird.Descritpion); err != nil {
			return nil, err
		}
		birds = append(birds, bird)
	}
	return birds, nil
}

var store Store

func InitStore(s Store) {
	store = s
}
