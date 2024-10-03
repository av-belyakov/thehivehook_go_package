package supportingfunctions

import (
	"fmt"
	"runtime"
)

// GetMemAlloc выводит информацию о затрачиваемой памяти приложения
func GetAppMemAlloc() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return fmt.Sprintf("%d KB\n", m.Alloc/1024)
}
