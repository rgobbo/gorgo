package gorgo

import (
	"regexp"
	"fmt"
	"github.com/spf13/cast"
	"strings"
	"strconv"
)

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var alphaNumericRegex = regexp.MustCompile("^[a-zA-Z0-9_]*$")
var numericRegexString  = regexp.MustCompile("^[-+]?[0-9]+(?:\\.[0-9]+)?$")
var numberRegex = regexp.MustCompile("^[0-9.]*$")

func isEmail(email string) bool {
	return emailRegexp.MatchString(email)
}

func isCpf(cpf string) bool {
	cpf = strings.Replace(cpf, ".", "", -1)
	cpf = strings.Replace(cpf, "-", "", -1)
	if len(cpf) != 11 {
		return false
	}
	var eq bool
	var dig string
	for _, val := range cpf {
		if len(dig) == 0 {
			dig = string(val)
		}
		if string(val) == dig {
			eq = true
			continue
		}
		eq = false
		break
	}
	if eq {
		return false
	}
	i := 10
	sum := 0
	for index := 0; index < len(cpf)-2; index++ {
		pos, _ := strconv.Atoi(string(cpf[index]))
		sum += pos * i
		i--
	}
	prod := sum * 10
	mod := prod % 11
	if mod == 10 {
		mod = 0
	}
	digit1, _ := strconv.Atoi(string(cpf[9]))
	if mod != digit1 {
		return false
	}
	i = 11
	sum = 0
	for index := 0; index < len(cpf)-1; index++ {
		pos, _ := strconv.Atoi(string(cpf[index]))
		sum += pos * i
		i--
	}
	prod = sum * 10
	mod = prod % 11
	if mod == 10 {
		mod = 0
	}
	digit2, _ := strconv.Atoi(string(cpf[10]))
	if mod != digit2 {
		return false
	}
	return true

}

func isCnpj(cnpj string) bool {
	cnpj = strings.Replace(cnpj, ".", "", -1)
	cnpj = strings.Replace(cnpj, "-", "", -1)
	cnpj = strings.Replace(cnpj, "/", "", -1)
	if len(cnpj) != 14 {
		return false
	}
	algs := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	var algProdCpfDig1 = make([]int, 12, 12)
	for key, val := range algs {
		intParsed, _ := strconv.Atoi(string(cnpj[key]))
		sumTmp := val * intParsed
		algProdCpfDig1[key] = sumTmp
	}
	sum := 0
	for _, val := range algProdCpfDig1 {
		sum += val
	}
	digit1 := sum % 11
	if digit1 < 2 {
		digit1 = 0
	} else {
		digit1 = 11 - digit1
	}
	char12, _ := strconv.Atoi(string(cnpj[12]))
	if char12 != digit1 {
		return false
	}
	algs = append([]int{6}, algs...)
	var algProdCpfDig2 = make([]int, 13, 13)
	for key, val := range algs {
		intParsed, _ := strconv.Atoi(string(cnpj[key]))
		sumTmp := val * intParsed
		algProdCpfDig2[key] = sumTmp
	}
	sum = 0
	for _, val := range algProdCpfDig2 {
		sum += val
	}
	digit2 := sum % 11
	if digit2 < 2 {
		digit2 = 0
	} else {
		digit2 = 11 - digit2
	}
	char13, _ := strconv.Atoi(string(cnpj[13]))
	if char13 != digit2 {
		return false
	}
	return true
}

func isAlphaNUmeric(str string) bool {
	return alphaNumericRegex.MatchString(str)
}

func isNumber(str string) bool {
 	return numericRegexString.MatchString(str)
}

func GetFunctions() FuncMap {
	funcs := FuncMap{
		"isemail" : isEmail,
		"iscpf" : isCpf,
		"iscnpj" : isCnpj,
		"isalphanumeric" : isAlphaNUmeric,
		"isnumber" : isNumber,
	}
	return funcs
}

func validateFields(collection string, data JSONDoc, mod *model, funcs FuncMap) error {
		if val, ok := mod.Tables[collection]; ok {
			for _, f := range val.Fields {

				if f.Validation != "" {
					fn := funcs[f.Validation].(func(string)bool)
					str := cast.ToString(data[f.Name])
					ret := fn(str)
					if ret == false {
						return fmt.Errorf("Validation [%s] error, field: %s - value: %v",f.Validation,f.Name, data[f.Name])
					}
				}
				if f.Required {
					if data[f.Name] == nil {
						return fmt.Errorf("Field [%s] required, receive %v",f.Name, data[f.Name])
					}
				}

				if f.Minlen > 0 {
					switch f.Type {
					case "string":
						str := cast.ToString(data[f.Name])
						if len(str) < f.Minlen {
							return fmt.Errorf("Field [%s] minlen %v and received %v",f.Name, f.Minlen, len(str))
						}
					case "int" :
						i := cast.ToInt(data[f.Name])
						if i < f.Minlen {
							return fmt.Errorf("Field [%s] minlen %v and received %v",f.Name, f.Minlen, i)
						}
					case "bigint" :
						i := cast.ToInt64(data[f.Name])
						if i < int64(f.Minlen) {
							return fmt.Errorf("Field [%s] minlen %v and received %v",f.Name, f.Minlen, i)
						}
					case "float" :
						fl := cast.ToFloat64(data[f.Name])
						if fl < float64(f.Minlen) {
							return fmt.Errorf("Field [%s] minlen %v and received %v",f.Name, f.Minlen, fl)
						}
					}
				}

				if f.Maxlen > 0 {
					switch f.Type {
					case "string":
						str := cast.ToString(data[f.Name])
						if len(str) > f.Maxlen {
							return fmt.Errorf("Field [%s] maxlen %v and received %v",f.Name, f.Maxlen, len(str))
						}
					case "int" :
						i := cast.ToInt(data[f.Name])
						if i > f.Maxlen {
							return fmt.Errorf("Field [%s] maxlen %v and received %v",f.Name, f.Maxlen, i)
						}
					case "bigint" :
						i := cast.ToInt64(data[f.Name])
						if i > int64(f.Maxlen) {
							return fmt.Errorf("Field [%s] maxlen %v and received %v",f.Name, f.Maxlen, i)
						}
					case "float" :
						fl := cast.ToFloat64(data[f.Name])
						if fl > float64(f.Maxlen) {
							return fmt.Errorf("Field [%s] maxlen %v and received %v",f.Name, f.Maxlen, fl)
						}
					}
				}

				if f.Default != nil && data[f.Name] == nil {
					data[f.Name] = f.Default
				}
				if f.Alias != "" {
					value := data[f.Name]
					delete(data, f.Name)
					data[f.Alias] = value
				}
			}

	}

	return nil
}
