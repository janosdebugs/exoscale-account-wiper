package nlbs

import "fmt"

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p Plugin) SetParameter(name string, value string) error {
	return fmt.Errorf("NLB deletion has no options")
}
