package repositories

import (
	"api/src/models"
	"database/sql"
)

type Publishes struct {
	db *sql.DB
}

func NewPublishRepository(db *sql.DB) *Publishes {
	return &Publishes{db}
}

func (repository Publishes) Create(publish models.Publish) (uint64, error) {
	statement, erro := repository.db.Prepare(
		"insert into publishes (title, content, author_id) values (?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	result, erro := statement.Exec(publish.Title, publish.Content, publish.AuthorID)
	if erro != nil {
		return 0, erro
	}

	lastInsertId, erro := result.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(lastInsertId), nil
}

func (repository Publishes) GetByID(publishID uint64) (models.Publish, error) {
	rows, erro := repository.db.Query(`
	select p.*, u.nick from 
	publishes p inner join users u
	on u.id = p.author_id where p.id = ?`,
		publishID,
	)
	if erro != nil {
		return models.Publish{}, erro
	}
	defer rows.Close()

	var publish models.Publish

	if rows.Next() {
		if erro = rows.Scan(
			&publish.ID,
			&publish.Title,
			&publish.Content,
			&publish.AuthorID,
			&publish.Likes,
			&publish.Created_at,
			&publish.AuthorNick,
		); erro != nil {
			return models.Publish{}, erro
		}
	}

	return publish, nil
}

func (repository Publishes) Get(usuarioID uint64) ([]models.Publish, error) {
	rows, erro := repository.db.Query(`
	select distinct p.*, u.nick from publishes p 
	inner join users u on u.id = p.author_id 
	inner join followers s on p.author_id = s.user_id 
	where u.id = ? or s.follower_id = ?
	order by 1 desc`,
		usuarioID, usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var publicacoes []models.Publish

	for rows.Next() {
		var publish models.Publish

		if erro = rows.Scan(
			&publish.ID,
			&publish.Title,
			&publish.Content,
			&publish.AuthorID,
			&publish.Likes,
			&publish.Created_at,
			&publish.AuthorNick,
		); erro != nil {
			return nil, erro
		}

		publicacoes = append(publicacoes, publish)
	}

	return publicacoes, nil
}

func (repository Publishes) Update(publishID uint64, publish models.Publish) error {
	statement, erro := repository.db.Prepare("update publishes set title = ?, content = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publish.Title, publish.Content, publishID); erro != nil {
		return erro
	}

	return nil
}

func (repository Publishes) Delete(publishID uint64) error {
	statement, erro := repository.db.Prepare("delete from publishes where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publishID); erro != nil {
		return erro
	}

	return nil
}

func (repository Publishes) GetByUser(usuarioID uint64) ([]models.Publish, error) {
	rows, erro := repository.db.Query(`
		select p.*, u.nick from publishes p
		join users u on u.id = p.author_id
		where p.author_id = ?`,
		usuarioID,
	)
	if erro != nil {
		return nil, erro
	}
	defer rows.Close()

	var publicacoes []models.Publish

	for rows.Next() {
		var publish models.Publish

		if erro = rows.Scan(
			&publish.ID,
			&publish.Title,
			&publish.Content,
			&publish.AuthorID,
			&publish.Likes,
			&publish.Created_at,
			&publish.AuthorNick,
		); erro != nil {
			return nil, erro
		}

		publicacoes = append(publicacoes, publish)
	}

	return publicacoes, nil
}

func (repository Publishes) Like(publishID uint64) error {
	statement, erro := repository.db.Prepare("update publishes set likes = likes + 1 where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(publishID); erro != nil {
		return erro
	}

	return nil
}

func (repository Publishes) Deslike(publishID uint64) error {
	statement, erro := repository.db.Prepare(`
		update publishes set likes = 
		CASE 
			WHEN likes > 0 THEN likes - 1
			ELSE 0 
		END
		where id = ?
	`)
	if erro != nil {
		return erro
	}

	if _, erro = statement.Exec(publishID); erro != nil {
		return erro
	}

	return nil
}
