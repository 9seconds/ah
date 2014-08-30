package homestorage


const (
	HOMEDIR_NAME = ".ah"
)

type HomeStorage struct {
	content map[string]string
	homeDirPath string
	parsed bool
}


func (hs *HomeStorage) Init(root string) {
    hs.homeDirPath = path.Join(root, HODEDIR_NAME)
	hs.parsed = false
	hs.content = make(map[string]string)
}

