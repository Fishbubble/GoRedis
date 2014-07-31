package levelredis

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	WrongKindError   = errors.New("wrong kind error")
	BadArgumentCount = errors.New("bad argument count")
	BadArgumentType  = errors.New("bad argument type")
	msitype          = reflect.TypeOf(make(map[string]interface{}))
)

const (
	dot = "."
)

// 提供面向document操作的map
// doc := New()
// doc.Set(jsonObj)
// doc.Get(fields)
type MapDoc struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewMapDoc(data map[string]interface{}) (m *MapDoc) {
	m = &MapDoc{}
	if m.data = data; m.data == nil {
		m.data = make(map[string]interface{})
	}
	return
}

// doc_set(key, {"name":"latermoon", "$rpush":["photos", "c.jpg", "d.jpg"], "$incr":["version", 1]})
func (m *MapDoc) Set(in map[string]interface{}) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// defer func() {
	// 	if v := recover(); v != nil {
	// 		if e, ok := v.(error); ok {
	// 			err = e
	// 		} else {
	// 			err = errors.New(fmt.Sprint(v))
	// 		}
	// 	}
	// }()
	for k, v := range in {
		if !strings.HasPrefix(k, "$") {
			parent, key, val, exist := m.findEntry(k, false)
			fmt.Println(parent, key, val, exist)
			// 检查是否可覆盖
			if exist {
				// 如果原数据是复合结构，则新数据必须同类型
				if m.isComplexType(val) && reflect.TypeOf(val).Kind() != reflect.TypeOf(v).Kind() {
					return errors.New("bad value type for `" + k + "`")
				}
				// 如果原数据是简单数据类型，则新数据不能是复杂结构
				if !m.isComplexType(val) && m.isComplexType(v) {
					return errors.New("bad value type for `" + k + "`")
				}
			} else {
				// 创建父对象
				parent, key, _, _ = m.findEntry(k, true)
			}
			parent[key] = v
			continue
		}
		action := k[1:]
		switch action {
		case "set":
			argmap := v.(map[string]interface{})
			for field, value := range argmap {
				parent, key, _, _ := m.findElement(field, true)
				parent[key] = value
			}
		case "rpush":
			argmap, ok := v.(map[string]interface{})
			if !ok {
				return errors.New("bad rpush format")
			}
			for field, value := range argmap {
				parent, key, _, _ := m.findElement(field, true)
				if _, ok := value.([]interface{}); !ok {
					return errors.New("bad rpush items format")
				}
				m.doRpush(parent, key, value.([]interface{}))
			}
		case "incr":
			argmap := v.(map[string]interface{})
			for field, value := range argmap {
				parent, key, _, _ := m.findElement(field, true)
				err = m.doIncr(parent, key, value)
			}
		case "del":
			arglist := v.([]interface{})
			for _, field := range arglist {
				parent, key, _, exist := m.findElement(field.(string), false)
				if exist {
					delete(parent, key)
				}
			}
		default:
		}
	}
	return
}

// 是否json中的数组或对象
func (m *MapDoc) isComplexType(val interface{}) (ok bool) {
	if val == nil {
		return false
	}
	kind := reflect.TypeOf(val).Kind()
	return kind == reflect.Slice || kind == reflect.Map
}

// doc_get(key, ["name", "setting.mute", "photos.$1"])
func (m *MapDoc) Get(fields ...string) (out map[string]interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out = make(map[string]interface{})
	if len(fields) == 0 || (len(fields) == 1 && fields[0] == "") {
		for k, v := range m.data {
			out[k] = v
		}
		return
	}

	for _, field := range fields {
		dst := out
		src := m.data
		// 逐个字段扫描copy
		pairs := strings.Split(field, dot)
		count := len(pairs)
		for i := 0; i < count; i++ {
			curkey := pairs[i]
			obj, ok := src[curkey]
			if !ok {
				break
			}
			// 定位到最后一个元素
			if i == count-1 {
				dst[curkey] = obj
				break
			}
			// 没到最后一个元素，但原始数据不是map[string]interface{}，表示出错
			if reflect.TypeOf(obj) != msitype {
				delete(dst, curkey)
				break
			}
			if dst[curkey] == nil {
				dst[curkey] = make(map[string]interface{})
			}
			src = src[curkey].(map[string]interface{})
			dst = dst[curkey].(map[string]interface{})
		}
	}
	return
}

func (m *MapDoc) findEntry(field string, create bool) (parent map[string]interface{}, key string, val interface{}, exist bool) {
	pairs := strings.Split(field, dot)
	parent = m.data
	for i := 0; i < len(pairs)-1; i++ {
		curkey := pairs[i]
		obj, ok := parent[curkey]
		if !ok || reflect.TypeOf(obj) != msitype {
			if create {
				parent[curkey] = make(map[string]interface{})
			} else {
				return nil, "", nil, false // bye
			}
		}
		parent = parent[curkey].(map[string]interface{})
	}
	exist = true
	// 定位最后的元素
	key = pairs[len(pairs)-1]
	val = parent[key]
	return
}

/**
 * 根据field路径查找元素
 * @param field 多级的field使用"."分隔
 * @return parent[key] == obj，其中 parent 目标元素父对象，必定是map[string]interface{}，key 目标元素key，obj，目标元素
 */
func (m *MapDoc) findElement(field string, createIfMissing bool) (parent map[string]interface{}, key string, obj interface{}, exist bool) {
	pairs := strings.Split(field, dot)
	parent = m.data
	for i := 0; i < len(pairs)-1; i++ {
		curkey := pairs[i]
		var ok bool
		_, ok = parent[curkey]
		// 初始化或覆盖
		if !ok || reflect.TypeOf(parent[curkey]) != msitype {
			if createIfMissing {
				parent[curkey] = make(map[string]interface{})
			} else {
				exist = false
				return
			}
		}
		parent = parent[curkey].(map[string]interface{})
	}
	exist = true
	key = pairs[len(pairs)-1]
	obj = parent[key]
	return
}

func (m *MapDoc) doRpush(parent map[string]interface{}, key string, elems []interface{}) (err error) {
	obj := parent[key]
	if obj != nil {
		for i := 0; i < len(elems); i++ {
			parent[key] = append(parent[key].([]interface{}), elems[i])
		}
	} else {
		parent[key] = elems
	}
	return
}

func (m *MapDoc) doIncr(parent map[string]interface{}, key string, value interface{}) (err error) {
	obj := parent[key]
	if obj == nil {
		parent[key] = value
	} else {
		oldint, e1 := toInt(obj)
		if e1 != nil {
			return errors.New(fmt.Sprintf("`%s` is not `int`", key))
		}

		incrint, e2 := toInt(value)
		if e2 != nil {
			return errors.New(fmt.Sprintf("incrment of `%s` is not `int`", key))
		}
		parent[key] = oldint + incrint
	}
	return
}

func toInt(obj interface{}) (n int, err error) {
	switch obj.(type) {
	case int:
		n = obj.(int)
	case float64:
		n = int(obj.(float64))
	default:
		err = BadArgumentType
	}
	return
}

func (m *MapDoc) String() string {
	b, _ := json.Marshal(m.data)
	return string(b)
}

func (m *MapDoc) Map() map[string]interface{} {
	return m.data
}
