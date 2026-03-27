package functions

func registerAll(r *Registry) {
	registerStringFuncs(r)
	registerNumberFuncs(r)
	registerLogicFuncs(r)
	registerListFuncs(r)
	registerMapFuncs(r)
	registerTypeFuncs(r)
	registerMiscFuncs(r)
}
