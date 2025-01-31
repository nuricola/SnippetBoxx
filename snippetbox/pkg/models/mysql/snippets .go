package mysql

import (
	"database/sql"
	"errors"
	"github.com/nuricola/snippetbox/pkg/models"
	
)

// SnippetModel - Определяем тип который обертывает пул подключения sql.DB
type SnippetModel struct {
	DB  *sql.DB
}


// Insert - Метод для создания новой заметки в базе дынных.
func (m *SnippetModel) Insert(title,content,expires string )(int, error){
	stmt := `INSERT INTO Snippets(title,content,created,expires)
			VALUES(?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIME(),INTERVAL ? DAY))`

	result , err := m.DB.Exec(stmt,title , content, expires)
	if err != nil{
		return 0,err
	}

	id , err := result.LastInsertId()
	if err != nil{
		return 0, err
	}

	return int(id),nil
	
}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func(m *SnippetModel) Get(id int)(*models.Snippet,error){
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt,id)

	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil{
		if errors.Is(err,sql.ErrNoRows){
			return nil, models.ErrNoRecord
		} else {
			return nil , err
		}
	}

	return s, nil

}

// Latest - Метод возвращает последние 10 заметок.
func (m *SnippetModel) Latest()([]*models.Snippet,error){
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	var snippets []*models.Snippet

	for rows.Next(){
		s:= &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
