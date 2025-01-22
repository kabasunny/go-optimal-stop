package optimization

// オプションを設定するための構造体と関数を定義
type resultOptions struct {
	ModelName   string
	SignalCount int
}

// resultOptions 構造体を変更するための関数の型定義
type ResultOption func(*resultOptions)

func WithModelName(name string) ResultOption {
	return func(opts *resultOptions) {
		opts.ModelName = name
	}
}

func WithSignalCount(count int) ResultOption {
	return func(opts *resultOptions) {
		opts.SignalCount = count
	}
}
