package parser

import (
	"salami/compiler/types"
)

func (p *Parser) handleFieldLine() error {
	fieldNameToken := p.currentToken()
	p.advance()
	fieldValueToken := p.currentToken()
	p.advance()

	if p.currentToken().Type != types.Newline && p.currentToken().Type != types.EOF {
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
	case "Value":
		err := p.handleValueField(fieldNameToken, fieldValueToken)
		if err != nil {
			return err
		}
	default:
		return p.parseError(fieldNameToken, "unhandled error")
	}
	return nil
}

func (p *Parser) handleResourceTypeField(fieldNameToken *types.Token, fieldValueToken *types.Token) error {
	if p.currentObjectTypeIs(Unset) {
		p.setCurrentObjectType(Resource)
	}
	if !p.currentObjectTypeIs(Resource) {
		return p.parseError(fieldNameToken, "Resource type field can only be used on resource")
	}
	p.currentResource().ResourceType = types.ResourceType(fieldValueToken.Value)
	return nil
}

func (p *Parser) handleLogicalNameField(fieldNameToken *types.Token, fieldValueToken *types.Token) error {
	if p.currentObjectTypeIs(Unset) {
		return p.parseError(fieldNameToken, "Logical name field must be preceded by Resource type field")
	}
	if !p.currentObjectTypeIs(Resource) {
		return p.parseError(fieldNameToken, "Logical name field can only be used on resource")
	}
	p.currentResource().LogicalName = types.LogicalName(fieldValueToken.Value)
	return nil
}

func (p *Parser) handleDescriptionField(fieldNameToken *types.Token, fieldValueToken *types.Token) error {
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

func (p *Parser) handleNameField(fieldNameToken *types.Token, fieldValueToken *types.Token) error {
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

func (p *Parser) handleValueField(fieldNameToken *types.Token, fieldValueToken *types.Token) error {
	if p.currentObjectTypeIs(Unset) {
		return p.parseError(
			fieldNameToken,
			"ambiguous object type. Ensure object has Resource type field or @variable decorator before Value",
		)
	}
	if p.currentObjectTypeIs(Resource) {
		p.addFieldLineToNaturalLanguage(fieldNameToken, fieldValueToken)
	} else {
		p.currentVariable().Value = fieldValueToken.Value
	}
	return nil
}

func (p *Parser) addFieldLineToNaturalLanguage(fieldNameToken *types.Token, fieldValueToken *types.Token) {
	p.addLineToNaturalLanguage(fieldNameToken.Value + ": " + fieldValueToken.Value)
}
