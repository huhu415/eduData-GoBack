package domain

type AddcouresStruct struct {
	LoginForm
	Color   string      `json:"color"`
	Coures  string      `json:"coures" binding:"required"`
	Teacher string      `json:"teacher"`
	Time    []TimeEntry `json:"time" binding:"required"`
}

type TimeEntry struct {
	Checkboxs  []int  `json:"checkboxs" binding:"required"`
	MultiIndex []int  `json:"multiIndex" binding:"required"`
	Place      string `json:"place"`
}
