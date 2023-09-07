package main

type Params struct {
	// Minimum component absolute value
	ComponentMin uint
	// Maximum component absolute value
	ComponentMax uint
}

// SetDefault - Set default param values
func (p *Params) SetDefault() {
	p.ComponentMin = 0
	p.ComponentMax = 10
}

// Parse - Parse parameters from an ini file
// (only specified params in the file will change existing params)
func (p *Params) Parse(filepath string) (e error) {
	/*var file *os.File
	if file, e = os.Open(filepath); e != nil {
		return e
	}*/
	// TODO
	return nil
}
