package extension

type ParseExtension interface {
	Apply()
}

type TagifyExtension interface {
	Apply()
}

type ScoreExtension interface {
	Apply()
}
