package sqlbuilder

// Buildable is an interface for performing the package-level functions
// which all require a package-level Dialect to be configured using a
// call to SetDialect. Buildable allows for using one or more independent
// Dialects at the same time within the same codebase.
type Buildable interface {
	AlterTable(tbl Table) AlterTableBuilder
	CreateTable(tbl Table) CreateTableBuilder
	CreateIndex(tbl Table) CreateIndexBuilder
	Delete(from Table) DeleteBuilder
	Insert(into Table) InsertBuilder
	Select(from Table) SelectBuilder
}

type buildable struct {
	dialect Dialect
}

// NewBuildable constructs a new Buildable instance, configured with the given
// Dialect
func NewBuildable(d Dialect) Buildable {
	b := &buildable{dialect: d}
	return b
}

func (b *buildable) AlterTable(tbl Table) AlterTableBuilder {
	return alterTable(tbl, b.dialect)
}

func (b *buildable) CreateTable(tbl Table) CreateTableBuilder {
	return createTable(tbl, b.dialect)
}

func (b *buildable) CreateIndex(tbl Table) CreateIndexBuilder {
	return createIndex(tbl, b.dialect)
}

func (b *buildable) Delete(from Table) DeleteBuilder {
	return deleteFn(from, b.dialect)
}

func (b *buildable) Insert(into Table) InsertBuilder {
	return insert(into, b.dialect)
}

func (b *buildable) Select(from Table) SelectBuilder {
	return selectFn(from, b.dialect)
}
