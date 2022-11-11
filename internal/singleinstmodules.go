package internal

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"runtime/debug"
	"sync"
	"syscall"
)

var (
	singleInst *SingleInstModules = nil
	once       sync.Once
)

type SingleInstModules struct {
	modules         []Module
	moduleMap       map[uint64]Module
	moduleTypes     map[reflect.Type]struct{}
	coreChangedChan chan Module
	rex             *regexp.Regexp
	mu              sync.Mutex
	wg              sync.WaitGroup
	moduleId        uint64
}

func (single *SingleInstModules) Construct() {
	if single.coreChangedChan == nil {
		single.coreChangedChan = make(chan Module, 30)
	}

	if single.moduleTypes == nil {
		single.moduleTypes = make(map[reflect.Type]struct{})
	}

	if single.moduleId == 0 {
		single.moduleId = 10000
	}

	if single.moduleMap == nil {
		single.moduleMap = make(map[uint64]Module)
	}
}

func (single *SingleInstModules) Register(module Module) bool {
	moduleType := reflect.TypeOf(module)

	single.mu.Lock()
	defer single.mu.Unlock()

	if _, ok := single.moduleTypes[moduleType]; ok {
		return false
	}

	single.moduleTypes[moduleType] = struct{}{}
	single.modules = append(single.modules, module)

	module.ModuleConstruct()

	return true
}

func (single *SingleInstModules) GetModuleName(module Module) string {
	if module != nil {
		return reflect.TypeOf(module).Elem().Name()
	}
	panic("unknown module name")
}

func (single *SingleInstModules) Run(isTest bool) {

	// first construct
	defer single.Destruct()
	defer single.Shutdown()
	log.Println("============run modules============")
	single.constructAndRun()

	/*
		if wait {
			<-single.loop()

			// final destroy
			single.Shutdown()
		}

	*/
	log.Println("============all done============")
	if isTest {
		return
	}

	for {
		select {
		case module := <-single.coreChangedChan:
			if module != nil {
				if canRestart, ok := module.(ModuleCanRestart); ok {
					if ok := canRestart.ModuleRestart(); ok {
						if canAfterRestart, ok := module.(ModuleCanAfterRestart); ok {
							canAfterRestart.AfterRestart()
						}
					}
				}
			}

		case <-single.loop():

			// final destroy

			return
		}
	}
}

func (single *SingleInstModules) RestartModule(module Module) {
	single.coreChangedChan <- module
}

func (single *SingleInstModules) constructAndRun() {
	single.mu.Lock()
	defer single.mu.Unlock()
	numModules := len(single.modules)
	// first construct
	//for _, module := range single.modules {
	//	module.ModuleConstruct()
	//}

	// final run
	for i := numModules - 1; i >= 0; i-- {
		single.runOneModule(single.modules[i])
	}
}

func (single *SingleInstModules) runOneModule(module Module) {
	if module == nil {
		return
	}

	var runFuncList []func()
	moduleType := reflect.TypeOf(module)
	for i := 0; i < moduleType.NumMethod(); i++ {
		//  lexicographic order
		method := moduleType.Method(i)
		methodName := method.Name
		if !single.isRunFunc(methodName) {
			continue
		}

		newRunFunc := func() {
			defer recoverOnPanic()
			defer func() {
				r := recover()

				if canAfterRun, ok := module.(ModuleCanAfterRun); ok {
					canAfterRun.ModuleAfterRun(methodName)
				}

				onPanic(r)
			}()

			if canBeforeRun, ok := module.(ModuleCanBeforeRun); ok {
				canBeforeRun.ModuleBeforeRun(methodName)
			}

			// call module.Run*()
			//single.wg.Add(1)

			//defer single.wg.Done()
			method.Func.Call([]reflect.Value{reflect.ValueOf(module)})

		}

		runFuncList = append(runFuncList, newRunFunc)
	}

	for _, runFunc := range runFuncList {
		runFunc()
	}
}
func recoverOnPanic() {
	r := recover()
	onPanic(r)
}

func onPanic(r interface{}) {
	if r != nil {
		log.Printf("moudle panicked: %v\n%s", r, debug.Stack())
	}
}

func (single *SingleInstModules) isRunFunc(methodName string) bool {
	if single.rex == nil {
		single.rex = regexp.MustCompile("^ModuleRun")
	}
	return single.rex.MatchString(methodName)
}

func (single *SingleInstModules) Shutdown() {
	// deps final shutdown
	log.Println("============shutdown modules============")
	single.mu.Lock()
	defer single.mu.Unlock()
	numModules := len(single.modules)
	for i := numModules - 1; i >= 0; i-- {
		module := single.modules[i]
		if canDown, ok := module.(ModuleCanShutdown); ok {
			canDown.ModuleShutdown()
		}
	}
	//single.wg.Wait()
}

func (single *SingleInstModules) Destruct() {
	// deps final destroy
	log.Println("============destroy modules============")
	single.mu.Lock()
	defer single.mu.Unlock()
	numModules := len(single.modules)
	for i := numModules - 1; i >= 0; i-- {
		module := single.modules[i]
		module.ModuleDestruct()
	}

	single.moduleTypes = nil
	single.modules = nil
	single.rex = nil
}

func (single *SingleInstModules) loop() chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return quit
}

func GetSingleInst() *SingleInstModules {
	if singleInst == nil {
		once.Do(func() {
			singleInst = new(SingleInstModules)
			singleInst.Construct()
		})
	}

	return singleInst
}
