package di

import "github.com/make-software/ces-go-parser"

type CESParserAware struct {
	parser *ces.EventParser
}

func (s *CESParserAware) SetCESParser(parser *ces.EventParser) {
	s.parser = parser
}

func (s *CESParserAware) GetCESParser() *ces.EventParser {
	return s.parser
}
