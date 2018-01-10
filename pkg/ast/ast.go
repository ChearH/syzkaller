// Copyright 2017 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

// Package ast parses and formats sys files.
package ast

// Pos represents source info for AST nodes.
type Pos struct {
	File string
	Off  int // byte offset, starting at 0
	Line int // line number, starting at 1
	Col  int // column number, starting at 1 (byte count)
}

// Description contains top-level nodes of a parsed sys description.
type Description struct {
	Nodes []Node
}

// Node is AST node interface.
type Node interface {
	Info() (pos Pos, typ string, name string)
	// Clone makes a deep copy of the node.
	// If newPos is not zero, sets Pos of all nodes to newPos.
	// If newPos is zero, Pos of nodes is left intact.
	Clone(newPos Pos) Node
	// Walk calls callback cb for all child nodes of this node.
	// Note: it's not recursive. Use Recursive helper for recursive walk.
	Walk(cb func(Node))
}

// Top-level AST nodes:

type NewLine struct {
	Pos Pos
}

func (n *NewLine) Info() (Pos, string, string) {
	return n.Pos, tok2str[tokNewLine], ""
}

type Comment struct {
	Pos  Pos
	Text string
}

func (n *Comment) Info() (Pos, string, string) {
	return n.Pos, tok2str[tokComment], ""
}

type Include struct {
	Pos  Pos
	File *String
}

func (n *Include) Info() (Pos, string, string) {
	return n.Pos, tok2str[tokInclude], ""
}

type Incdir struct {
	Pos Pos
	Dir *String
}

func (n *Incdir) Info() (Pos, string, string) {
	return n.Pos, tok2str[tokInclude], ""
}

type Define struct {
	Pos   Pos
	Name  *Ident
	Value *Int
}

func (n *Define) Info() (Pos, string, string) {
	return n.Pos, tok2str[tokDefine], n.Name.Name
}

type Resource struct {
	Pos    Pos
	Name   *Ident
	Base   *Type
	Values []*Int
}

func (n *Resource) Info() (Pos, string, string) {
	return n.Pos, tok2str[tokResource], n.Name.Name
}

type Call struct {
	Pos      Pos
	Name     *Ident
	CallName string
	NR       uint64
	Args     []*Field
	Ret      *Type
}

func (n *Call) Info() (Pos, string, string) {
	return n.Pos, "syscall", n.Name.Name
}

type Struct struct {
	Pos      Pos
	Name     *Ident
	Fields   []*Field
	Attrs    []*Ident
	Comments []*Comment
	IsUnion  bool
}

func (n *Struct) Info() (Pos, string, string) {
	typ := "struct"
	if n.IsUnion {
		typ = "union"
	}
	return n.Pos, typ, n.Name.Name
}

type IntFlags struct {
	Pos    Pos
	Name   *Ident
	Values []*Int
}

func (n *IntFlags) Info() (Pos, string, string) {
	return n.Pos, "flags", n.Name.Name
}

type StrFlags struct {
	Pos    Pos
	Name   *Ident
	Values []*String
}

func (n *StrFlags) Info() (Pos, string, string) {
	return n.Pos, "string flags", n.Name.Name
}

type TypeDef struct {
	Pos  Pos
	Name *Ident
	Type *Type
}

func (n *TypeDef) Info() (Pos, string, string) {
	return n.Pos, "type", n.Name.Name
}

// Not top-level AST nodes:

type Ident struct {
	Pos  Pos
	Name string
}

func (n *Ident) Info() (Pos, string, string) {
	return n.Pos, tok2str[tokIdent], n.Name
}

type String struct {
	Pos   Pos
	Value string
}

func (n *String) Info() (Pos, string, string) {
	return n.Pos, tok2str[tokString], ""
}

type Int struct {
	Pos Pos
	// Only one of Value, Ident, CExpr is filled.
	Value    uint64
	ValueHex bool // says if value was in hex (for formatting)
	Ident    string
	CExpr    string
}

func (n *Int) Info() (Pos, string, string) {
	return n.Pos, tok2str[tokInt], ""
}

type Type struct {
	Pos Pos
	// Only one of Value, Ident, String is filled.
	Value    uint64
	ValueHex bool
	Ident    string
	String   string
	// Part after COLON (for ranges and bitfields).
	HasColon  bool
	Pos2      Pos
	Value2    uint64
	Value2Hex bool
	Ident2    string
	Args      []*Type
}

func (n *Type) Info() (Pos, string, string) {
	return n.Pos, "type", n.Ident
}

type Field struct {
	Pos      Pos
	Name     *Ident
	Type     *Type
	NewBlock bool // separated from previous fields by a new line
	Comments []*Comment
}

func (n *Field) Info() (Pos, string, string) {
	return n.Pos, "arg/field", n.Name.Name
}
