package ctk

// TODO: support full path traversal

import (
	"bytes"
	"fmt"

	"github.com/tdewolff/parse/v2"
	tcss "github.com/tdewolff/parse/v2/css"

	"github.com/go-curses/cdk/lib/sync"
)

type cStyleSheet struct {
	Lexer      *tcss.Lexer
	Rules      []*StyleSheetRule
	MediaRules []*StyleSheetMedia

	sync.RWMutex
}

func newStyleSheet() *cStyleSheet {
	ss := &cStyleSheet{
		Lexer:      nil,
		Rules:      make([]*StyleSheetRule, 0),
		MediaRules: make([]*StyleSheetMedia, 0),
	}
	return ss
}

func newStyleSheetFromString(css string) (*cStyleSheet, error) {
	ss := newStyleSheet()
	if err := ss.ParseString(css); err != nil {
		return nil, err
	}
	return ss, nil
}

func (s *cStyleSheet) String() string {
	str := ""
	s.RLock()
	for _, r := range s.Rules {
		str += r.String() + "\n"
	}
	for _, m := range s.MediaRules {
		str += m.String() + "\n"
	}
	s.RUnlock()
	return str
}

func (s *cStyleSheet) ApplyStylesTo(w Widget) {
	selector := w.CssFullPath()
	styles := s.SelectProperties(selector)
	s.RLock()
	for s, _ := range styles {
		for k, v := range styles[s] {
			if err := w.SetCssPropertyFromStyle(k+":"+s, v.Value); err != nil {
				w.LogErr(err)
			}
		}
	}
	s.RUnlock()
}

func (s *cStyleSheet) SelectProperties(path string) (properties map[string]map[string]*StyleSheetProperty) {
	properties = make(map[string]map[string]*StyleSheetProperty)
	selector := newStyleSheetSelectorFromPath(path)
	s.RLock()
	for _, rule := range s.Rules {
		rSelect := rule.Selector
		grouped := parseStyleSheetSelectorGroup(rSelect)
		for _, gSelector := range grouped {
			if selector.Match(gSelector) {
				for _, elem := range rule.Properties {
					if _, ok := properties[gSelector.State]; !ok {
						properties[gSelector.State] = make(map[string]*StyleSheetProperty)
					}
					properties[gSelector.State][elem.Key] = elem
				}
			}
		}
	}
	s.RUnlock()
	return
}

func (s *cStyleSheet) ParseString(source string) (err error) {
	s.Lock()
	defer s.Unlock()
	s.Lexer = tcss.NewLexer(parse.NewInput(bytes.NewBufferString(source)))
	for {
		tt, data := s.Lexer.Next()
		switch tt {
		case tcss.ErrorToken:
			err = nil
			return
		case tcss.WhitespaceToken:
			continue // nop
		case tcss.CommentToken:
			continue // nop, ignore actual comments
		case tcss.AtKeywordToken:
			// data == "@media"
			if cssMediaRule, err := s.recurseMedia(); err != nil {
				return err
			} else {
				s.MediaRules = append(s.MediaRules, cssMediaRule)
			}
		default:
			// data == selector
			if cssRule, err := s.recurseRule(tt, data); err != nil {
				return err
			} else {
				s.Rules = append(s.Rules, cssRule)
			}
		}
	}
}

// consume up to (and including) the first opening bracket
func (s *cStyleSheet) recurseMedia() (mediaRule *StyleSheetMedia, err error) {
	mediaRule = &StyleSheetMedia{}
	var ruleMode bool
	for {
		tt, data := s.Lexer.Next()
		switch tt {
		case tcss.ErrorToken:
			return
		case tcss.LeftParenthesisToken:
			mediaRule.Conditions += "("
			continue // nop
		case tcss.RightParenthesisToken:
			mediaRule.Conditions += ")"
			continue // nop
		case tcss.WhitespaceToken:
			continue // nop
		case tcss.RightBraceToken:
			continue // nop
		case tcss.LeftBraceToken:
			ruleMode = true
			continue
		case tcss.LeftBracketToken, tcss.RightBracketToken, tcss.ColonToken, tcss.NumberToken, tcss.DimensionToken, tcss.DelimToken, tcss.IdentToken:
			if ruleMode {
				if cssRule, err := s.recurseRule(tt, data); err != nil {
					return nil, err
				} else {
					mediaRule.Rules = append(mediaRule.Rules, cssRule)
				}
			} else {
				mediaRule.Conditions += string(data)
			}
			continue // key / value parsing
		default:
			return nil, fmt.Errorf("recurseMediaRule: unexpected token type: %v (%v)", tt, data)
		}
	}
}

// consume up to (and including) the first opening curly brace
func (s *cStyleSheet) recurseRule(tt tcss.TokenType, data []byte) (cssRule *StyleSheetRule, err error) {
	cssRule = &StyleSheetRule{}
	for {
		switch tt {
		case tcss.ErrorToken:
			return
		case tcss.LeftBraceToken:
			if properties, err := s.recurseKeyValues(); err != nil {
				return nil, err
			} else {
				for _, v := range properties {
					cssRule.Properties = append(cssRule.Properties, v)
				}
			}
			return
		case tcss.RightBraceToken:
			return // end of rule block
		case tcss.WhitespaceToken:
			cssRule.Selector += " "
			tt, data = s.Lexer.Next()
			continue // nop
		case tcss.LeftBracketToken, tcss.RightBracketToken, tcss.CommentToken, tcss.DelimToken, tcss.HashToken, tcss.ColonToken, tcss.NumberToken, tcss.IdentToken:
			cssRule.Selector += string(data)
			tt, data = s.Lexer.Next()
			continue // key / value parsing
		case tcss.CommaToken:
			cssRule.Selector += ","
			tt, data = s.Lexer.Next()
			continue // key / value parsing
		default:
			return nil, fmt.Errorf("RecurseRule: unexpected token type: %v (%v)", tt, data)
		}
	}
}

// consume up to (and including) the first closing curly brace, returning the
// (unaltered) key/value pairs accumulated
func (s *cStyleSheet) recurseKeyValues() (properties map[string]*StyleSheetProperty, err error) {
	properties = make(map[string]*StyleSheetProperty)
	var key, value string
	var vType tcss.TokenType
	var isValue bool
	for {
		tt, data := s.Lexer.Next()
		switch tt {
		case tcss.LeftBraceToken:
			continue // ignore opening braces
		case tcss.RightBraceToken:
			return // closing brace completes
		case tcss.ColonToken:
			isValue = true
			continue // colons transition from key to value parsing
		case tcss.SemicolonToken:
			isValue = false
			properties[key] = &StyleSheetProperty{
				Key:   key,
				Value: value,
				Type:  vType,
			}
			key, value = "", ""
			vType = tcss.ErrorToken
			continue // semicolons transition from current pair to new pair
		case tcss.WhitespaceToken:
			continue // nop
		case tcss.LeftBracketToken, tcss.RightBracketToken, tcss.DelimToken, tcss.DimensionToken, tcss.HashToken, tcss.IdentToken:
			if isValue {
				vType = tt
				value += string(data)
			} else {
				key += string(data)
			}
			continue // key / value parsing
		default:
			return nil, fmt.Errorf("RecurseKeyValues: unexpected token type: %v (%v)", tt, data)
		}
	}
}
