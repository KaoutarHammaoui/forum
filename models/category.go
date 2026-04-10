package models

import "forum/database"

type Category struct {
	IdCat int
	Name  string
}

func GetAllCategories() ([]Category, error) {
	categories := []Category{}

	rows, err := database.DB.Query("SELECT id, name FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		category := Category{}
		err := rows.Scan(&category.IdCat, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func InsertPostCategory(postID, categoryID int) error {
	_, err := database.DB.Exec("INSERT INTO post_category (post_id, category_id) VALUES (?, ?)", postID, categoryID)
	return err
}
