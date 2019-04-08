package profile

import "fmt"

func format(schema, table string) string {
	return fmt.Sprintf("%s.%s", schema, table)
}

var (
	profileTableSection = format(schemaAcl, "section")
	profileTableContent = format(schemaAcl, "content")
)
