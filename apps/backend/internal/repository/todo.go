package repository

import (
	"context"
	"fmt"

	"github.com/aneesh1213/tasker/internal/model/todo"
	"github.com/aneesh1213/tasker/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TodoRepository struct {
	server *server.Server
}

func NewTodoRepository(server *server.Server) *TodoRepository {
	return &TodoRepository{server: server}
}

func (r *TodoRepository) CreateTodo(ctx context.Context, userID string, payload *todo.CreateTodoPayload) (*todo.Todo, error) {
	stmt := `
		INSERT INTO
			todos (
				user_id,
				title,
				description,
				priority,
				due_date,
				parent_todo_id,
				category_id,
				metadata
			)
		VALUES
			(
				@user_id,
				@title,
				@description,
				@priority,
				@due_date,
				@parent_todo_id,
				@category_id,
				@metadata
			)
		RETURNING
		*
	`

	priority := todo.PriorityMedium
	if payload.Priority != nil {
		priority = *payload.Priority
	}

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":        userID,
		"title":          payload.Title,
		"description":    payload.Description,
		"priority":       priority,
		"due_date":       payload.DueDate,
		"parent_todo_id": payload.ParentTodoID,
		"category_id":    payload.CategoryID,
		"metadata":       payload.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create todo query for user_id=%s title=%s: %w", userID, payload.Title, err)
	}

	todoItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[todo.Todo])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:todos for user_id=%s title=%s: %w", userID, payload.Title, err)
	}

	return &todoItem, nil
}

func (r *TodoRepository) GetTodoByID(ctx context.Context, userID string, todoID uuid.UUID) {
	stmt := `
		SELECT
			t.*,
			CASE
				WHEN c.id IS NOT NULL THEN to_jsonb(camel (c))
				ELSE NULL
			END AS category,
			COALESCE(
				jsonb_agg(
					to_jsonb(camel (child))
					ORDER BY
						child.sort_order ASC,
						child.created_at ASC
				) FILTER (
					WHERE
						child.id IS NOT NULL
				),
				'[]'::JSONB
			) AS children,
			COALESCE(
				jsonb_agg(
					to_jsonb(camel (com))
					ORDER BY
						com.created_at ASC
				) FILTER (
					WHERE
						com.id IS NOT NULL
				),
				'[]'::JSONB
			) AS comments,
			COALESCE(
					jsonb_agg(
						to_jsonb(camel (att))
						ORDER BY
							att.created_at DESC
					) FILTER (
						WHERE
							att.id IS NOT NULL
					),
					'[]'::JSONB
				) AS attachments
		FROM
			todos t
			LEFT JOIN todo_categories c ON c.id=t.category_id
			AND c.user_id=@user_id
			LEFT JOIN todos child ON child.parent_todo_id=t.id
			AND child.user_id=@user_id
			LEFT JOIN todo_comments com ON com.todo_id=t.id
			AND com.user_id=@user_id
			LEFT JOIN todo_attachments att ON att.todo_id=t.id
		WHERE
			t.id=@id
			AND t.user_id=@user_id
		GROUP BY
			t.id,
			c.id
	`
}
