package tian

type Info1 struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Newslist []struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Id       int    `json:"id"`
		TypeId   int    `json:"type_id"`
		TypeName string `json:"type_name"`
		CpName   string `json:"cp_name"`
		Zuofa    string `json:"zuofa"`
		Texing   string `json:"texing"`
		Tishi    string `json:"tishi"`
		Tiaoliao string `json:"tiaoliao"`
		Yuanliao string `json:"yuanliao"`
		En       string `json:"en"`
		Zh       string `json:"zh"`
		Saying   string `json:"saying"`
		Transl   string `json:"transl"`
		Source   string `json:"source"`
	} `json:"newslist"`
}
