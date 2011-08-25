package pages
import(
	"os"
)
type WebContext interface {
	WriteString(content string)
	GetParams()map[string]string
	Write([]byte) (int, os.Error)
	WriteHeader(int)
}

