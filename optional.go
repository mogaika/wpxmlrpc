package wp

type Optional map[string]interface{}

func NewOptional() Optional {
	return make(Optional)
}

func (o Optional) Set(name string, val interface{}) Optional {
	o[name] = val
	return o
}

func (o Optional) SetIf(cond bool, name string, val interface{}) Optional {
	if cond {
		o[name] = val
	}
	return o
}

func (o Optional) AddMap(m map[string]interface{}) Optional {
	if m != nil {
		for k, v := range m {
			o[k] = v
		}
	}
	return o
}
