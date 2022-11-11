package singleinstmodule

import (
	"fmt"
	"github.com/dan-and-dna/singleinstmodule/internal"
)

type Module = internal.Module
type SingleInstModules = internal.SingleInstModules
type SingleInstModuleCore = internal.SingleInstModuleCore
type ModuleCore = internal.ModuleCore

func Register(module Module) {
	ok := internal.GetSingleInst().Register(module)
	if !ok {
		panic(fmt.Sprintf("[module] %s is already exists", GetModuleName(module)))
	}

}

func GetModuleName(module Module) string {
	return internal.GetSingleInst().GetModuleName(module)
}

func Run(isTest bool) {
	internal.GetSingleInst().Run(isTest)
}

func RestartModule(module Module) {
	internal.GetSingleInst().RestartModule(module)
}
