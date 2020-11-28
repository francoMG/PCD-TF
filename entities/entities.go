package entities

type Pacient struct {
    Malignant float64 `json:"malignant"`
    RadiusMea float64 `json:"radius_mea"`
    TextureMean   float64 `json:"texture_mean"`
    PerimeterMean float64 `json:"perimeter_mean"`
    AreaMean  float64 `json:"area_mean"`
    SmoothnessMea float64 `json:"smoothness_mea"`
    CompactnessMean   float64 `json:"compactness_mean"`
    ConcavityMean float64 `json:"concavity_mean"`
    ConcavePointsMea float64 `json:"concave_points_me"`
    SymmetryMean  float64 `json:"symmetry_mean"`
    FractalDimensionMean float64 `json:"fractal_dimension_mean"`
    RadiusSe  float64 `json:"radius_se"`
    TextureSe float64 `json:"texture_se"`
    PerimeterSe   float64 `json:"perimeter_se"`
    AreaS float64 `json:"area_s"`
    SmoothnessSe  float64 `json:"smoothness_se"`
    CompactnessSe float64 `json:"compactness_se"`
    ConcavitySe   float64 `json:"concavity_se"`
    ConcavePointsSe  float64 `json:"concave_points_se"`
    SymmetryS float64 `json:"symmetry_s"`
    FractalDimensionSe   float64 `json:"fractal_dimension_se"`
    RadiusWorst   float64 `json:"radius_worst"`
    TextureWorst  float64 `json:"texture_worst"`
    PerimeterWors float64 `json:"perimeter_wors"`
    AreaWorst float64 `json:"area_worst"`
    SmoothnessWorst   float64 `json:"smoothness_worst"`
    CompactnessWorst  float64 `json:"compactness_worst"`
    ConcavityWors float64 `json:"concavity_wors"`
    ConcavePointsWorst   float64 `json:"concave_points_worst"`
    SymmetryWorst float64 `json:"symmetry_worst"`
    FractalDimensionWors float64 `json:"fractal_dimension_wors"`
}