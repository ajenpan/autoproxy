package pool

type Spider interface {
	Work() (string, error)
}

func Crawl() {

	sp := &HensonGetter{}

	result, err := sp.Work()
	if err != nil {
		return
	}

	for _, v := range result {

	}

	// result:= chan *CheckReport
	// ipChan := make(chan *models.IP, 2000)

}
