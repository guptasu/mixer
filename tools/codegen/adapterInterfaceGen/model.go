package main

type Model struct {
	// top level fields
	Name        string
	Check       bool
	PackageName string
	VarietyName string

	// types
	TypeFullName string

	// imports
	Imports []string

	ConstructorFields []FieldInfo
}

type FieldInfo struct {
	Name string
	Type TypeInfo
}

type TypeInfo struct {
	Name   string
	IsExpr bool

	IsMap     bool
}

const FullNameOfExprMessage = "*istio_mixer_v1_config_template.Expr"
