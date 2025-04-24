package persistence

const (
	defaultPage     = 1
	defaultPageSize = 10
)

func getValidPage(page int) int {
	if page < 1 {
		return defaultPage
	}

	return page
}

func getValidSize(size int) int {
	if size < 1 {
		return defaultPageSize
	}

	return size
}

func countOffset(page, size int) int {
	return (page - 1) * size
}

func getTotalPages(total, size int) int {
	totalPages := total / size

	if total%size > 0 {
		totalPages++
	}

	return totalPages
}
