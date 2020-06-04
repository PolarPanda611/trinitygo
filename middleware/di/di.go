package di

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/gin-gonic/gin"
)

// New DI middleware
func New(app application.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := fmt.Sprintf("%v@%v", c.Request.Method, c.FullPath())
		runtimeKeyMap := application.DecodeHTTPRuntimeKey(c, app.RuntimeKeys())
		tContext := app.ContextPool().Acquire(app, runtimeKeyMap, c)
		defer func() {
			//release trinity go context obj
			app.ContextPool().Release(tContext)
		}()
		controller, sharedInstance := app.ControllerPool().GetController(method, tContext, app, c)
		defer func() {
			for _, v := range sharedInstance {
				app.InstancePool().Release(v)
			}
		}()
		validators := app.ControllerPool().GetControllerValidators(method)
		for _, v := range validators {
			v(tContext)
			if c.IsAborted() {
				return
			}
		}
		funcName, ok := app.ControllerPool().GetControllerFuncName(method)
		if ok && funcName == "" || !ok {
			funcName = c.Request.Method
		}
		currentMethod, ok := reflect.TypeOf(controller).MethodByName(funcName)
		if !ok {
			panic("controller has no method ")
		}
		// validation passed , run handler
		if currentMethod.Type.NumIn() < -1 {
			panic("wrong method")
		} else if currentMethod.Type.NumIn() == 1 {
			var inParam []reflect.Value                            // 构造函数入参 ， 入参1 ， transport指针对象
			inParam = append(inParam, reflect.ValueOf(controller)) // 传入transport对象
			currentMethod.Func.Call(inParam)                       // 调用transport函数，传入参数
		} else {
			var inParam []reflect.Value                            // 构造函数入参 ， 入参1 ， transport指针对象
			inParam = append(inParam, reflect.ValueOf(controller)) // 传入transport对象

			for i := 1; i < currentMethod.Type.NumIn(); i++ {
				if currentMethod.Type.In(i).Kind() != reflect.Struct {
					panic(fmt.Sprintf("controller in-param type expect struct , actual %v , wrong type", currentMethod.Type.In(i).Kind()))
				}
				destVal := reflect.Indirect(reflect.New(currentMethod.Type.In(i)))
				for index := 0; index < destVal.NumField(); index++ {
					val := destVal.Field(index)
					if !val.CanSet() {
						panic(fmt.Sprintf("param : %v is not exported , cannot set", currentMethod.Type.In(i).Field(index).Name))
					}
					// check if header param
					if headerParam, isExist := currentMethod.Type.In(i).Field(index).Tag.Lookup("header_param"); isExist {
						headerValString := c.GetHeader(headerParam)
						switch currentMethod.Type.In(i).Field(index).Type.Kind() {
						case reflect.Int64:
							headerVal, err := strconv.ParseInt(headerValString, 10, 64)
							if err != nil {
								panic(err)
							}
							val.Set(reflect.ValueOf(headerVal))
							break
						case reflect.Int32:
							headerVal, err := strconv.ParseInt(headerValString, 10, 32)
							if err != nil {
								panic(err)
							}
							val.Set(reflect.ValueOf(headerVal))
							break
						case reflect.Int16:
							headerVal, err := strconv.ParseInt(headerValString, 10, 16)
							if err != nil {
								panic(err)
							}
							val.Set(reflect.ValueOf(headerVal))
							break
						case reflect.Int:
							headerVal, err := strconv.Atoi(headerValString)
							if err != nil {
								panic(err)
							}
							val.Set(reflect.ValueOf(headerVal))
							break
						case reflect.String:
							val.Set(reflect.ValueOf(headerValString))
							break
						default:
							panic("Unsupported type, only support int64 , int32 , int16 , int , string ")
						}
					}
					// check if path param
					if pathParam, isExist := currentMethod.Type.In(i).Field(index).Tag.Lookup("path_param"); isExist {
						paramValString, ok := c.Params.Get(pathParam)
						if !ok {
							panic(fmt.Sprintf("%v param not exist ", pathParam))
						}
						switch currentMethod.Type.In(i).Field(index).Type.Kind() {
						case reflect.Int64:
							paramVal, err := strconv.ParseInt(paramValString, 10, 64)
							if err != nil {
								panic(err)
							}
							val.Set(reflect.ValueOf(paramVal))
							break
						case reflect.Int32:
							paramVal, err := strconv.ParseInt(paramValString, 10, 32)
							if err != nil {
								panic(err)
							}
							val.Set(reflect.ValueOf(paramVal))
							break
						case reflect.Int:
							paramVal, err := strconv.Atoi(paramValString)
							if err != nil {
								panic(err)
							}
							val.Set(reflect.ValueOf(paramVal))
							break
						case reflect.String:
							val.Set(reflect.ValueOf(paramValString))
							break
						default:
							panic("Unsupported type")
						}
					}
					// check if query param
					if queryParam, isExist := currentMethod.Type.In(i).Field(index).Tag.Lookup("query_param"); isExist {
						if queryParam == "" {
							switch currentMethod.Type.In(i).Field(index).Type.Kind() {
							case reflect.String:
								val.Set(reflect.ValueOf(c.Request.URL.RawQuery))
								break
							default:
								panic("Unsupported type , only support string ")
							}
						} else {
							queryValString := c.Query(queryParam)
							switch currentMethod.Type.In(i).Field(index).Type.Kind() {
							case reflect.Int64:
								queryVal, err := strconv.ParseInt(queryValString, 10, 64)
								if err != nil {
									panic(err)
								}
								val.Set(reflect.ValueOf(queryVal))
								break
							case reflect.Int32:
								queryVal, err := strconv.ParseInt(queryValString, 10, 32)
								if err != nil {
									panic(err)
								}
								val.Set(reflect.ValueOf(queryVal))
								break
							case reflect.Int16:
								queryVal, err := strconv.ParseInt(queryValString, 10, 16)
								if err != nil {
									panic(err)
								}
								val.Set(reflect.ValueOf(queryVal))
								break
							case reflect.Int:
								queryVal, err := strconv.Atoi(queryValString)
								if err != nil {
									panic(err)
								}
								val.Set(reflect.ValueOf(queryVal))
								break
							case reflect.String:
								val.Set(reflect.ValueOf(queryValString))
								break
							default:
								panic("Unsupported type , only support int64 , int32 , int16 , int , string ")
							}
						}

					}
					// check if body param
					if bodyParam, isExist := currentMethod.Type.In(i).Field(index).Tag.Lookup("body_param"); isExist {
						respBytes, err := ioutil.ReadAll(c.Request.Body)
						c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(respBytes))
						if err != nil {
							panic(err)
						}
						if bodyParam == "" {
							switch currentMethod.Type.In(i).Field(index).Type.Kind() {
							case reflect.String:
								val.Set(reflect.ValueOf(string(respBytes)))
								break
							case reflect.Struct, reflect.Slice:
								targetVal := reflect.New(currentMethod.Type.In(i).Field(index).Type).Interface()
								if err := c.BindJSON(targetVal); err != nil {
									panic(err)
								}
								val.Set(reflect.Indirect(reflect.ValueOf(targetVal)))
								break
							case reflect.Map:
								if fmt.Sprintf("%v", currentMethod.Type.In(i).Field(index).Type) != "map[string]interface {}" {
									panic("map only support map[string]interface{}")
								}
								bodyVal := make(map[string]interface{})
								if len(respBytes) > 0 {
									if err := json.Unmarshal(respBytes, &bodyVal); err != nil {
										panic(err)
									}
								}
								val.Set(reflect.ValueOf(bodyVal))
								break
							case reflect.Interface:
								var bodyVal interface{}
								if len(respBytes) > 0 {
									if err := json.Unmarshal(respBytes, &bodyVal); err != nil {
										panic(err)
									}
								}
								val.Set(reflect.ValueOf(bodyVal))
								break
							default:
								panic("Unsupported type , only support string , struct ,Slice ,  map[string]interface{} , interface{}")
							}
						} else {
							bodyVal := make(map[string]interface{})
							if len(respBytes) > 0 {
								if err := json.Unmarshal(respBytes, &bodyVal); err != nil {
									panic(err)
								}
							}
							value, err := bodyParamConverter(bodyParam, currentMethod.Type.In(i).Field(index).Type, bodyVal)
							if err != nil {
								panic(err)
							}
							val.Set(reflect.ValueOf(value))
						}

					}
				}
				inParam = append(inParam, destVal)
			}
			currentMethod.Func.Call(inParam) // 调用transport函数，传入参数
		}

	}
}

func bodyParamConverter(destName string, destType reflect.Type, parentVal map[string]interface{}) (interface{}, error) {
	value, ok := parentVal[destName]
	if !ok {
		return nil, fmt.Errorf("key %v not exist", destName)
	}
	switch destType.Kind() {
	case reflect.Int64:
		convertedValue, ok := value.(int64)
		if !ok {
			return nil, fmt.Errorf("key %v convert to int64 error ", destName)
		}
		return convertedValue, nil
	case reflect.Int32:
		convertedValue, ok := value.(int32)
		if !ok {
			return nil, fmt.Errorf("key %v convert to int32 error ", destName)
		}
		return convertedValue, nil
	case reflect.Int:
		convertedValue, ok := value.(int)
		if !ok {
			return nil, fmt.Errorf("key %v convert to int error ", destName)
		}
		return convertedValue, nil
	case reflect.String:
		convertedValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("key %v convert to string error ", destName)
		}
		return convertedValue, nil
	case reflect.Struct:
		c, _ := json.Marshal(value)
		targetVal := reflect.New(destType).Interface()
		decoder := json.NewDecoder(bytes.NewReader(c))
		if err := decoder.Decode(targetVal); err != nil {
			return nil, err
		}
		return reflect.Indirect(reflect.ValueOf(targetVal)).Interface(), nil
	default:
		return nil, fmt.Errorf("type %v not support", destType)
	}

}
