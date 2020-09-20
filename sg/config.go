package sg

import "fmt"

func (p *Plugin) GetParameters() map[string]string {
	return make(map[string]string)
}

func (p Plugin) SetParameter(name string, value string) error {
	return fmt.Errorf("security group deletion has no options")
}
