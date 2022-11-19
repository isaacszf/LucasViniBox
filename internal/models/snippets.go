package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet type to hold the data for an individual snippet. The struct fields
// correspond to the MySQL snippets table fields
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// / SnippetModel wraps a sql.DB connection pool
type SnippetModel struct{ Database *sql.DB }

// This function will insert a new snippet into the database
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	// SQL statement that we want to execute
	// The "?" acts as a placeholder
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Executing the statement using hhe "Exec" function. This function will return a
	// sql.Result, which contains some basic information about what happened when
	// the statement was executed
	result, err := m.Database.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Getting the ID field of our newly inserted record in the snippets table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// This function will return a snippet based on its ID
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires >
	UTC_TIMESTAMP() AND id = ?`

	// Using "QueryRow" function on the connection pool to execute our SQL statement.
	// This returns a pointer to a sql.Row object which holds the result  from the db
	row := m.Database.QueryRow(stmt, id)
	s := &Snippet{} // Initializing a pointer to a new zeroed Snippet struct

	// Using "row.Scan" to copy the values from each filed in sql.Row to the corresponding
	// field in the Snippet struct
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// This function will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires >
	UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Using "Query" method on the connection pool to execute our SQL statement.
	// This returns a sql.Rows containing the result of out query
	rows, err := m.Database.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Ensuring that sql.Rows resultset is always properly closed
	// before the "Latest() method returns"
	defer rows.Close()

	snippets := []*Snippet{}

	// Using "row.Next" to iterate through the rows in the resultset. This prepares the
	// first (and ten each subsequent) row to be acted on by the "row.Scan" method
	for rows.Next() {
		s := &Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// Append it to the slice of snippets
		snippets = append(snippets, s)
	}

	// When the "rows.Next" finishes, we call "rows.Err" to retrieve any error that
	// was encountered during the iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
