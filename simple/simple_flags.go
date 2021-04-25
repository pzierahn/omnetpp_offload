package simple

import "fmt"

type StringFlagList []string

func (list *StringFlagList) String() string {
	return fmt.Sprint(*list)
}

func (list *StringFlagList) Set(value string) error {
	*list = append(*list, value)
	return nil
}
