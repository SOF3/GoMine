package interfaces

type ILevel interface {
	GetServer() IServer
	GetName() string
	GetDimensions() map[string]IDimension
	AddDimension(IDimension)
	DimensionExists(string) bool
	RemoveDimension(string) bool
	SetDefaultDimension(IDimension)
	GetDefaultDimension() IDimension
	TickLevel()
	GetGameRules() map[string]IGameRule
	GetGameRule(string) IGameRule
	IsGenerated() bool
	SetGenerator(IGenerator)
	GetGenerator() IGenerator
	GenerateChunks()
}

type IGameRule interface {
	GetName() string
	GetValue() interface{}
	SetValue(interface{}) bool
}
