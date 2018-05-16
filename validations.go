package gorgo

import (
	"regexp"
	"fmt"
	"github.com/spf13/cast"
)

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var cpfRegex = regexp.MustCompile(`/^\d{3}\.\d{3}\.\d{3}\-\d{2}$/`)
var cnpjRegex = regexp.MustCompile(`/^\d{2}\.\d{3}\.\d{3}\/\d{4}\-\d{2}$/`)
var alphaNumericRegex = regexp.MustCompile("^[a-zA-Z0-9_]*$")
var numericRegexString  = regexp.MustCompile("^[-+]?[0-9]+(?:\\.[0-9]+)?$")
var numberRegex = regexp.MustCompile("^[0-9.]*$")

func isEmail(email string) bool {
	return emailRegexp.MatchString(email)
}

func isCpf(cpf string) bool {
	return cpfRegex.MatchString(cpf)
}

func isCnpj(cnpj string) bool {
	return cnpjRegex.MatchString(cnpj)
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
