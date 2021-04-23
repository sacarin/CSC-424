package main

type budget struct {
	Amount   float64
	Category category
}

func (b *budget) insert(client_id int) error {
	exist, err := rowExist(`SELECT amount FROM budget WHERE client_id = $1
		AND cat_id = $2`, client_id, b.Category.ID)
	if exist {
		b.purge(client_id)
	}
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO budget (client_id, cat_id, amount)"+
		"VALUES ($1, $2, $3)", client_id, b.Category.ID, b.Amount)
	if err != nil {
		return err
	}

	return nil
}

func (b *budget) purge(client_id int) error {
	_, err := db.Exec("DELETE FROM budget WHERE client_id = $1 AND cat_id = $2",
		client_id, b.Category.ID)
	if err != nil {
		return err
	}
	return nil
}
