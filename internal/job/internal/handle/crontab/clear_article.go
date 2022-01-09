package crontab

type ClearArticle struct {
}

func NewClearArticle() *ClearArticle {
	return &ClearArticle{}
}

func (c *ClearArticle) Handle() error {
	return nil
}
