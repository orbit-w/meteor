package metadata

import (
	"context"
	"encoding/json"
	"strings"
)

/*
   @Author: orbit-w
   @File: metadata
   @2023 11月 周日 17:15
*/

type metaDataKey struct{}

type MD map[string]string

func (ins *MD) GetValue(key string) (v string, exist bool) {
	key = strings.ToLower(key)
	md := *ins
	v, exist = md[key]
	return
}

func NewMetaContext(father context.Context, m map[string]string) context.Context {
	md := MD{}
	for k, v := range m {
		md[strings.ToLower(k)] = v
	}
	return context.WithValue(father, metaDataKey{}, md)
}

func FromMetaContext(ctx context.Context) (md MD, ok bool) {
	md, ok = ctx.Value(metaDataKey{}).(MD)
	if !ok {
		return nil, false
	}
	return
}

func Marshal(m MD) ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Unmarshal(data []byte, dst *MD) error {
	return json.Unmarshal(data, dst)
}
