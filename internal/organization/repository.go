package organization

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrgNode struct {
	ID           string     `json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	PositionName *string    `json:"position_name,omitempty"`
	DepartmentName *string  `json:"department_name,omitempty"`
	Email        *string    `json:"email,omitempty"`
	PhotoURL     *string    `json:"photo_url,omitempty"`
	Children     []*OrgNode `json:"children"`
}

type OrgRepository struct {
	pool *pgxpool.Pool
}

func NewOrgRepository(pool *pgxpool.Pool) *OrgRepository {
	return &OrgRepository{pool: pool}
}

func (r *OrgRepository) GetTree(ctx context.Context, companyID string) ([]*OrgNode, error) {
	query := `
		SELECT
			e.id, e.first_name, e.last_name,
			p.name, d.name, e.email, e.photo_url, e.manager_id
		FROM employees e
		LEFT JOIN positions p ON p.id = e.position_id
		LEFT JOIN departments d ON d.id = e.department_id
		WHERE e.company_id = $1 AND e.status = 'active'
		ORDER BY e.last_name, e.first_name`

	rows, err := r.pool.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type rawEmployee struct {
		ID             string
		FirstName      string
		LastName       string
		PositionName   *string
		DepartmentName *string
		Email          *string
		PhotoURL       *string
		ManagerID      *string
	}

	var all []rawEmployee
	for rows.Next() {
		var e rawEmployee
		if err := rows.Scan(
			&e.ID, &e.FirstName, &e.LastName,
			&e.PositionName, &e.DepartmentName, &e.Email, &e.PhotoURL, &e.ManagerID,
		); err != nil {
			return nil, err
		}
		all = append(all, e)
	}

	nodeMap := make(map[string]*OrgNode)
	for _, e := range all {
		nodeMap[e.ID] = &OrgNode{
			ID:             e.ID,
			FirstName:      e.FirstName,
			LastName:       e.LastName,
			PositionName:   e.PositionName,
			DepartmentName: e.DepartmentName,
			Email:          e.Email,
			PhotoURL:       e.PhotoURL,
			Children:       []*OrgNode{},
		}
	}

	var roots []*OrgNode
	for _, e := range all {
		node := nodeMap[e.ID]
		if e.ManagerID != nil {
			if parent, ok := nodeMap[*e.ManagerID]; ok {
				parent.Children = append(parent.Children, node)
			} else {
				roots = append(roots, node)
			}
		} else {
			roots = append(roots, node)
		}
	}

	return roots, nil
}

func (r *OrgRepository) HasCycle(ctx context.Context, employeeID, newManagerID string) (bool, error) {
	if employeeID == newManagerID {
		return true, nil
	}

	current := newManagerID
	visited := map[string]bool{employeeID: true}

	for {
		visited[current] = true

		var managerID *string
		err := r.pool.QueryRow(ctx,
			`SELECT manager_id FROM employees WHERE id = $1`, current,
		).Scan(&managerID)
		if err != nil {
			return false, nil
		}

		if managerID == nil {
			return false, nil
		}

		if *managerID == employeeID {
			return true, nil
		}

		if visited[*managerID] {
			return false, nil
		}

		current = *managerID
	}
}
