package templates

import (
	"bytes"
	"embed"
	"io/fs"
	"strings"
	"text/template"
)

const SystemTemplatePath = "system.tmpl"
const UpdateResourceTemplatePath = "update/resource.tmpl"
const UpdateVariableTemplatePath = "update/variable.tmpl"
const CreateResourceTemplatePath = "create/resource.tmpl"
const CreateVariableTemplatePath = "create/variable.tmpl"
const FixResourceValidationIssueTemplatePath = "fix_validation_issue/resource.tmpl"
const FixVariableValidationIssueTemplatePath = "fix_validation_issue/variable.tmpl"

func PopulateTemplate(templatePath string, dataStruct interface{}) (string, error) {
	templateString, err := ReadTemplateFile(templatePath)
	if err != nil {
		return "", err
	}

	result, err := replaceTemplateVariables(templateString, dataStruct)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func replaceTemplateVariables(templateString string, dataStruct interface{}) (string, error) {
	tmpl, err := template.New("template").Parse(templateString)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, dataStruct)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

//go:embed create update fix_validation_issue system.tmpl
var templatesFS embed.FS

func ReadTemplateFile(filePath string) (string, error) {
	data, err := fs.ReadFile(templatesFS, filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
