package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDropTable(t *testing.T) {
	a := assert.New(t)
	table1, _ := NewTable(
		"TABLE_A",
		IntColumn("id", CO_PrimaryKey),
		IntColumn("test1"),
		IntColumn("test2"),
	)

	query, args, err := DropTable(table1).ToSql()
	a.Equal(`DROP TABLE "TABLE_A";`, query)
	a.Equal([]interface{}{}, args)
	a.Nil(err)
}
