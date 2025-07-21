package repository

import "github.com/jackc/pgx/v5/pgxpool"

// CalculatorRepo provides access to static pricing data for different class packs.
// Currently the data is hard-coded but could later be loaded from a database.
//
// Example JSON shape returned by the handler:
// {
//   "৬ষ্ঠ": {"ডায়মন্ড প্যাক":800, "গোল্ড প্যাক":600, "সিলভার প্যাক":400},
//   "৭ম": {...},
//   ...
// }
//
// Bengali words are preserved to perfectly match the mobile application copy.
// Consumers should treat keys as opaque strings.
//
// All values express gem amounts.
//
// As the data is read-only, no mutex is required.

type CalculatorRepo struct {
    DB *pgxpool.Pool // kept for future flexibility; currently unused
}

// NewCalculatorRepo constructs a CalculatorRepo. The db parameter can be nil.
func NewCalculatorRepo(db *pgxpool.Pool) *CalculatorRepo {
    return &CalculatorRepo{DB: db}
}

// pricingMap holds gems required for each package of a class/year.
var pricingMap = map[string]map[string]int{
    "৬ষ্ঠ": {
        "ডায়মন্ড প্যাক": 700, // 5% of 8000 * 1.75
        "গোল্ড প্যাক": 438,    // 5% of 5000 * 1.75
        "গোল্ড প্যাক (রেকর্ডেড)": 350, // 5% of 4000 * 1.75
    },
    "৭ম": {
        "ডায়মন্ড প্যাক": 700, // 5% of 8000 * 1.75
        "গোল্ড প্যাক": 438,    // 5% of 5000 * 1.75
        "গোল্ড প্যাক (রেকর্ডেড)": 350, // 5% of 4000 * 1.75
    },
    "৮ম": {
        "ডায়মন্ড প্যাক": 700, // 5% of 8000 * 1.75
        "গোল্ড প্যাক": 438,    // 5% of 5000 * 1.75
        "গোল্ড প্যাক (রেকর্ডেড)": 350, // 5% of 4000 * 1.75
    },
    "৯ম": {
        "ডায়মন্ড প্যাক": 875, // 5% of 10000 * 1.75
        "গোল্ড প্যাক": 613,    // 5% of 7000 * 1.75
        "গোল্ড প্যাক (রেকর্ডেড)": 350, // 5% of 4000 * 1.75
    },
    "১০ম": {
        "ডায়মন্ড প্যাক": 656, // 5% of 7500 * 1.75
        "গোল্ড প্যাক": 350,    // 5% of 4000 * 1.75
        "গোল্ড প্যাক (রেকর্ডেড)": 394, // 5% of 4500 * 1.75
    },
}

// GetPricing returns the entire pricing map.
func (r *CalculatorRepo) GetPricing() map[string]map[string]int {
    return pricingMap
}

// GetPrice fetches the gem count for a given class and package. The second return indicates presence.
func (r *CalculatorRepo) GetPrice(class, pack string) (int, bool) {
    if packs, ok := pricingMap[class]; ok {
        p, ok := packs[pack]
        return p, ok
    }
    return 0, false
}
