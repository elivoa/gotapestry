package register

import ()

// var TemplateKeyMap = &TempalteKeyMapStruct{
// 	Keymap: map[string]*TemplateUnit{},
// }

// type TempalteKeyMapStruct struct {
// 	l      sync.RWMutex
// 	Keymap map[string]*TemplateUnit // template key as key
// }

// structs

// type TemplateUnit struct {
// 	Key               string // template key.
// 	FilePath          string
// 	ContentOrigin     string `json:"-"`
// 	ContentTransfered string `json:"-"`
// 	IsCached          bool   `json:"-"`

// 	Blocks map[string]*BlockUnit

// 	// todo components?
// }

// Note: Component in blocks are directly belong to block's container.

// func (t *TemplateCache) Get(protonType reflect.Type) (*TemplateUnit, error) {
// 	t.l.RLock()
// 	defer t.l.RUnlock()
// 	if unit, ok := t.Templates[protonType]; ok {
// 		return unit, nil
// 	}
// 	return nil, errors.New("Template not loaded.")
// }

// // Get cached TemplateUnit by template key.
// func (t *TemplateCache) GetByKey(key string) (*TemplateUnit, error) {
// 	t.l.RLock()
// 	defer t.l.RUnlock()
// 	if unit, ok := t.Keymap[key]; ok {
// 		return unit, nil
// 	}
// 	return nil, errors.New("Template not loaded.")
// }

// func (t *TemplateCache) GetBlock(protonType reflect.Type, blockName string) (*BlockUnit, error) {
// 	t.l.RLock()
// 	defer t.l.RUnlock()
// 	if unit, ok := t.Templates[protonType]; ok {
// 		if nil == unit {
// 			return nil, errors.New("Error: Templates are nil, can't has blocks.")
// 		}
// 		if bu, okb := unit.Blocks[blockName]; okb {
// 			return bu, nil
// 		}
// 		return nil, errors.New(fmt.Sprintf("Block '%v' not found.", blockName))
// 	}
// 	return nil, errors.New("Template not loaded.")
// }

// func (t *TemplateCache) GetBlockByKey(key string, blockName string) (*BlockUnit, error) {
// 	t.l.RLock()
// 	defer t.l.RUnlock()
// 	if unit, ok := t.Keymap[key]; ok {
// 		if nil == unit {
// 			return nil, errors.New("Error: Templates are nil, can't has blocks.")
// 		}
// 		if bu, okb := unit.Blocks[blockName]; okb {
// 			return bu, nil
// 		}
// 		return nil, errors.New(fmt.Sprintf("Block '%v' not found.", blockName))
// 	}
// 	return nil, errors.New("Template not loaded.")
// }
