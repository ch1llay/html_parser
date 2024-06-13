package main

import (
	"fmt"
	parser "html_parser/internal/parser"
	"strings"
)

func main() {
	html := `<html><body><div class="container" id="main"><p data-info="greeting">Hello, world!</p></div></body></html>`
	p := parser.NewParser(html)
	tree := p.Prepare()

	fmt.Println("Root node:", tree.Root.Name)
	fmt.Println("Root text:", tree.Root.Text)
	fmt.Println("Nodes by tag 'div':", len(tree.ElemsByTag["div"]))
	fmt.Println("Nodes by class 'container':", len(tree.ElemsByClass["container"]))
	fmt.Println("Nodes by id 'main':", len(tree.ElemsById["main"]))

	// Вывод текстового содержимого всех узлов
	printNodeText(tree.Root, 0)
}

func printNodeText(node *parser.Node, level int) {
	if node == nil {
		return
	}

	fmt.Printf("%s<%s", strings.Repeat(" ", level*2), node.Name)
	for attr, val := range node.Attributes {
		fmt.Printf(" %s=\"%s\"", attr, val)
	}
	fmt.Printf(">\n")

	if node.Text != "" {
		fmt.Printf("%s%s\n", strings.Repeat(" ", level*2+2), node.Text)
	}
	for _, child := range node.Children {
		printNodeText(child, level+1)
	}
	fmt.Printf("%s</%s>\n", strings.Repeat(" ", level*2), node.Name)
}
