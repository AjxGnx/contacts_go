package dto

type Paginate struct {
	Page  int
	Limit int
}

func (p *Paginate) SetDefaultLimitAndPage() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.Limit == 0 {
		p.Limit = 10
	}
}
