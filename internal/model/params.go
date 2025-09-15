package model

type PaginationParams struct {
	IsPaginate bool
	Page       int
	PerPage    int
}

func (p *PaginationParams) Validate() {

	if p.Page < 1 {
		p.Page = 1
	}

	// default to 5
	if p.PerPage < 1 {
		p.PerPage = 5
	}

}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}
