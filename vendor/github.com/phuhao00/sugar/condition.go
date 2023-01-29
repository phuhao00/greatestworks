package sugar

func Ternary[T any](condition bool, ifOutput T, elseOutput T) T {

	if condition {
		return ifOutput
	}

	return elseOutput
}

type IfElse[T any] struct {
	Result T
	Ok     bool
}

func If[T any](condition bool, result T) *IfElse[T] {

	if condition {
		return &IfElse[T]{result, true}
	}

	var t T
	return &IfElse[T]{t, false}
}

func IfFn[T any](condition bool, resultFn func() T) *IfElse[T] {

	if condition {
		return &IfElse[T]{resultFn(), true}
	}

	var t T
	return &IfElse[T]{t, false}
}

func (i *IfElse[T]) ElseIf(condition bool, result T) *IfElse[T] {

	if condition && !i.Ok {
		i.Result = result
		i.Ok = true
	}

	return i
}

func (i *IfElse[T]) ElseIfFn(condition bool, resultFn func() T) *IfElse[T] {
	if condition && !i.Ok {
		i.Result = resultFn()
		i.Ok = true
	}

	return i
}

func (i *IfElse[T]) Else(result T) T {
	if i.Ok {
		return i.Result
	}

	return result
}

func (i *IfElse[T]) ElseFn(resultFn func() T) T {
	if i.Ok {
		return i.Result
	}

	return resultFn()
}

type SwitchCase[T comparable, R any] struct {
	Predicate T
	Result    R
	Ok        bool
}

func Switch[T comparable, R any](predicate T) *SwitchCase[T, R] {
	var result R

	return &SwitchCase[T, R]{
		predicate,
		result,
		false,
	}
}

func (s *SwitchCase[T, R]) Case(value T, result R) *SwitchCase[T, R] {
	if !s.Ok && s.Predicate == value {
		s.Result = result
		s.Ok = true
	}
	return s
}

func (s *SwitchCase[T, R]) CaseFn(value T, resultFn func() R) *SwitchCase[T, R] {
	if !s.Ok && value == s.Predicate {
		s.Result = resultFn()
		s.Ok = true
	}

	return s
}

func (s *SwitchCase[T, R]) Default(result R) R {
	if !s.Ok {
		s.Result = result
	}

	return s.Result
}

func (s *SwitchCase[T, R]) DefaultFn(resultFn func() R) R {
	if !s.Ok {
		s.Result = resultFn()
	}

	return s.Result
}
