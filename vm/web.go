package vm

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"net/http"
	"ryanlang/compiler"
	"ryanlang/compiler/instruction"
	"ryanlang/object"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type debuggerNotificationSignal int
type debuggerNotification struct {
	signal  debuggerNotificationSignal
	payload interface{}
}

const (
	vmStateChanged debuggerNotificationSignal = iota
	vmHaltMsg
	vmBreakpointsUpdated
)

type webObject struct {
	Type string      `json:"type"`
	Id   int         `json:"id"`
	Data interface{} `json:"data"`
}
type webStackItem struct {
	Index  int       `json:"index"`
	Object webObject `json:"object"`
}
type webLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}
type webObjectCodeInstruction struct {
	Addr        int          `json:"addr"`
	Instruction string       `json:"instruction"`
	Comment     string       `json:"comment"`
	Location    *webLocation `json:"location"`
}
type webObjectCode struct {
	Instructions []webObjectCodeInstruction `json:"instructions"`
}
type webFile struct {
	Name     string `json:"name"`
	Contents string `json:"contents"`
}
type webMessageType int

const (
	webMessageVmState webMessageType = iota
	webMessageCodeDump
	webMessageHaltMsg
)

type webMessage struct {
	Typ     webMessageType `json:"typ"`
	Payload interface{}    `json:"payload"`
}
type webBreakpointLocation struct {
	ObjectID int `json:"objectID"`
	Address  int `json:"address"`
}
type webBreakpoints struct {
	Locations []webBreakpointLocation `json:"locations"`
}
type webState struct {
	State       state          `json:"state"`
	Breakpoints webBreakpoints `json:"breakpoints"`
}
type webSymbolInfo struct {
	Id       int          `json:"id"`
	Name     string       `json:"name"`
	Location *webLocation `json:"location"`
}
type webCodeDumpFrame struct {
	Cp       int             `json:"cp"`
	Bsp      int             `json:"bsp"`
	ObjectID int             `json:"objectID"`
	Locals   []webSymbolInfo `json:"locals"`
	Foreigns []webSymbolInfo `json:"foreigns"`
}
type webCodeDump struct {
	Objects []webObject        `json:"objects"`
	Frames  []webCodeDumpFrame `json:"frames"`
	Files   []webFile          `json:"files"`
	Stack   []webStackItem     `json:"stack"`
}

type webDebugger struct {
	v             *VM
	client        *webClient
	logger        *zap.SugaredLogger
	notifications chan debuggerNotification
	commands      chan webCommand
}
type webClient struct {
	ws *websocket.Conn
}

type requestType int

const (
	requestTypeDumpState requestType = iota
	requestTypeNext
	requestTypeRun
	requestTypeStepOut
	requestTypeToggleBreakpoint
)

type webRequest struct {
	Typ     requestType `json:"typ"`
	Payload interface{} `json:"payload"`
}

func newWebDebugger(v *VM) *webDebugger {
	logger, _ := zap.NewDevelopment()

	return &webDebugger{
		v:             v,
		client:        nil,
		logger:        logger.Sugar(),
		notifications: make(chan debuggerNotification),
	}
}

func (wd *webDebugger) webDumpCode(obj *object.Code, dd *compiler.DebugData) (*webObjectCode, error) {
	result := &webObjectCode{}
	cp := 0
	for cp < len(obj.Code) {
		inst, n, err := instruction.Read(obj.Code[cp:])
		if err != nil {
			return nil, err
		}

		var loc *webLocation
		var comment string
		switch inst.Op() {
		case instruction.OpAnnotation:
			annotationObject, ok := wd.v.objects.Get(inst.(instruction.Annotation).Index)
			if ok {
				comment = annotationObject.(*object.String).Value
			}
		case instruction.OpPushConstant:
			obj, ok := wd.v.objects.Get(inst.(instruction.PushConstant).Index)
			if ok {
				comment = fmt.Sprintf("%s %s", obj.Type().String(), obj.String())
			}
		case instruction.OpPushLocalRef:
			comment = fmt.Sprintf("local \"%s\"", dd.LocalName(int(inst.(instruction.PushLocal).Index)))
		case instruction.OpStoreLocal:
			comment = fmt.Sprintf("local \"%s\"", dd.LocalName(int(inst.(instruction.StoreLocal).Index)))
		case instruction.OpPushForeign:
			comment = fmt.Sprintf("foreign \"%s\"", dd.ForeignName(int(inst.(instruction.PushForeign).Index)))
		case instruction.OpStoreForeign:
			comment = fmt.Sprintf("foreign \"%s\"", dd.ForeignName(int(inst.(instruction.StoreForeign).Index)))
		}

		entry := dd.SearchSource(cp)
		if entry != nil {
			l := entry.Location()
			if l != nil {
				loc = &webLocation{
					File:   l.File,
					Line:   l.Line,
					Column: l.Column,
				}
			}
		}
		result.Instructions = append(result.Instructions, webObjectCodeInstruction{
			Addr:        cp,
			Instruction: inst.String(),
			Comment:     comment,
			Location:    loc,
		})
		cp += n
	}
	return result, nil
}
func (wd *webDebugger) webDumpFiles() []webFile {
	result := []webFile{}
	for fn, contents := range wd.v.sourceFiles {
		result = append(result, webFile{
			Name:     fn,
			Contents: contents,
		})
	}
	return result
}
func (wd *webDebugger) send(msg *webMessage, client *webClient) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if err := client.ws.WriteMessage(websocket.BinaryMessage, b); err != nil {
		return err
	}
	return nil
}
func (wd *webDebugger) sendVmState(client *webClient) error {
	codeToObjectIDs := wd.getCodeToObjectIDs()
	bpLocations := make([]webBreakpointLocation, 0)
	for code, addrs := range wd.v.bp.locations {
		for addr, ok := range addrs {
			if !ok {
				continue
			}
			bpLocations = append(bpLocations, webBreakpointLocation{
				ObjectID: codeToObjectIDs[code],
				Address:  addr,
			})
		}
	}

	return wd.send(&webMessage{
		Typ: webMessageVmState,
		Payload: webState{
			State: wd.v.state,
			Breakpoints: webBreakpoints{
				Locations: bpLocations,
			},
		},
	}, client)
}
func (wd *webDebugger) sendHaltMsg(client *webClient, err error) error {
	return wd.send(&webMessage{
		Typ:     webMessageHaltMsg,
		Payload: err.Error(),
	}, client)
}
func (wd *webDebugger) getCodeToObjectIDs() map[*object.Code]int {
	codeToObjectIDs := map[*object.Code]int{}
	for i := uint16(0); i < wd.v.objects.Len(); i++ {
		obj, _ := wd.v.objects.Get(i)
		if obj.Type() == object.CODE {
			codeToObjectIDs[obj.(*object.Code)] = int(i)
		}
	}
	return codeToObjectIDs
}
func (wd *webDebugger) sendCodeDump(client *webClient) error {
	s := webCodeDump{
		Objects: make([]webObject, 0),
		Files:   wd.webDumpFiles(),
		Stack:   make([]webStackItem, 0),
	}
	codeToObjectIDs := wd.getCodeToObjectIDs()

	for i := uint16(0); i < wd.v.objects.Len(); i++ {
		obj, _ := wd.v.objects.Get(i)
		wobj := webObject{
			Type: obj.Type().String(),
			Id:   int(i),
			Data: nil,
		}

		if obj.Type() == object.CODE {
			dumpedCode, err := wd.webDumpCode(obj.(*object.Code), wd.v.debugData[int(i)])
			if err != nil {
				return err
			}
			wobj.Data = dumpedCode
		}

		s.Objects = append(s.Objects, wobj)
	}

	frames := make([]webCodeDumpFrame, 0)
	for fi := 0; fi <= wd.v.fp; fi++ {
		frame := wd.v.frames[fi]

		objId, ok := codeToObjectIDs[frame.cl.Code]
		if !ok {
			objId = -1
		}

		var locs []webSymbolInfo
		var fors []webSymbolInfo
		if dd, ok := wd.v.debugData[objId]; ok {
			for i, s := range dd.Locals {
				locs = append(locs, webSymbolInfo{
					Id:   i,
					Name: s.Name,
				})
			}
			for i, s := range dd.Foreigns {
				fors = append(fors, webSymbolInfo{
					Id:   i,
					Name: s.Name,
				})
			}
		}

		frames = append(frames, webCodeDumpFrame{
			Cp:       frame.cpe,
			ObjectID: objId,
			Bsp:      frame.bsp,
			Locals:   locs,
			Foreigns: fors,
		})
	}
	s.Frames = frames

	for si := 0; si <= wd.v.sp; si++ {
		value := *wd.v.stack[si]

		s.Stack = append(s.Stack, webStackItem{
			Index: si,
			Object: webObject{
				Type: value.Type().String(),
				Data: value.String(),
			},
		})
	}

	return wd.send(&webMessage{
		Typ:     webMessageCodeDump,
		Payload: s,
	}, client)
}

func (wd *webDebugger) handleRequest(request *webRequest, client *webClient) error {
	switch request.Typ {
	case requestTypeDumpState:
		wd.sendVmState(client)
		wd.sendCodeDump(client)
	case requestTypeNext:
		if wd.commands != nil {
			wd.commands <- webCommand{typ: webCommandNext}
		}
	case requestTypeRun:
		if wd.commands != nil {
			wd.commands <- webCommand{typ: webCommandRun}
		}
	case requestTypeStepOut:
		if wd.commands != nil {
			wd.commands <- webCommand{typ: webCommandStepOut}
		}
	case requestTypeToggleBreakpoint:
		if wd.commands != nil {
			wd.commands <- webCommand{
				typ: webCommandToggleBreakpoint,
				payload: webCommandPayloadToggleBreakpoint{
					ObjectID: int(request.Payload.(map[string]interface{})["objectID"].(float64)),
					Address:  int(request.Payload.(map[string]interface{})["address"].(float64)),
				},
			}
		}
	default:
		return fmt.Errorf("unknown request type: %d", request.Typ)
	}
	return nil
}
func (wd *webDebugger) webReader() {
	defer wd.removeClient()
	for {
		mt, b, err := wd.client.ws.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, 1001, 1005, 1006) {
				wd.logger.With("err", err).Error("failed to read mesage")
			}
			break
		}
		if mt != websocket.BinaryMessage {
			wd.logger.Error("dont know how to handle non-binary message requests")
			continue
		}
		var request webRequest
		if err = json.Unmarshal(b, &request); err != nil {
			wd.logger.With("err", err).Error("cannot parse json request")
			continue
		}
		if err = wd.handleRequest(&request, wd.client); err != nil {
			wd.logger.With("err", err).Error("cannot handle request")
			continue
		}
	}
}
func (wd *webDebugger) webServeWs(w http.ResponseWriter, r *http.Request) {
	if wd.client != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	wd.client = &webClient{
		ws: ws,
	}

	//go wd.webWriter()
	wd.webReader()
}
func (wd *webDebugger) removeClient() {
	if wd.client == nil {
		return
	}
	c := wd.client
	wd.client = nil
	c.ws.Close()
	//close(c.notifications)
}
func (wd *webDebugger) notificationListener() {
	for {
		select {
		case notification := <-wd.notifications:
			switch notification.signal {
			case vmStateChanged:
				if wd.client != nil {
					wd.sendVmState(wd.client)
					if wd.v.state == statePaused {
						wd.sendCodeDump(wd.client)
					}
				}
			case vmHaltMsg:
				if wd.client != nil {
					wd.sendHaltMsg(wd.client, notification.payload.(error))
				}
			case vmBreakpointsUpdated:
				if wd.client != nil {
					wd.sendVmState(wd.client)
				}
			default:
				panic("unknown signal for web debugger")
			}
		}
	}
}
func (wd *webDebugger) Notify(signal debuggerNotificationSignal, payload interface{}) {
	wd.notifications <- debuggerNotification{
		signal:  signal,
		payload: payload,
	}
}
func (wd *webDebugger) start() {
	go wd.notificationListener()
	http.HandleFunc("/ws", wd.webServeWs)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type webCommandPayloadToggleBreakpoint struct {
	ObjectID int `json:"objectID"`
	Address  int `json:"address"`
}
type webCommand struct {
	typ     webCommandType
	payload interface{}
}
type webCommandType int

const (
	webCommandNext webCommandType = iota
	webCommandRun
	webCommandStepOut
	webCommandToggleBreakpoint
)

func (wd *webDebugger) Listen() <-chan webCommand {
	wd.commands = make(chan webCommand)
	return wd.commands
}
func (wd *webDebugger) StopListening() {
	close(wd.commands)
	wd.commands = nil
}
