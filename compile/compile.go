package compile

import (
	"github.com/patrickz98/project.go.omnetpp/omnetpp"
	"github.com/patrickz98/project.go.omnetpp/simple"
)

func (compiler Compiler) Opp() (opp omnetpp.OmnetProject) {
	oppConf := omnetpp.Config{
		OppConfig: compiler.OppConfig,
		Path:      compiler.SimulationBase,
	}

	opp = omnetpp.New(&oppConf)

	return
}

func (compiler Compiler) Clean() (err error) {

	opp := compiler.Opp()
	err = opp.Clean()
	if err != nil {
		return
	}

	compiler.cleanedFiles, err = simple.ListDir(opp.Path)

	return
}

func (compiler Compiler) Compile() (err error) {

	opp := compiler.Opp()
	err = opp.Setup(false)
	if err != nil {
		return
	}

	compiler.compiledFiles, err = simple.ListDir(opp.Path)
	if err != nil {
		return
	}

	return
}

func (compiler Compiler) CompiledFiles() (files map[string]bool) {
	files = simple.DirDiff(compiler.cleanedFiles, compiler.compiledFiles)
	return
}
