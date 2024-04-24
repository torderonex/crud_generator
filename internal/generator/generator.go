package generator

import (
	"fmt"
	"strings"

	"github.com/torderonex/crudgenerator/internal/parser"
)

type CodeGenerator struct {
	Struct parser.StructWithFields
	Module string
}

func (g *CodeGenerator) generateRepository() string {
	code := fmt.Sprintf(`package postgres

import (
	"fmt"
	"%s"
	"github.com/jmoiron/sqlx"
)

type %s struct {
	Db *sqlx.DB
}

func New%s(db *sqlx.DB) *%s {
	return &%s{db}
}`, g.Module+"/"+g.Struct.Package, g.Struct.RepositoryName(), g.Struct.RepositoryName(), g.Struct.RepositoryName(), g.Struct.RepositoryName())

	return code
}

func (g *CodeGenerator) generateCreateFunction() string {
	query := fmt.Sprintf("\"INSERT INTO %s VALUES ($1,$2) RETURNING id\"", g.Struct.GetTableName())
	fn := fmt.Sprintf(`func (r %s) Create%s(%s) (int,error){
	var id int
	query := fmt.Sprintf(%s)
	row := r.Db.QueryRow(query, %s)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}`, g.Struct.RepositoryName(), g.Struct.Name, g.Struct.AsArgument(), query, g.Struct.FieldArgsToString())
	return fn
}

func (g *CodeGenerator) generateReadAllFunction() string {
	query := fmt.Sprintf("\"SELECT * from %s\"", g.Struct.GetTableName())
	fn := fmt.Sprintf(`func (r %s) GetAll%ss() ([]%s, error) {
	var c []%s
	query := %s
	err := r.Db.Select(&c, query)
	return c, err
}`, g.Struct.RepositoryName(), g.Struct.Name, g.Struct.NameWithPackage(), g.Struct.NameWithPackage(), query)
	return fn
}

func (g *CodeGenerator) generateReadOneFunction() string {
	query := fmt.Sprintf("\"SELECT * from %s WHERE id = $1\"", g.Struct.GetTableName())
	fn := fmt.Sprintf(`func (r %s) Get%sById(id int) (%s, error) {
	var c %s
	query := %s
	err := r.Db.Get(&c, query, id)
	return c, err
}`, g.Struct.RepositoryName(), g.Struct.Name, g.Struct.NameWithPackage(), g.Struct.NameWithPackage(), query)
	return fn
}

func (g *CodeGenerator) generateDeleteFunction() string {
	query := fmt.Sprintf("\"DELETE from %s WHERE id = $1\"", g.Struct.GetTableName())
	fn := fmt.Sprintf(`func (r %s) Delete%sById(id int) error{
	query := %s
	_, err := r.Db.Exec(query, id)
	return err
}`, g.Struct.RepositoryName(), g.Struct.Name, query)
	return fn
}

func (g *CodeGenerator) GenerateCRUD() string {
	funcs := []func() string{
		g.generateRepository,
		g.generateCreateFunction,
		g.generateReadAllFunction,
		g.generateReadOneFunction,
		g.generateDeleteFunction,
	}

	var res strings.Builder
	for _, f := range funcs {
		res.WriteString(f())
		res.WriteString("\n\n")
	}
	return res.String()
}
