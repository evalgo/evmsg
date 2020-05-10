package evmsg

func (m *Message) Value(key string) interface{} {
	err := CheckRequiredKeys(m, []string{key})
	if err != nil {
		return nil
	}
	for _, val := range m.Data.([]interface{}) {
		if value, ok := val.(map[string]interface{})[key]; ok {
			return value
		}
	}
	return nil
}

func (m *Message) Values() map[string]interface{} {
	if m.Data == nil {
		return map[string]interface{}{}
	}
	if len(m.Data.([]interface{})) == 0 {
		return map[string]interface{}{}
	}
	return m.Data.([]interface{})[0].(map[string]interface{})
}
