package db

// ReservationsTable is a CRUD interface for accessing the "reservation" table
type ReservationsTable struct {
	RethinkCRUD
}

func (r *ReservationsTable) IsUsernameUsed(name string) (bool, error) {
	count, err := r.FindByAndCount("name", name)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (r *ReservationsTable) IsEmailUsed(email string) (bool, error) {
	count, err := r.FindByAndCount("email", email)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
