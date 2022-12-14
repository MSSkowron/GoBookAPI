package storage

import (
	"database/sql"
	"log"

	"github.com/MSSkowron/GoBankAPI/model"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*model.User) error
	GetUserByEmail(string) (*model.User, error)
	GetBooks() ([]*model.Book, error)
	CreateBook(*model.Book) error
	GetBookByID(int) (*model.Book, error)
	DeleteBookByID(int) error
	UpdateBook(*model.Book) error
}

type PostgresSQLStorage struct {
	db *sql.DB
}

func NewPostgresSQLStorage() (*PostgresSQLStorage, error) {
	connStr := "user=gobookapiuser dbname=postgres password=gobookapipassword sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	postgresSQLStore := &PostgresSQLStorage{
		db: db,
	}

	if err := postgresSQLStore.init(); err != nil {
		return nil, err
	}

	return postgresSQLStore, nil
}

func (s *PostgresSQLStorage) init() error {
	if err := s.createUsersTable(); err != nil {
		return err
	}

	if err := s.createBooksTable(); err != nil {
		return err
	}

	return nil
}

func (s *PostgresSQLStorage) createUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id INT GENERATED ALWAYS AS IDENTITY,
		email varchar(50) NOT NULL,
		password varchar(256) NOT NULL,
		first_name varchar(50) NOT NULL,
		last_name varchar(50) NOT NULL,
		age smallint NOT NULL,
		PRIMARY KEY(id)
	);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresSQLStorage) createBooksTable() error {
	query := `CREATE TABLE IF NOT EXISTS books (
		id INT GENERATED ALWAYS AS IDENTITY,
		author  varchar(100) NOT NULL,
		title varchar(100) NOT NULL, 
		PRIMARY KEY(id)
	);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresSQLStorage) CreateUser(user *model.User) error {
	query := `insert into "users" (email, password, first_name, last_name, age) values ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(query, user.Email, user.Password, user.FirstName, user.LastName, user.Age)
	if err != nil {
		log.Println("[PostgresSQLStorage] Error while inserting new user: " + err.Error())
		return err
	}

	log.Println("[PostgresSQLStorage] Inserted new user")

	return nil
}

func (s *PostgresSQLStorage) GetUserByEmail(email string) (*model.User, error) {
	query := `select * from users where email=$1`

	row := s.db.QueryRow(query, email)

	user := &model.User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age); err != nil {
		return nil, err
	}

	log.Println("[PostgresSQLStorage] User correctly pulled from database")

	return user, nil
}

func (s *PostgresSQLStorage) CreateBook(book *model.Book) error {
	query := `insert into "books" (author, title) values ($1, $2)`

	_, err := s.db.Exec(query, book.Author, book.Title)
	if err != nil {
		log.Println("[PostgresSQLStorage] Error while inserting new book: " + err.Error())
		return err
	}

	log.Println("[PostgresSQLStorage] Inserted new book")

	return nil
}

func (s *PostgresSQLStorage) GetBooks() ([]*model.Book, error) {
	query := `select * from books`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	books := []*model.Book{}
	for rows.Next() {
		book := &model.Book{}
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	log.Println("[PostgresSQLStorage] Books correctly pulled from database")

	return books, nil
}

func (s *PostgresSQLStorage) GetBookByID(id int) (*model.Book, error) {
	query := `select * from books where id=$1`

	row := s.db.QueryRow(query, id)

	book := &model.Book{}
	if err := row.Scan(&book.ID, &book.Title, &book.Author); err != nil {
		return nil, err
	}

	log.Println("[PostgresSQLStorage] Book correctly pulled from database")

	return book, nil
}

func (s *PostgresSQLStorage) DeleteBookByID(id int) error {
	query := `delete from books where id=$1`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	log.Println("[PostgresSQLStorage] Book correctly pulled from database")

	return nil
}

func (s *PostgresSQLStorage) UpdateBook(book *model.Book) error {
	query := `UPDATE books SET author = $1, title= $2 WHERE id = $3;`

	_, err := s.db.Exec(query, book.Author, book.Title, book.ID)
	if err != nil {
		return err
	}

	log.Println("[PostgresSQLStorage] Book correctly updated")

	return nil
}
