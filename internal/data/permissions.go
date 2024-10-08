package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

type PermissionDatabaseModel struct {
	DB *sql.DB
}

func (m *PermissionDatabaseModel) GetAllForUser(id int64) (Permissions, error) {
	query := `
		SELECT permissions.code
		FROM permissions
		INNER JOIN users_permissions
		ON users_permissions.permission_id = permissions.id
		WHERE users_permissions.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions
	for rows.Next() {
		var code string
		err := rows.Scan(&code)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, code)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (m *PermissionDatabaseModel) AddForUser(userID int64, codes ...string) error {
	// the 2nd line of the query creates a sub-table created from the user ID as the first field
	// and the corresponding ids for the permission codes in the codes array.
	query := `
		INSERT INTO users_permissions
		SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID, pq.Array(codes))
	return err
}
