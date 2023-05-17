package llmchain

var (
	chainsMap map[string]Chain
)

// Reg reg Chain for later use.
func RegChain(c Chain) {

	if chainsMap == nil {
		chainsMap = make(map[string]Chain)
	}
	//same name will be repalced by later reg oper
	chainsMap[c.Name()] = c

}

// Get return Chain by name
func GetChain(name string) (Chain, bool) {

	c, ok := chainsMap[name]

	// if !ok { //TODO
	// 	return &BaseChat{}, true
	// }

	return c, ok

}
