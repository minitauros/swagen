package swagger

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"text/template"

	sql2 "github.com/minitauros/swagen/sql"
	"github.com/pkg/errors"
)

type (
	// ServiceInfo describes the service for which we are generating Swagger YML.
	ServiceInfo struct {
		Name string
		Host string
	}

	// Generator is the service that takes care of the Swager YML generation.
	Generator struct {
		TableService sql2.TableService

		// The following values will be written by parsing a config file into this struct.
		ServiceInfo ServiceInfo
		Resources   map[string]Resource // [table name]resource: The resources that we are generating Swagger YML for.
	}

	// TemplateData contains the data that is used when rendering the full Swagger YML.
	TemplateData struct {
		ServiceInfo ServiceInfo
		Resources   string
		Definitions string
	}

	// Resource describes a resource which we grant API access to.
	Resource struct {
		Title      string
		Path       string
		Definition Definition
		Get        Request
	}

	// Request describes the a request with which the API grants access to a resource.
	Request struct {
		Params []RequestParam
	}

	// RequestParam is a parameter that can be used in a request.
	RequestParam struct {
		In          string
		Name        string
		Required    bool
		Type        FieldType
		Description string
	}

	// Definition is the definition of a resource.
	Definition struct {
		Name   string
		Fields []DefinitionField
	}

	// DefinitionField is a field of a definition.
	DefinitionField struct {
		Name       string
		Type       FieldType
		IsNullable bool
	}

	// FieldType describes the type of a (definition) field.
	FieldType struct {
		Name            string
		ExtraProperties map[string]string
	}
)

// getTypeFromSqlColumnType returns the field type for the given SQL column.
func getTypeFromSqlColumnType(col *sql.ColumnType) (FieldType, error) {
	switch col.DatabaseTypeName() {
	case "VARCHAR":
		return FieldType{"string", nil}, nil
	case "TEXT":
		return FieldType{"string", nil}, nil
	case "CHAR":
		return FieldType{"string", nil}, nil
	case "DATETIME":
		return FieldType{"string", map[string]string{
			"format": "date-time",
		}}, nil
	case "TINYINT":
		return FieldType{"boolean", map[string]string{
			"format": "int64",
		}}, nil
	case "SMALLINT":
		return FieldType{"integer", map[string]string{
			"format": "int64",
		}}, nil
	case "INT":
		return FieldType{"integer", map[string]string{
			"format": "int64",
		}}, nil
	case "BIGINT":
		return FieldType{"integer", map[string]string{
			"format": "int64",
		}}, nil
	}
	return FieldType{}, errors.WithStack(fmt.Errorf("type " + col.DatabaseTypeName() + " not supported"))
}

// Generate generates the Swagger YML, which is returned as a string.
func (g *Generator) Generate() (string, error) {
	resources, err := g.generateResources()
	if err != nil {
		return "", err
	}

	fullTempl, err := template.New("resources").Parse(fullTemplate)
	if err != nil {
		return "", errors.WithStack(err)
	}

	resourceTempl, err := template.New("resources").Parse(resourceTemplate)
	if err != nil {
		return "", errors.WithStack(err)
	}

	definitionTempl, err := template.New("definitions").Parse(definitionTemplate)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// Write resources part of template.
	resourcesBuf := bytes.Buffer{}
	for _, resource := range resources {
		err = resourceTempl.Execute(&resourcesBuf, resource)
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	// Write definitions part of template.
	definitionsBuf := bytes.Buffer{}
	for _, resource := range resources {
		err = definitionTempl.Execute(&definitionsBuf, resource.Definition)
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	// Put everything together in the full template.
	data := TemplateData{
		ServiceInfo: g.ServiceInfo,
		Resources:   resourcesBuf.String(),
		Definitions: definitionsBuf.String(),
	}
	fullBuf := bytes.Buffer{}
	err = fullTempl.Execute(&fullBuf, data)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return fullBuf.String(), err
}

// generateResources uses the resources set on the generator (see field documentation) to generate resource objects.
func (g *Generator) generateResources() ([]Resource, error) {
	var resources []Resource

	for tableName, resource := range g.Resources {
		cols, err := g.TableService.GetColumns(tableName)
		if err != nil {
			return nil, err
		}

		var fields []DefinitionField
		for _, col := range cols {
			fieldType, err := getTypeFromSqlColumnType(col)
			if err != nil {
				return nil, err
			}

			isNullable, ok := col.Nullable()
			if !ok {
				return nil, errors.New("not ok")
			}

			fields = append(fields, DefinitionField{
				Name:       col.Name(),
				Type:       fieldType,
				IsNullable: isNullable,
			})
		}

		resource.Path = strings.ReplaceAll(tableName, "_", "-")
		resource.Definition.Fields = fields

		resources = append(resources, resource)
	}

	return resources, nil
}
