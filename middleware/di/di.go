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

// stringConverter
/*
	@word the word need to be converted
	@destVal the target value need to convert to

	stringConverter will convert to destVal according to
	the type of destVal
*/
func stringConverter(word string, destVal *reflect.Value) error {
	switch destVal.Type().Kind() {
	case reflect.Int64:
		paramVal, err := strconv.ParseInt(word, 10, 64)
		if err != nil {
			return err
		}
		destVal.Set(reflect.ValueOf(paramVal))
		break
	case reflect.Int32:
		paramVal, err := strconv.ParseInt(word, 10, 32)
		if err != nil {
			return err
		}
		destVal.Set(reflect.ValueOf(paramVal))
		break
	case reflect.Int16:
		paramVal, err := strconv.ParseInt(word, 10, 16)
		if err != nil {
			return err
		}
		destVal.Set(reflect.ValueOf(paramVal))
		break
	case reflect.Int:
		paramVal, err := strconv.Atoi(word)
		if err != nil {
			return err
		}
		destVal.Set(reflect.ValueOf(paramVal))
		break
	case reflect.Bool:
		paramVal, err := strconv.ParseBool(word)
		if err != nil {
			return err
		}
		destVal.Set(reflect.ValueOf(paramVal))
		break
	case reflect.String:
		destVal.Set(reflect.ValueOf(word))
		break
	default:
		return fmt.Errorf("Unsupported type")
	}
	return nil
}

// bodyParamConverter
/*
	@bodyVal the father val of
	@key the key name of the parentVal
	@destType the type of dest type
	bodyParamConverter will get the key from the body value
	and convert the value to the dest type value
*/
func bodyParamConverter(bodyVal map[string]interface{}, key string, destType reflect.Type) (interface{}, error) {
	value, ok := bodyVal[key]
	if !ok {
		return nil, fmt.Errorf("key %v not exist", key)
	}
	switch destType.Kind() {
	case reflect.Int64:
		convertedValue, ok := value.(int64)
		if !ok {
			return nil, fmt.Errorf("key %v convert to int64 error ", key)
		}
		return convertedValue, nil
	case reflect.Int32:
		convertedValue, ok := value.(int32)
		if !ok {
			return nil, fmt.Errorf("key %v convert to int32 error ", key)
		}
		return convertedValue, nil
	case reflect.Int:
		convertedValue, ok := value.(int)
		if !ok {
			return nil, fmt.Errorf("key %v convert to int error ", key)
		}
		return convertedValue, nil
	case reflect.String:
		convertedValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("key %v convert to string error ", key)
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

// New DI middleware
func New(app application.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := fmt.Sprintf("%v@%v", c.Request.Method, c.FullPath())
		runtimeKeyMap := application.DecodeHTTPRuntimeKey(c, app.RuntimeKeys())
		tContext := app.ContextPool().Acquire(app, runtimeKeyMap, c)
		defer func() {
			if tContext.AutoFree() {
				//release trinity go context obj
				app.ContextPool().Release(tContext)
			}
		}()
		controller, _, toFreeInstance := app.ControllerPool().GetController(method, tContext, app, c)
		defer func() {
			for _, v := range toFreeInstance {
				app.InstancePool().Release(v)
			}
		}()
		validators := app.ControllerPool().GetControllerValidators(method)
		for _, v := range validators {
			if err := v(tContext); err != nil {
				tContext.HTTPResponseInternalErr(err)
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
		var responseValue []reflect.Value
		if currentMethod.Type.NumIn() < -1 {
			panic("wrong method")
		} else if currentMethod.Type.NumIn() == 1 {
			var inParam []reflect.Value                            // 构造函数入参 ， 入参1 ， transport指针对象
			inParam = append(inParam, reflect.ValueOf(controller)) // 传入transport对象
			responseValue = currentMethod.Func.Call(inParam)       // 调用transport函数，传入参数
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
						if err := stringConverter(headerValString, &val); err != nil {
							panic(err)
						}
					}
					// check if path param
					if pathParam, isExist := currentMethod.Type.In(i).Field(index).Tag.Lookup("path_param"); isExist {
						paramValString, ok := c.Params.Get(pathParam)
						if !ok {
							panic(fmt.Sprintf("%v param not exist ", pathParam))
						}
						if err := stringConverter(paramValString, &val); err != nil {
							panic(err)
						}
					}
					// check if query param
					if queryParam, isExist := currentMethod.Type.In(i).Field(index).Tag.Lookup("query_param"); isExist {
						if queryParam == "" {
							switch val.Type().Kind() {
							case reflect.String:
								val.Set(reflect.ValueOf(c.Request.URL.RawQuery))
								break
							default:
								panic("Unsupported type , only support string ")
							}
						} else {
							queryValString := c.Query(queryParam)
							if err := stringConverter(queryValString, &val); err != nil {
								panic(err)
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
							switch val.Type().Kind() {
							case reflect.String:
								val.Set(reflect.ValueOf(string(respBytes)))
								break
							case reflect.Struct, reflect.Slice:
								// if is []byte
								if fmt.Sprintf("%v", currentMethod.Type.In(i).Field(index).Type) == "[]uint8" {
									val.Set(reflect.Indirect(reflect.ValueOf(respBytes)))
									break
								} else {
									targetVal := reflect.New(currentMethod.Type.In(i).Field(index).Type).Interface()
									if err := c.BindJSON(targetVal); err != nil {
										panic(err)
									}
									val.Set(reflect.Indirect(reflect.ValueOf(targetVal)))
									break
								}

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
								panic("Unsupported type , only support string , struct ,Slice ,  map[string]interface{} , interface{} , []byte")
							}
						} else {
							bodyVal := make(map[string]interface{})
							if len(respBytes) > 0 {
								if err := json.Unmarshal(respBytes, &bodyVal); err != nil {
									panic(err)
								}
							}
							value, err := bodyParamConverter(bodyVal, bodyParam, currentMethod.Type.In(i).Field(index).Type)
							if err != nil {
								panic(err)
							}
							val.Set(reflect.ValueOf(value))
						}

					}
				}
				inParam = append(inParam, destVal)
			}
			responseValue = currentMethod.Func.Call(inParam) // 调用transport函数，传入参数
		}
		switch len(responseValue) {
		case 0:
			return
		case 1:
			if err, ok := responseValue[0].Interface().(error); ok {
				if err != nil {
					tContext.HTTPResponseInternalErr(err)
					return
				}
			}
			tContext.HTTPResponse(responseValue[0].Interface(), nil)
			return
		case 2:
			if err, ok := responseValue[1].Interface().(error); ok {
				if err != nil {
					tContext.HTTPResponseInternalErr(err)
					return
				}
			}
			tContext.HTTPResponse(responseValue[0].Interface(), nil)
			return
		default:
			panic("wrong res type , first out should be response value , second out should be error ")
		}

	}
}
