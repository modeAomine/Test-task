package Service

import (
	"database/sql"
	"errors"
	"tests/DataBase"
	"tests/Model"
)

func GetWardrobeById(id int) (*Model.Wardrobe, error) {
	var wardrobe Model.Wardrobe
	err := DataBase.DB.QueryRow("SELECT id, title, quantity, price, description, height, width, depth, filename, link FROM Wardrobe WHERE id = $1", id).Scan(
		&wardrobe.ID,
		&wardrobe.Title,
		&wardrobe.Quantity,
		&wardrobe.Price,
		&wardrobe.Description,
		&wardrobe.Height,
		&wardrobe.Width,
		&wardrobe.Depth,
		&wardrobe.Filename,
		&wardrobe.Link)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &wardrobe, nil
}

func CreateWardrobe(w *Model.Wardrobe) error {
	requiredFields := map[string]interface{}{
		"title":       w.Title,
		"quantity":    w.Quantity,
		"price":       w.Price,
		"description": w.Description,
		"height":      w.Height,
		"width":       w.Width,
		"depth":       w.Depth,
		"filename":    w.Filename,
		"link":        w.Link,
	}

	for field, value := range requiredFields {
		if value == "" || value == 0 {
			return errors.New(field + " is required")
		}
	}

	query := `INSERT INTO Wardrobe (title, quantity, price, description, height, width, depth, filename, link) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	err := DataBase.DB.QueryRow(query, w.Title, w.Quantity, w.Price, w.Description, w.Height, w.Width, w.Depth, w.Filename, w.Link).Scan(&w.ID)
	if err != nil {
		return err
	}

	return nil
}

func UpdateWardrobe(w *Model.Wardrobe) error {
	storedWardrobe, err := GetWardrobeById(w.ID)
	if err != nil {
		return err
	}

	if storedWardrobe == nil {
		return errors.New("Wardrobe not found")
	}

	storedWardrobe.Title = w.Title
	storedWardrobe.Quantity = w.Quantity
	storedWardrobe.Price = w.Price
	storedWardrobe.Description = w.Description
	storedWardrobe.Height = w.Height
	storedWardrobe.Width = w.Width
	storedWardrobe.Depth = w.Depth
	storedWardrobe.Filename = w.Filename
	storedWardrobe.Link = w.Link

	_, err = DataBase.DB.Exec(`
		UPDATE Wardrobe 
		SET title = $1, quantity = $2, price = $3, description = $4, height = $5, width = $6, depth = $7, filename = $8, link = $9
		WHERE id = $10`,
		storedWardrobe.Title,
		storedWardrobe.Quantity,
		storedWardrobe.Price,
		storedWardrobe.Description,
		storedWardrobe.Height,
		storedWardrobe.Width,
		storedWardrobe.Depth,
		storedWardrobe.Filename,
		storedWardrobe.Link,
		storedWardrobe.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteWardrobe(id int) error {
	_, err := DataBase.DB.Exec("DELETE FROM Wardrobe WHERE id = $1", id)
	return err
}

func GetAllWardrobe() ([]Model.Wardrobe, error) {
	var wardrobes []Model.Wardrobe
	rows, err := DataBase.DB.Query("SELECT id, title, quantity, price, description, height, width, depth, filename, link FROM Wardrobe")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var wardrobe Model.Wardrobe
		err := rows.Scan(&wardrobe.ID, &wardrobe.Title, &wardrobe.Quantity, &wardrobe.Price, &wardrobe.Description, &wardrobe.Height, &wardrobe.Width, &wardrobe.Depth, &wardrobe.Filename, &wardrobe.Link)
		if err != nil {
			return nil, err
		}
		wardrobes = append(wardrobes, wardrobe)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return wardrobes, nil
}
