package handling

import (
	"amaterasu/cqrs/infrastructure/utils"
	v2messaging "amaterasu/cqrs/infrastructure/v2.messaging"
	"errors"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

//
const (
	PREFIX = "Handle"
	NUMIN  = 2
	NUMOUT = 1
)

// the handle func for specific command in conmand handler class
// ex:
/*
	type fooCommandHandler struct{} // handler
	type foo1 struct{} // command
	type foo2 struct{} // command

	func (fh fooCommandHandler) handleFoo1(f foo1) error {

	}

	func (fh fooCommandHandler) handleFoo2(f foo2) error {

	}
*/
type commandHandlerFunc func(msg v2messaging.Message) error

// dispatch command to specific command handler
// note : 1 command is handled only by 1 commandhandler, failed if 2 or 3 ...
type CommandDispatcher struct {
	// list of command handler
	commandHandlers map[reflect.Type]v2messaging.ICommandHandler
	// list type of command, for tricky golang :))
	commandTypes map[string]reflect.Type
	// handler func to handle command
	// ex [{foo1,handleFoo1} - {foo2,handleFoo2} - {...,...}]
	commandHandlerFuncs map[string]commandHandlerFunc
}

// create new command dispatcher
func New_CommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		commandHandlers:     make(map[reflect.Type]v2messaging.ICommandHandler),
		commandTypes:        make(map[string]reflect.Type),
		commandHandlerFuncs: make(map[string]commandHandlerFunc),
	}
}

// Registers the specified command handler.
func (cd *CommandDispatcher) Register(commandHandler v2messaging.ICommandHandler) error {

	handlerType := reflect.TypeOf(commandHandler)
	numHandlerMethods := handlerType.NumMethod() // count all method

	if numHandlerMethods == 0 {
		return errors.New("must have handleCommand method, please read the document please")
	}

	handleFuncCount := 0 // for count handle command func inside commandHandler

	for i := 0; i < numHandlerMethods; i++ {
		// get a method in commandHandler
		handlerMethod := handlerType.Method(i)

		if !strings.HasPrefix(handlerMethod.Name, PREFIX) || // skip method without [PREFIX]
			handlerMethod.Type.NumIn() != NUMIN || // only handle method should have [NUMIN] arguments,
			handlerMethod.Type.NumOut() != NUMOUT { // [NUMOUT] outputs
			continue
		}
		// yeahhh, we have handle command func in command handler, so lucky
		handleFuncCount++
		// ex:
		// func (fh fooCommandHandler) handleFoo1(f foo1) error {}
		// func (fh fooCommandHandler) handleFoo2(f foo2) error {}
		handlerFunc := func(msg v2messaging.Message) error {

			response := handlerMethod.Func.Call([]reflect.Value{
				reflect.ValueOf(commandHandler), // fh param
				reflect.ValueOf(msg).Elem(),     // f param
			})

			// can use class if this class is nil
			// ex: func (fh fooCommandHandler) handleFoo1(f foo1) error {}
			// if fh nil, can't call `handleFoo1` func
			if len(response) > 0 && !response[0].IsNil() {
				err := response[0].Interface().(error)
				return err
			}

			return nil
		}

		// 1 because param in golang start at 1, 0 is struct type
		commandType := handlerMethod.Type.In(1) // get command type

		// 1 command is handled only by 1 commandhandler, failed if 2 or 3 ...
		if _, ok := cd.commandHandlers[commandType]; ok {
			return errors.New("the command handled by the received handler already has a registered handler")
		}

		// command type as string
		commandTypeName := utils.GetTypeName2(commandType)
		// store
		cd.commandHandlers[commandType] = commandHandler
		cd.commandTypes[commandTypeName] = commandType
		cd.commandHandlerFuncs[commandTypeName] = handlerFunc

	}

	if handleFuncCount == 0 {
		return errors.New("must have handleCommand method, please read the document please")
	}

	return nil
}

// Processes the message by calling the registered handler.
func (cd *CommandDispatcher) ProcessMessage(msg v2messaging.Envelope) error {
	handlerFunc, found := cd.commandHandlerFuncs[msg.MsgType]

	if !found {
		return errors.New("not found consumer handler")
	}

	entry, exists := cd.commandTypes[msg.MsgType]

	if !exists {
		return errors.New("not found command type")
	}

	command := reflect.New(entry).Interface().(v2messaging.Message)
	config := &mapstructure.DecoderConfig{
		DecodeHook:       utils.MapTimeFromJSON,
		TagName:          "json",
		Result:           command,
		WeaklyTypedInput: true,
	}

	decoder, errDecoder := mapstructure.NewDecoder(config)
	if errDecoder != nil {
		return errors.New("config mapstructure fail")
	}

	errDecode := decoder.Decode(msg.Body)
	if errDecode != nil {
		return errors.New("decode message fail")
	}

	return handlerFunc(command)
}
