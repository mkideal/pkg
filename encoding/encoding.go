package encoding

import (
	"bytes"
	"text/scanner"
)

type Node interface {
	Pos() scanner.Position
}

type Comment struct {
	Slash scanner.Position
	Text  string
}

func (c *Comment) Pos() scanner.Position { return c.Slash }

type CommentGroup struct {
	List []*Comment
}

func (g *CommentGroup) Pos() scanner.Position { return g.List[0].Pos() }

func (g *CommentGroup) Text() string {
	if g == nil || len(g.List) == 0 {
		return ""
	}
	var buf bytes.Buffer
	for _, c := range g.List {
		if buf.Len() > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(c.Text)
	}
	return buf.String()
}

type Parser struct {
	scanner *scanner.Scanner
	mode    uint

	Pos scanner.Position
	Tok rune
	Lit string

	Comments    []*CommentGroup
	LeadComment *CommentGroup
	LineComment *CommentGroup
}

func (p *Parser) Init(s *scanner.Scanner) {
	p.scanner = s
}

func (p *Parser) Next() {
	p.LeadComment = nil
	p.LineComment = nil
	prev := p.Pos
	p.next0()

	if p.Tok == scanner.Comment {
		var comment *CommentGroup
		var endline int

		if p.Pos.Line == prev.Line {
			comment, endline = p.consumeCommentGroup(0)
			if p.Pos.Line != endline {
				p.LineComment = comment
			}
		}

		endline = -1
		for p.Tok == scanner.Comment {
			comment, endline = p.consumeCommentGroup(1)
		}

		if endline+1 == p.Pos.Line {
			p.LeadComment = comment
		}
	}
}

func (p *Parser) next0() {
	p.Tok = p.scanner.Scan()
	p.Pos = p.scanner.Pos()
	p.Lit = p.scanner.TokenText()
}

func (p *Parser) consumeComment() (comment *Comment, endline int) {
	endline = p.Pos.Line
	if len(p.Lit) > 0 && p.Lit[1] == '*' {
		for i := 0; i < len(p.Lit); i++ {
			if p.Lit[i] == '\n' {
				endline++
			}
		}
	}

	comment = &Comment{Slash: p.Pos, Text: p.Lit}
	p.next0()

	return
}

func (p *Parser) consumeCommentGroup(n int) (comments *CommentGroup, endline int) {
	var list []*Comment
	endline = p.Pos.Line
	for p.Tok == scanner.Comment && p.Pos.Line <= endline+n {
		var comment *Comment
		comment, endline = p.consumeComment()
		list = append(list, comment)
	}

	// add comment group to the comments list
	comments = &CommentGroup{List: list}
	p.Comments = append(p.Comments, comments)

	return
}
