package parser

import (
	"html_parser/internal/client"
	"strings"
)

type Parser struct {
	Tree *ParsingTree
	html string
}

type ParsingTree struct {
	Root         *Node
	ElemsByTag   map[string][]*Node
	ElemsByClass map[string][]*Node
	ElemsById    map[string][]*Node
}

type Node struct {
	Name            string
	Class           string
	Id              string
	Text            string
	AllChildrenText string
	Attributes      map[string]string
	Children        []*Node
}

func NewParser(html string) *Parser {
	return &Parser{html: html}
}

func NewParserByUrl(url string) *Parser {
	html := client.GetHtml(url) // Это заглушка, необходимо реализовать функцию client.GetHtml
	return &Parser{html: html}
}

func (p *Parser) Prepare() *ParsingTree {
	html := p.html
	tokens := tokenize(html)
	root := parseTokens(tokens)

	tree := &ParsingTree{
		Root:         root,
		ElemsByTag:   make(map[string][]*Node),
		ElemsByClass: make(map[string][]*Node),
		ElemsById:    make(map[string][]*Node),
	}

	p.Tree = tree

	// Заполнение карт
	fillMaps(root, tree)

	return tree
}

func tokenize(html string) []string {
	// Простая токенизация HTML, делим по угловым скобкам и текст
	tokens := []string{}
	currentToken := ""
	inTag := false

	for i := 0; i < len(html); i++ {
		char := html[i]
		if char == '<' {
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
			inTag = true
			currentToken += string(char)
		} else if char == '>' {
			currentToken += string(char)
			tokens = append(tokens, currentToken)
			currentToken = ""
			inTag = false
		} else {
			currentToken += string(char)
			if !inTag && (i == len(html)-1 || html[i+1] == '<') {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
		}
	}
	return tokens
}

func parseTokens(tokens []string) *Node {
	stack := []*Node{}
	var root *Node

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if strings.HasPrefix(token, "<") {
			if strings.HasPrefix(token, "</") {
				// Закрывающий тег
				if len(stack) > 1 {
					node := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					parent := stack[len(stack)-1]
					parent.Children = append(parent.Children, node)
				}
			} else {
				// Открывающий тег
				tagContent := strings.Trim(token, "<>")
				tagParts := strings.Fields(tagContent)
				tagName := tagParts[0]

				node := &Node{Name: tagName, Attributes: make(map[string]string)}
				for _, part := range tagParts[1:] {
					if strings.HasPrefix(part, "Class=") {
						node.Class = strings.Trim(part[len("Class="):], `"`)
						node.Attributes["class"] = node.Class
					} else if strings.HasPrefix(part, "id=") {
						node.Id = strings.Trim(part[len("id="):], `"`)
						node.Attributes["id"] = node.Id
					} else {
						attrParts := strings.SplitN(part, "=", 2)
						if len(attrParts) == 2 {
							attrName := attrParts[0]
							attrValue := strings.Trim(attrParts[1], `"`)
							node.Attributes[attrName] = attrValue
						}
					}
				}

				if len(stack) == 0 {
					root = node
				} else {
					stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
				}
				stack = append(stack, node)
			}
		} else {
			// Текстовый контент
			if len(stack) > 0 {
				stack[len(stack)-1].Text += token
			}
		}
	}

	// Обработка корневого узла
	if len(stack) > 0 {
		root = stack[0]
		for len(stack) > 1 {
			node := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, node)
		}
	}

	return root
}

func fillMaps(node *Node, tree *ParsingTree) {
	if node == nil {
		return
	}

	tree.ElemsByTag[node.Name] = append(tree.ElemsByTag[node.Name], node)

	if node.Class != "" {
		tree.ElemsByClass[node.Class] = append(tree.ElemsByClass[node.Class], node)
	}

	if node.Id != "" {
		tree.ElemsById[node.Id] = append(tree.ElemsById[node.Id], node)
	}

	for _, child := range node.Children {
		fillMaps(child, tree)
	}
}
