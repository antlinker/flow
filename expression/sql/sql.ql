query = fn(query,args...) {
	return sqlctx.Query(__ctx__,query,args...)
}
count = fn(query,args...) {
	return sqlctx.Count(__ctx__,query,args...)
}
one = fn(query,args...) {
	return sqlctx.One(__ctx__,query,args...)
}
querySliceStr = fn(query,key,args...) {
	return SliceStr(sqlctx.Query(__ctx__,query,args...),key)
}
export query,count,one