package constant

type duration struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

var (
	Slots = map[uint]duration{
		1: {
			StartTime: "00:00",
			EndTime:   "01:00",
		},
		2: {
			StartTime: "01:00",
			EndTime:   "02:00",
		},
		3: {
			StartTime: "02:00",
			EndTime:   "03:00",
		},
		4: {
			StartTime: "03:00",
			EndTime:   "04:00",
		},
		5: {
			StartTime: "04:00",
			EndTime:   "05:00",
		},
		6: {
			StartTime: "05:00",
			EndTime:   "06:00",
		},
		7: {
			StartTime: "06:00",
			EndTime:   "07:00",
		},
		8: {
			StartTime: "07:00",
			EndTime:   "08:00",
		},
		9: {
			StartTime: "08:00",
			EndTime:   "09:00",
		},
		10: {
			StartTime: "09:00",
			EndTime:   "10:00",
		},
		11: {
			StartTime: "10:00",
			EndTime:   "11:00",
		},
		12: {
			StartTime: "11:00",
			EndTime:   "12:00",
		},
		13: {
			StartTime: "12:00",
			EndTime:   "13:00",
		},
		14: {
			StartTime: "13:00",
			EndTime:   "14:00",
		},
		15: {
			StartTime: "14:00",
			EndTime:   "15:00",
		},
		16: {
			StartTime: "15:00",
			EndTime:   "16:00",
		},
		17: {
			StartTime: "16:00",
			EndTime:   "17:00",
		},
		18: {
			StartTime: "17:00",
			EndTime:   "18:00",
		},
		19: {
			StartTime: "18:00",
			EndTime:   "19:00",
		},
		20: {
			StartTime: "19:00",
			EndTime:   "20:00",
		},
		21: {
			StartTime: "20:00",
			EndTime:   "21:00",
		},
		22: {
			StartTime: "21:00",
			EndTime:   "22:00",
		},
		23: {
			StartTime: "22:00",
			EndTime:   "23:00",
		},
		24: {
			StartTime: "23:00",
			EndTime:   "00:00",
		},
	}
)
