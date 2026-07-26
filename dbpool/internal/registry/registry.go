package registry

var (
	stables []any
	gtables = make(map[string][]any)
)

func RegisterGlobal(dbname string, tables ...any) {
	gtables[dbname] = append(gtables[dbname], tables...)
}

func RegisterShards(tabs ...any) {
	stables = append(stables, tabs...)
}

func GetShardsTables() []any {
	return stables
}

func GetGlobalTables(dbname string) []any {
	return gtables[dbname]
}
