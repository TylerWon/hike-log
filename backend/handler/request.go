package handler

type createHikeRequest struct {
	TrailName     string  `json:"trailName" binding:"required"`
	Date          string  `json:"date" binding:"required,datetime=2006-01-02"`
	Notes         string  `json:"notes" binding:"required"`
	Rating        float32 `json:"rating" binding:"required,gte=0,lte=5,halfstep"`
	Difficulty    float32 `json:"difficulty" binding:"required,gte=0,lte=10,halfstep"`
	Distance      float32 `json:"distance" binding:"required"`
	ElevationGain uint    `json:"elevationGain" binding:"required"`
	Duration      uint    `json:"duration" binding:"required"`
	AllTrailsUrl  string  `json:"allTrailsUrl" binding:"required,url"`
}
