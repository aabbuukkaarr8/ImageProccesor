package service

type Service struct {
	repo      Repository
	minio     StorageMinio
	producer  Producer
	processor Machine
}

func NewService(
	s StorageMinio,
	p Producer,
	proc Machine,
	r Repository,
) *Service {
	return &Service{
		minio:     s,
		producer:  p,
		processor: proc,
		repo:      r,
	}
}
