package obsidian

type Vault struct {
	Path string
}

func NewVault(path string) *Vault {
	return &Vault{Path: path}
}

// TODO: Реализовать функции для создания и редактирования Markdown-файлов
