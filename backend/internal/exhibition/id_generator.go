package exhibition

// IDGenerator は、ExhibitionのIDを生成する役割を表す。
//
// Exhibition側は、
// UUIDなのか連番なのかランダム文字列なのかを知らない。
//
// 「IDを生成できる」ことだけ知っていればよい。
type IDGenerator interface {
	Generate() string
}
