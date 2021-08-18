package stringif

func Substrings(p string)string{
	var s2 string
	for i:=8;i<len(p);i++{
		s2=s2+string(p[i])
	}
	return s2
}
