package tool

import (
	"encoding/json"
	"strconv"
)

func Float64toA(a float64) string {
	return strconv.FormatFloat(a, 'f', 2, 64)
}
func AtoFloat64(a string) float64 {
	t, _ := strconv.ParseFloat(a, 64)
	return float64(t)
}

func AtoUin32(a string) uint32 {
	t, _ := strconv.ParseUint(a, 10, 32)
	return uint32(t)
}

func Atoint64(a string) int64 {
	t, _ := strconv.ParseInt(a, 10, 64)
	return int64(t)
}

func Atoint(a string) int {
	t, _ := strconv.ParseInt(a, 10, 64)
	return int(t)
}

func Uint32ToA(u uint32) string {
	return strconv.FormatUint(uint64(u), 10)
}
func Int64ToA(u int64) string {
	return strconv.FormatInt(int64(u), 10)
}

func StructToJson[T any](a T) string {
	res, _ := json.Marshal(a)
	return string(res)
}

func JsonToStruct[T any](as string) (T, error) {
	//var a map[string]interface{}
	a := new(T)
	err := json.Unmarshal([]byte(as), a)
	return *a, err
}
