package wit

import "errors"

// Differentiate sets the import or export direction of each
// [Interface], [TypeDef], and [Function] in a [Resolve],
// duplicating each [WorldItem] as necessary to differentiate
// between imported and exported items.
func (r *Resolve) Differentiate() error {
	// Differentiate all imports first
	for _, w := range r.Worlds {
		w.Imports.All()(func(_ string, i WorldItem) bool {
			switch i := i.(type) {
			case *Interface:
				r.importInterface(i)
			case *TypeDef:
				r.importTypeDef(i)
			case *Function:
				r.importFunction(i)
			}
			return true
		})
	}

	// Then exports
	for _, w := range r.Worlds {
		var err error
		w.Exports.All()(func(name string, i WorldItem) bool {
			switch i := i.(type) {
			case *Interface:
				i = r.exportInterface(i)
				w.Exports.Set(name, i)
			case *TypeDef:
				err = errors.New("exported type in world " + w.Name)
				return false
			case *Function:
				i = r.exportFunction(i)
				w.Exports.Set(name, i)
			}
			return true
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolve) importInterface(i *Interface) {
	i.Direction = Imported
	i.TypeDefs.All()(func(_ string, t *TypeDef) bool {
		r.importTypeDef(t)
		return true
	})
	i.Functions.All()(func(_ string, f *Function) bool {
		r.importFunction(f)
		return true
	})
}

func (r *Resolve) importTypeDef(t *TypeDef) {
	t.Direction = Imported
}

func (r *Resolve) importFunction(f *Function) {
	f.Direction = Imported
}

func (r *Resolve) exportInterface(i *Interface) *Interface {
	if i.Direction == Imported {
		clone := *i
		i = &clone
		r.Interfaces = append(r.Interfaces, i)
	}
	i.Direction = Exported
	i.TypeDefs.All()(func(name string, t *TypeDef) bool {
		t = r.exportTypeDef(t)
		i.TypeDefs.Set(name, t)
		return true
	})
	i.Functions.All()(func(name string, f *Function) bool {
		f = r.exportFunction(f)
		i.Functions.Set(name, f)
		return true
	})
	return i
}

func (r *Resolve) exportTypeDef(t *TypeDef) *TypeDef {
	if t.Direction == Imported {
		clone := *t
		t = &clone
		r.TypeDefs = append(r.TypeDefs, t)
	}
	t.Direction = Exported
	return t
}

func (r *Resolve) exportFunction(f *Function) *Function {
	if f.Direction == Imported {
		clone := *f
		f = &clone
	}
	f.Direction = Exported
	return f
}
