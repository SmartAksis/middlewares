package request_utils

import (
	"fmt"
	"math"
)

type dataQueryStringConverterInterface interface {
	Convert(value interface{}) string
}
type stringConverter struct {next dataQueryStringConverterInterface
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

func (converter stringConverter) Convert(value interface{}) string {
	result, instanceOf := value.(string)
	if !instanceOf {
		var next dataQueryStringConverterInterface
		next=converter.next
		if next != nil {
			return next.Convert(value)
		}
	}
	return fmt.Sprintf("%s", result)
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
		if math.Mod(result, 1.0) == 0 {
			return fmt.Sprintf("%d", int64(result))
		} else {
			return fmt.Sprintf("%f", result)
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
		if math.Mod(float64(result), 1.0) == 0 {
			return fmt.Sprintf("%d", int64(result))
		} else {
			return fmt.Sprintf("%f", result)
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



