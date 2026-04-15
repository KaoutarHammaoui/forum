package models

import "forum/database"

type Category struct {
	IdCat int
	Name  string
}

func GetAllCategory() ([]Category, error) {
	categories := []Category{}
	query := "SELECT id, name FROM category"
	lignes, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer lignes.Close()

	for lignes.Next() {
		category := Category{}
		err := lignes.Scan(&category.IdCat, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := lignes.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func GetCategoryByName(name string) (Category, error) {
	category := Category{}
	query := "SELECT id FROM category WHERE name = ?"
	row := database.DB.QueryRow(query, name)
	err := row.Scan(&category.IdCat)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func InsertPostCategory(postID int64, categoryID int) error {
	_, err := database.DB.Exec("INSERT INTO post_category (post_id, category_id) VALUES (?, ?)", postID, categoryID)
	return err
}
