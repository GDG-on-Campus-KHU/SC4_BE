package services

import (
	"database/sql"
	"errors"

	"github.com/GDG-on-Campus-KHU/SC4_BE/db"
)

var ErrNoExistingSupplies = errors.New("no existing supplies found")

type SuppliesService struct{}

func NewSuppliesService() *SuppliesService {
	return &SuppliesService{}
}

func (s *SuppliesService) GetUserSupplies(userID int) (map[string]bool, error) {
	supplies := make(map[string]bool)

	query := `
        SELECT s.name, CASE 
            WHEN uc.is_checked IS NULL THEN false
            WHEN uc.is_checked = 1 THEN true
            ELSE false 
        END as status
        FROM supplies s
        LEFT JOIN user_checklist uc ON s.id = uc.supply_id AND uc.user_id = ?
        ORDER BY s.id
    `

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var status bool
		if err := rows.Scan(&name, &status); err != nil {
			return nil, err
		}
		supplies[name] = status
	}

	return supplies, nil
}

func (s *SuppliesService) SaveUserSupplies(userID int, supplies map[string]bool) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM user_checklist WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	for name, isChecked := range supplies {
		var supplyID int
		err := tx.QueryRow("SELECT id FROM supplies WHERE name = ?", name).Scan(&supplyID)
		if err != nil {
			continue
		}

		_, err = tx.Exec(
			"INSERT INTO user_checklist (user_id, supply_id, is_checked) VALUES (?, ?, ?)",
			userID, supplyID, isChecked,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SuppliesService) UpdateUserSupplies(userID int, supplies map[string]bool) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM user_checklist WHERE user_id = ?)",
		userID,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return ErrNoExistingSupplies
	}

	for name, isChecked := range supplies {
		var supplyID int
		err := tx.QueryRow("SELECT id FROM supplies WHERE name = ?", name).Scan(&supplyID)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return err
		}

		var itemExists bool
		err = tx.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM user_checklist WHERE user_id = ? AND supply_id = ?)",
			userID, supplyID,
		).Scan(&itemExists)
		if err != nil {
			return err
		}

		if itemExists {
			_, err = tx.Exec(
				"UPDATE user_checklist SET is_checked = ? WHERE user_id = ? AND supply_id = ?",
				isChecked, userID, supplyID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
