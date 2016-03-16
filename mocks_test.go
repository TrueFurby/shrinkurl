package main

type mockStorage struct {
}

func NewMockStorage() *mockStorage {
	return &mockStorage{}
}

func (s *mockStorage) UrlStore() UrlStore {
	return s
}

func (s *mockStorage) Update(u *Url) error {
	return nil
}

func (s *mockStorage) Get(hash string) (u Url, err error) {
	return Url{}, nil
}

func (s *mockStorage) GetByUrl(url string) (u Url, err error) {
	return Url{}, nil
}

func (s *mockStorage) Remove(hash string) error {
	return nil
}
