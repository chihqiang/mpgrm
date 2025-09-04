package httpx

// Paginate 分页处理工具函数
// 泛型 T 可以是任意类型
// fetchPage: 分页函数，接收页码，返回当前页数据列表和错误
// processItem: 对每个分页结果项的处理回调函数
// 返回值:
// - error: 如果分页函数返回错误，则中断并返回该错误
func Paginate[T any](fetchPage func(page int) ([]T, error), processItem func(T)) error {
	page := 1
	for {
		items, err := fetchPage(page) // 调用分页函数
		if err != nil {
			return err // 遇到错误立即返回
		}
		if len(items) == 0 {
			break // 数据为空表示分页结束
		}

		for _, item := range items {
			processItem(item) // 处理每条数据
		}

		page++
	}

	return nil // 正常结束返回 nil
}
