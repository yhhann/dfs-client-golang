package conf

import "flag"

var (
	flagMap map[string]*flag.Flag
)

func InitFlag() {
	if flagMap == nil {
		flagMap = make(map[string]*flag.Flag)

		if !flag.Parsed() {
			flag.Parse()
		}

		flag.VisitAll(func(f *flag.Flag) {
			flagMap[f.Name] = f
		})
	}
}

func Update(name string, value string) {
	if flagMap == nil {
		return
	}

	if f, ok := flagMap[name]; ok {
		f.Value.Set(value)
	}
}
