package main

type mockStorage struct {
	urls map[string]Url
}

func NewMockStorage() *mockStorage {
	return &mockStorage{make(map[string]Url)}
}

func (s *mockStorage) UrlStore() UrlStore {
	return s
}

func (s *mockStorage) Update(u *Url) error {
	s.urls[u.Hash] = *u
	return nil
}

func (s *mockStorage) Get(hash string) (Url, error) {
	u := s.urls[hash]
	return u, nil
}

func (s *mockStorage) GetByUrl(url string) (Url, error) {
	for _, u := range s.urls {
		if u.Url == url {
			return u, nil
		}
	}
	return Url{}, nil
}

func (s *mockStorage) Remove(hash string) error {
	delete(s.urls, hash)
	return nil
}
