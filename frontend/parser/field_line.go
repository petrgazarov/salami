package parser

import (
	commonTypes "salami/common/types"
	frontendTypes "salami/frontend/types"
)

func (p *Parser) handleFieldLine() error {
	fieldNameToken := p.currentToken()
	p.advance()
	fieldValueToken := p.currentToken()
	p.advance()

	if p.currentToken().Type != frontendTypes.Newline && p.currentToken().Type != frontendTypes.EOF {
		return p.parseError(p.currentToken())
	}

	switch fieldNameToken.Value {
	case "Resource type":
		err := p.handleResourceTypeField(fieldNameToken, fieldValueToken)
		if err != nil {
			return err
		}
	case "Logical name":
		err := p.handleLogicalNameField(fieldNameToken, fieldValueToken)
		if err != nil {
			return err
		}
	case "Description":
		err := p.handleDescriptionField(fieldNameToken, fieldValueToken)
		if err != nil {
			return err
		}
	case "Name":
		err := p.handleNameField(fieldNameToken, fieldValueToken)
		if err != nil {
			return err
		}
	case "Default":
		err := p.handleDefaultField(fieldNameToken, fieldValueToken)
		if err != nil {
			return err
		}
	default:
		return p.parseError(fieldNameToken, "unhandled error")
	}
	return nil
}

func (p *Parser) handleResourceTypeField(
	fieldNameToken *frontendTypes.Token,
	fieldValueToken *frontendTypes.Token,
) error {
	if p.currentObjectTypeIs(Unset) {
		p.setCurrentObjectType(Resource)
	}
	if !p.currentObjectTypeIs(Resource) {
		return p.parseError(fieldNameToken, "Resource type field can only be used on resource")
	}
	p.currentResource().ResourceType = commonTypes.ResourceType(fieldValueToken.Value)
	return nil
}

func (p *Parser) handleLogicalNameField(
	fieldNameToken *frontendTypes.Token,
	fieldValueToken *frontendTypes.Token,
) error {
	if p.currentObjectTypeIs(Unset) {
		return p.parseError(fieldNameToken, "Logical name field must be preceded by Resource type field")
	}
	if !p.currentObjectTypeIs(Resource) {
		return p.parseError(fieldNameToken, "Logical name field can only be used on resource")
	}
	p.currentResource().LogicalName = commonTypes.LogicalName(fieldValueToken.Value)
	return nil
}

func (p *Parser) handleDescriptionField(
	fieldNameToken *frontendTypes.Token,
	fieldValueToken *frontendTypes.Token,
) error {
	if p.currentObjectTypeIs(Unset) {
		return p.parseError(
			fieldNameToken,
			"ambiguous object type. Ensure object has Resource type field or @variable decorator before Description",
		)
	}
	if p.currentObjectTypeIs(Resource) {
		p.addFieldLineToNaturalLanguage(fieldNameToken, fieldValueToken)
	} else {
		p.currentVariable().Description = fieldValueToken.Value
	}
	return nil
}

func (p *Parser) handleNameField(
	fieldNameToken *frontendTypes.Token,
	fieldValueToken *frontendTypes.Token,
) error {
	if p.currentObjectTypeIs(Unset) {
		return p.parseError(
			fieldNameToken,
			"ambiguous object type. Ensure object has Resource type field or @variable decorator before Name",
		)
	}
	if p.currentObjectTypeIs(Resource) {
		p.addFieldLineToNaturalLanguage(fieldNameToken, fieldValueToken)
	} else {
		p.currentVariable().Name = fieldValueToken.Value
	}
	return nil
}

func (p *Parser) handleDefaultField(
	fieldNameToken *frontendTypes.Token,
	fieldValueToken *frontendTypes.Token,
) error {
	if p.currentObjectTypeIs(Unset) {
		return p.parseError(
			fieldNameToken,
			"ambiguous object type. Ensure object has Resource type field or @variable decorator before Default",
		)
	}
	if p.currentObjectTypeIs(Resource) {
		p.addFieldLineToNaturalLanguage(fieldNameToken, fieldValueToken)
	} else {
		p.currentVariable().Default = fieldValueToken.Value
	}
	return nil
}

func (p *Parser) addFieldLineToNaturalLanguage(
	fieldNameToken *frontendTypes.Token,
	fieldValueToken *frontendTypes.Token,
) {
	p.addLineToNaturalLanguage(fieldNameToken.Value + ": " + fieldValueToken.Value)
}
