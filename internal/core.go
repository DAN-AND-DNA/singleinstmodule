package internal

type Module interface {
	ModuleConstruct() // for module
	ModuleDestruct()  // for module
	ModuleLock() ModuleCore
	ModuleUnlock()
}

type ModuleCanRestart interface {
	ModuleRestart() bool // for module
}

type ModuleCanBeforeRun interface {
	ModuleBeforeRun(string) // for run
}

type ModuleCanAfterRun interface {
	ModuleAfterRun(string) // for run
}

type ModuleCanShutdown interface {
	ModuleShutdown() // for run
}

type ModuleCanAfterRestart interface {
	AfterRestart()
}

type ModuleCore interface {
	Lock()
	Unlock()
}
