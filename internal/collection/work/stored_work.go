package work

type StoredWork struct {
	*Work
	Path       string
	AssetNames []string
}
