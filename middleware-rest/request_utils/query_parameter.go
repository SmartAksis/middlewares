package request_utils

import (
	"fmt"
	"strconv"
	"strings"
)

type dataQueryStringConverterInterface interface {
	Convert(value interface{}) string
}

type int64Converter struct {next dataQueryStringConverterInterface
}
type int32Converter struct {next dataQueryStringConverterInterface
}
type int16Converter struct {next dataQueryStringConverterInterface
}
type int8Converter struct {next dataQueryStringConverterInterface
}
type intConverter struct {next dataQueryStringConverterInterface
}
type float64Converter struct {next dataQueryStringConverterInterface
}
type float32Converter struct {next dataQueryStringConverterInterface
}

func (converter int64Converter) Convert(value interface{}) string {
	result, instanceOf := value.(int64)
	if instanceOf {
		return fmt.Sprintf("%d", result)
	} else {
		var next dataQueryStringConverterInterface
		next=converter.next
		if next != nil {
			return next.Convert(value)
		}
	}
	return fmt.Sprintf("%s", result)
}

func (converter int32Converter) Convert(value interface{}) string {
	result, instanceOf := value.(int32)
	if instanceOf {
		return fmt.Sprintf("%d", result)
	} else {
		var next dataQueryStringConverterInterface
		next=converter.next
		if next != nil {
			return next.Convert(value)
		}
	}
	return fmt.Sprintf("%s", result)
}

func (converter int16Converter) Convert(value interface{}) string {
	result, instanceOf := value.(int16)
	if instanceOf {
		return fmt.Sprintf("%d", result)
	} else {
		var next dataQueryStringConverterInterface
		next=converter.next
		if next != nil {
			return next.Convert(value)
		}
	}
	return fmt.Sprintf("%s", result)
}

func (converter int8Converter) Convert(value interface{}) string {
	result, instanceOf := value.(int8)
	if instanceOf {
		return fmt.Sprintf("%d", result)
	} else {
		var next dataQueryStringConverterInterface
		next=converter.next
		if next != nil {
			return next.Convert(value)
		}
	}
	return fmt.Sprintf("%s", result)
}

func (converter intConverter) Convert(value interface{}) string {
	result, instanceOf := value.(int)
	if instanceOf {
		return fmt.Sprintf("%d", result)
	} else {
		var next dataQueryStringConverterInterface
		next=converter.next
		if next != nil {
			return next.Convert(value)
		}
	}
	return fmt.Sprintf("%s", result)
}

func (converter float64Converter) Convert(value interface{}) string {
	result, instanceOf := value.(float64)
	if instanceOf {
		stringValue:=strconv.FormatFloat(result, 'E', -1, 64)
		strings.Contains(stringValue, ".")
		if strings.Contains(stringValue, ".") {
			return fmt.Sprintf("%f", result)
		} else {
			return fmt.Sprintf("%d", int64(result))
		}
	} else {
		var next dataQueryStringConverterInterface
		next=converter.next
		if next != nil {
			return next.Convert(value)
		}
	}
	return fmt.Sprintf("%s", result)
}

func (converter float32Converter) Convert(value interface{}) string {
	result, instanceOf := value.(float32)
	if instanceOf {
		return fmt.Sprintf("%f", result)
	} else {
		var next dataQueryStringConverterInterface
		next=converter.next
		if next != nil {
			return next.Convert(value)
		}
	}
	return fmt.Sprintf("%s", result)
}



