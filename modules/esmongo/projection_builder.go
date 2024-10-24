package esmongo

type Projection int

const (
	ProjectionInclude = 1
	ProjectionExclude = 0
)

type ProjectionBuilder map[string]Projection

func (pb ProjectionBuilder) Include(key string) ProjectionBuilder {
	pb[key] = ProjectionInclude
	return pb
}

func (pb ProjectionBuilder) Exclude(key string) ProjectionBuilder {
	pb[key] = ProjectionExclude
	return pb
}

func (pb ProjectionBuilder) Clear() ProjectionBuilder {
	clear(pb)
	return pb
}
