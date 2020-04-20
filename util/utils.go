package util

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RFC3339FullDate for rfc full date
const RFC3339FullDate = "2006-01-02"

// GetCurrentTime get current time
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetCurrentTimeString get current time with string
func GetCurrentTimeString(timeType string) string {
	return GetCurrentTime().Format(timeType)
}

// GetCurrentTimeUnix get current time with unix time
func GetCurrentTimeUnix() int64 {
	return GetCurrentTime().Unix()
}

// CheckFileIsExist : check file if exist ,exist -> true , not exist -> false  ,
/**
 * @param filename string ,the file name need to check
 * @return boolean string
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//StringInSlice if value in stringlist
func StringInSlice(value string, stringSlice []string) bool {
	for _, v := range stringSlice {
		if v == value {
			return true
		}
	}
	return false
}

//SliceInSlice if slice in slice
func SliceInSlice(sliceToCheck []string, slice []string) bool {
	for _, v := range sliceToCheck {
		if !StringInSlice(v, slice) {
			return false
		}
	}
	return true
}

//RecordErrorLevelTwo login error and print line , func , and error to gin context
func RecordErrorLevelTwo() (uintptr, string, int) {
	funcName, file, line, _ := runtime.Caller(2)
	return funcName, file, line
}

// Getparentdirectory : get parent directory of the path ,
/*
 * @param path string  ,the path you want to get parent directory
 * @return string  , the parent directory you need
 */
func Getparentdirectory(path string, level int) string {
	return strings.Join(strings.Split(path, "/")[0:len(strings.Split(path, "/"))-level], "/")
}

//DeleteExtraSpace remove extra space
func DeleteExtraSpace(s string) string {
	s1 := strings.Replace(s, "	", " ", -1)      //替换tab为空格
	regstr := "\\s{2,}"                         //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)            //编译正则表达式
	s2 := make([]byte, len(s1))                 //定义字符数组切片
	copy(s2, s1)                                //将字符串复制到切片
	spcIndex := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spcIndex) > 0 {                     //找到适配项
		s2 = append(s2[:spcIndex[0]+1], s2[spcIndex[1]:]...) //删除多余空格
		spcIndex = reg.FindStringIndex(string(s2))           //继续在字符串中搜索
	}
	return string(s2)
}

// GetTypeName to get struct type name
func GetTypeName(myvar interface{}, isToLowerCase bool) string {
	name := ""
	t := reflect.TypeOf(myvar)
	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}
	if isToLowerCase {
		name = strings.ToLower(name)
	}
	return name

}

// GenerateSnowFlakeID generate snowflake id
func GenerateSnowFlakeID(nodenumber int64) int64 {

	// Create a new Node with a Node number of 1
	node, _ := snowflake.NewNode(nodenumber)

	// Generate a snowflake ID.
	id := node.Generate().Int64()
	return id

}

// GetFreePort get one free port
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// GetServiceName get service name which will register to service mesh
func GetServiceName(projectName string) string {
	return fmt.Sprintf("%v", projectName)
}

// GetServiceID get service name which will register to service mesh
func GetServiceID(projectName string, projectVersion string, ServiceIP string, ServicePort int) string {
	ServiceName := GetServiceName(projectName)
	return fmt.Sprintf("%v-%v-%v-%v", ServiceName, projectVersion, ServiceIP, ServicePort)
}

// AddExtraSpaceIfExist adds a separator
func AddExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

// GRPCErrIsUnknownWrapper unknown err wrapper
// only will wrapper the error which is not be encoded
// P.S :this wrapper is used in grpc
func GRPCErrIsUnknownWrapper(err error) error {
	if err == nil {
		return nil
	}
	_, ok := status.FromError(err)
	if ok {
		return err
	}
	newErr := status.Error(codes.Internal, err.Error())
	return newErr
}

// HTTPErrEncoder  encode http err
// only will wrapper the error which is not be encoded
// P.S :this wrapper is used in grpc
func HTTPErrEncoder(err error) error {
	if err == nil {
		return nil
	}
	_, ok := status.FromError(err)
	if ok {
		return err
	}
	newErr := status.Error(codes.Unknown, err.Error())
	return newErr
}

// HTTPErrDecoder decode
func HTTPErrDecoder(err error) (bool, map[string]string) {
	if err == nil {
		return false, nil
	}
	status, ok := status.FromError(err)
	if ok {
		newErr := make(map[string]string)
		newErr["code"] = status.Code().String()
		newErr["message"] = status.Message()
		return true, newErr
	}
	return false, nil
}
