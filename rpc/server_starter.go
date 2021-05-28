package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"reflect"
	"strings"
	errors2 "xblockchain/rpc/errors"

	"github.com/gorilla/websocket"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type ServerStarter struct {
	rpcserver  *rpc.Server
	serviceMap map[string]*service
}

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	numCalls  uint
}

type service struct {
	name    string                 // name of service
	rcvr    reflect.Value          // receiver of methods for the service
	typ     reflect.Type           // type of the receiver
	methods map[string]*methodType // registered methods
}

func (s *service) Call(methodName string, params interface{}) (interface{}, error) {
	mtype := s.methods[methodName]
	if mtype == nil {
		return nil, errors2.New(-32601, "Method not found")
	}
	function := mtype.method.Func
	argIsValue := false
	var argv reflect.Value
	if mtype.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(mtype.ArgType.Elem())
	} else {
		argv = reflect.New(mtype.ArgType)
		argIsValue = true
	}
	if argIsValue {
		argv = argv.Elem()
	}
	if params != nil {
		tk := reflect.TypeOf(params).Kind()
		switch tk {
		case reflect.Slice:
			paramsArr, _ := params.([]interface{})
			if len(paramsArr) != argv.NumField() {
				return nil, errors2.New(-32602, "Invalid params")
			}
			for i := 0; i < argv.NumField(); i++ {
				argv.Field(i).Set(reflect.ValueOf(paramsArr[i]))
			}
		case reflect.Map:
			paramsMap := params.(map[string]interface{})
			if len(paramsMap) != argv.NumField() {
				return nil, errors2.New(-32602, "Invalid params")
			}
			for k, v := range paramsMap {
				argvv := argv.FieldByName(k)
				if !argvv.IsValid() {
					continue
				}
				argvv.Set(reflect.ValueOf(v))
			}
		}
	}
	replyv := reflect.New(mtype.ReplyType.Elem())
	switch mtype.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(mtype.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(mtype.ReplyType.Elem(), 0, 0))
	}
	returnValues := function.Call([]reflect.Value{s.rcvr, argv, replyv})
	errInter := returnValues[0].Interface()
	if errInter != nil {
		e := errInter.(*errors2.JsonRPCError)
		return nil, e
	}
	return replyv.Interface(), nil
}

func NewServerStarter() (*ServerStarter, error) {
	return &ServerStarter{
		serviceMap: make(map[string]*service),
	}, nil
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return token.IsExported(t.Name()) || t.PkgPath() == ""
}
func suitableMethods(typ reflect.Type) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		if method.PkgPath != "" {
			continue
		}
		if mtype.NumIn() != 3 {
			continue
		}
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			continue
		}
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			continue
		}
		if !isExportedOrBuiltinType(replyType) {
			continue
		}
		if mtype.NumOut() != 1 {
			continue
		}
		if returnType := mtype.Out(0); returnType != typeOfError {
			continue
		}
		methods[mname] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
	}
	return methods

}
func (server *ServerStarter) Register(rcvr interface{}) error {
	return server.register(rcvr, "", false)
}

func (server *ServerStarter) RegisterName(name string, rcvr interface{}) error {
	return server.register(rcvr, name, true)
}

func (server *ServerStarter) register(rcvr interface{}, name string, useName bool) error {
	s := new(service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(s.rcvr).Type().Name()
	if useName {
		sname = name
	}
	if sname == "" {
		s := "rpc.Register: no service name for type " + s.typ.String()
		log.Print(s)
		return errors.New(s)
	}
	if !token.IsExported(sname) && !useName {
		s := "rpc.Register: type " + sname + " is not exported"
		log.Print(s)
		return errors.New(s)
	}
	s.name = sname
	s.methods = suitableMethods(s.typ)
	server.serviceMap[sname] = s

	return nil
}

func perParseRequestData(c *http.Request) (map[string]interface{}, error) {
	if "POST" != c.Method {
		return nil, errors.New("POST method excepted")
	}
	if nil == c.Body {
		return nil, fmt.Errorf("no POST data")
	}
	body, err := ioutil.ReadAll(c.Body)
	if err != nil {
		return nil, fmt.Errorf("errors while reading request body")
	}
	var data = make(map[string]interface{})
	decoder := json.NewDecoder(bytes.NewBuffer(body))
	decoder.UseNumber()
	err = decoder.Decode(&data)
	if nil != err {
		return nil, fmt.Errorf("errors parsing json request")
	}
	return data, nil
}

func (server *ServerStarter) handleJsonRPCRequest(data map[string]interface{}) (*int, interface{}, error) {
	// data, err := perParseRequestData(c)
	if len(data) == 0 || data != nil {
		return nil, nil, errors2.New(-32700, "Parse error")
	}
	idNumber, ok := data["id"].(json.Number)
	if !ok {
		return nil, nil, errors2.New(-32600, "Invalid Request")
	}
	id64, err := idNumber.Int64()
	if err != nil {
		return nil, nil, errors2.New(-32600, "Invalid Request")
	}
	id := int(id64)
	if data["jsonrpc"] != "2.0" {
		return &id, nil, errors2.New(-32600, "Invalid Request")
	}
	method, ok := data["method"].(string)
	mpake := strings.Split(method, ".")
	if !ok || len(mpake) != 2 {
		return &id, nil, errors2.New(-32601, "Method not found")
	}
	params := data["params"]
	service := server.serviceMap[mpake[0]]
	if service == nil {
		return &id, nil, errors2.New(-32601, "Method not found")
	}
	result, err := service.Call(mpake[1], params)
	if err != nil {
		return &id, nil, err
	}
	return &id, result, nil
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (server *ServerStarter) Run(addr string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			ws, err := upGrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			defer ws.Close()
			for {
				//读取ws中的数据
				mt, message, err := ws.ReadMessage()
				if err != nil {
					break
				}
				jsonMap := JSONToMap(string(message))
				id, data, err := server.handleJsonRPCRequest(jsonMap)

				infoMap := make(map[string]interface{})
				if err != nil {
					infoMap["jsonrpc"] = "2.0"
					infoMap["error"] = err
					infoMap["id"] = id
				} else {
					infoMap["jsonrpc"] = "2.0"
					infoMap["id"] = id
					infoMap["result"] = data
				}
				mjson, _ := json.Marshal(infoMap)
				//写入ws数据
				err = ws.WriteMessage(mt, mjson)
				if err != nil {
					break
				}
			}
		}
		asd, _ := perParseRequestData(r)
		id, data, err := server.handleJsonRPCRequest(asd)
		var infoDataJson = make(map[string]interface{})
		if err != nil {
			infoDataJson["jsonrpc"] = "2.0"
			infoDataJson["error"] = err
			infoDataJson["id"] = id
			mjson, _ := json.Marshal(infoDataJson)
			io.WriteString(w, string(mjson))
			return
		}
		infoDataJson["jsonrpc"] = "2.0"
		infoDataJson["id"] = id
		infoDataJson["result"] = data
		mjson, _ := json.Marshal(infoDataJson)
		io.WriteString(w, string(mjson))
	})
	return http.ListenAndServe(addr, nil)
}

func JSONToMap(str string) map[string]interface{} {

	var tempMap map[string]interface{}

	err := json.Unmarshal([]byte(str), &tempMap)

	if err != nil {
		panic(err)
	}

	return tempMap
}
