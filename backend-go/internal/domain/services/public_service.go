package services

import (
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/infrastructure/util"
)

type PublicService struct{}

func NewPublicService() *PublicService { return &PublicService{} }

// removed single GetRegionInfo in favor of batch endpoint

func (s *PublicService) GetRegions(regions []string) []struct {
	Region      string
	CountryCode int32
	Supported   bool
} {
	out := make([]struct {
		Region      string
		CountryCode int32
		Supported   bool
	}, 0, len(regions))
	for _, r := range regions {
		cc, ok := util.CountryCodeForRegion(r)
		out = append(out, struct {
			Region      string
			CountryCode int32
			Supported   bool
		}{Region: r, CountryCode: cc, Supported: ok})
	}
	return out
}
