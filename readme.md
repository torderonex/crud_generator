# Usage

- Download using `git clone https://github.com/torderonex/crud_generator`
- Compile the program with `go build -o crud_generate main.go`
- Add it to PATH
- Generate code with `crud_generate -i="path_to_go_file_with_entites" -o="path_to_output"`

# Example

\entites\user.go 
```go
package  entites

type  User  struct {
	Id  int
	Name  string
	Passowrd  string
}
```
```bash
crud_generate -i="\entites\user.go" -o="\repository"
```
\repository\user.go
```go
package postgres

import (
    "fmt"
    "crud_test/entites"
    "github.com/jmoiron/sqlx"
)

type UserRepository struct {
    Db * sqlx.DB
}

func NewUserRepository(db * sqlx.DB) * UserRepository {
    return &UserRepository {
        db
    }
}

func(r UserRepository) CreateUser(u entites.User)(int, error) {
    var id int
    query: = fmt.Sprintf("INSERT INTO users VALUES ($1,$2) RETURNING id")
    row: = r.Db.QueryRow(query, u.Id, u.Name, u.Passowrd)
    if err: = row.Scan( & id);
    err != nil {
        return 0, err
    }
    return id, nil
}

func(r UserRepository) GetAllUsers()([] entites.User, error) {
    var c[] entites.User
    query: = "SELECT * from users"
    err: = r.Db.Select( & c, query)
    return c, err
}

func(r UserRepository) GetUserById(id int)(entites.User, error) {
    var c entites.User
    query: = "SELECT * from users WHERE id = $1"
    err: = r.Db.Get( & c, query, id)
    return c, err
}

func(r UserRepository) DeleteUserById(id int) error {
    query: = "DELETE from users WHERE id = $1"
    _,
    err: = r.Db.Exec(query, id)
    return err
}
```
