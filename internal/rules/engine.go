// Package rules implements the aircraft matching and scoring engine.
package rules

import (
	"math"
	"sort"
	"strings"

	"github.com/Maicarons/airlinesim-autobuy/internal/config"
	"github.com/Maicarons/airlinesim-autobuy/internal/marketdata"
	"github.com/Maicarons/airlinesim-autobuy/internal/parser"
)

// MatchResult represents the result of matching an aircraft against rules.
type MatchResult struct {
	Aircraft  *parser.AircraftOffer `json:"aircraft"`
	Rule      *config.RuleConfig    `json:"rule"`
	Score     float64               `json:"score"`
	ShouldBuy bool                  `json:"should_buy"`
}

// Engine evaluates aircraft offers against configured rules.
type Engine struct {
	rules []config.RuleConfig
}

// New creates a new rules engine with the given rules.
func New(rules []config.RuleConfig) *Engine {
	return &Engine{rules: rules}
}

// UpdateRules updates the rules used by the engine.
func (e *Engine) UpdateRules(rules []config.RuleConfig) {
	e.rules = rules
}

// Evaluate checks an aircraft against all enabled rules and returns matches.
func (e *Engine) Evaluate(aircraft *parser.AircraftOffer) []*MatchResult {
	var results []*MatchResult

	for i := range e.rules {
		rule := e.rules[i]
		if !rule.Enabled {
			continue
		}

		if match := e.matchRule(aircraft, &rule); match != nil {
			results = append(results, match)
		}
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// matchRule checks if an aircraft matches a single rule.
func (e *Engine) matchRule(aircraft *parser.AircraftOffer, rule *config.RuleConfig) *MatchResult {
	match := rule.Match

	// Check aircraft type by name list
	if len(match.Types) > 0 && !contains(match.Types, aircraft.Type) {
		return nil
	}

	// Check aircraft type by type_id (map filter ID to type name)
	if match.TypeID != "" {
		typeInfo := marketdata.GetTypeByID(match.TypeID)
		if typeInfo != nil && !strings.EqualFold(aircraft.Type, typeInfo.Name) {
			return nil
		}
	}

	// Check aircraft family by family_id (map filter ID to family name)
	if match.FamilyID != "" {
		familyInfo := marketdata.GetFamilyByID(match.FamilyID)
		if familyInfo != nil {
			// Match if the aircraft type contains the family name or vice versa
			// e.g., family "737-600/700/800/900" contains "737-800" which is in type "Boeing 737-800 HGW (winglets)"
			// Split family name by "/" and check if any part matches
			familyParts := strings.Split(familyInfo.Name, "/")
			familyMatch := false
			for _, part := range familyParts {
				part = strings.TrimSpace(part)
				if part != "" && (strings.Contains(aircraft.Type, part) || strings.Contains(part, aircraft.Type)) {
					familyMatch = true
					break
				}
			}
			// Also check individual words in the type name against the family name
			if !familyMatch {
				for _, word := range strings.Fields(aircraft.Type) {
					if strings.Contains(familyInfo.Name, word) {
						familyMatch = true
						break
					}
				}
			}
			if !familyMatch {
				return nil
			}
		}
	}

	// Determine the effective price based on the selected financing
	effectivePrice := aircraft.Price
	if len(match.Financing) == 1 {
		switch match.Financing[0] {
		case "lease":
			if aircraft.LeaseDep > 0 {
				effectivePrice = aircraft.LeaseDep
			} else if aircraft.LeaseRate > 0 {
				effectivePrice = aircraft.LeaseRate
			}
		case "credit":
			if aircraft.DownPmt > 0 {
				effectivePrice = aircraft.DownPmt
			} else if aircraft.Install > 0 {
				effectivePrice = aircraft.Install
			}
		}
	}

	// Check price range
	if effectivePrice > 0 {
		if match.PriceRange.Min > 0 && effectivePrice < match.PriceRange.Min {
			return nil
		}
		if match.PriceRange.Max > 0 && effectivePrice > match.PriceRange.Max {
			return nil
		}
	}

	// Check age
	if match.MaxAge > 0 && aircraft.Age > match.MaxAge {
		return nil
	}

	// Check cycles
	if match.MaxCycles > 0 && aircraft.Cycles > match.MaxCycles {
		return nil
	}

	// Check condition
	if match.ConditionMin > 0 && aircraft.Condition < match.ConditionMin {
		return nil
	}

	// Check offer type
	if len(match.OfferTypes) > 0 && aircraft.OfferType != "" {
		if !contains(match.OfferTypes, aircraft.OfferType) {
			return nil
		}
	}

	// Calculate score
	score := e.calculateScore(aircraft, rule)

	return &MatchResult{
		Aircraft:  aircraft,
		Rule:      rule,
		Score:     score,
		ShouldBuy: score > 0 && rule.Action.AutoBuy,
	}
}

// calculateScore computes a desirability score for the aircraft-rule pair.
// Higher score = better deal.
func (e *Engine) calculateScore(aircraft *parser.AircraftOffer, rule *config.RuleConfig) float64 {
	score := 0.0

	// Determine the effective price based on the selected financing
	effectivePrice := aircraft.Price
	if len(rule.Match.Financing) == 1 {
		switch rule.Match.Financing[0] {
		case "lease":
			if aircraft.LeaseDep > 0 {
				effectivePrice = aircraft.LeaseDep
			} else if aircraft.LeaseRate > 0 {
				effectivePrice = aircraft.LeaseRate
			}
		case "credit":
			if aircraft.DownPmt > 0 {
				effectivePrice = aircraft.DownPmt
			} else if aircraft.Install > 0 {
				effectivePrice = aircraft.Install
			}
		}
	}

	// Price score: how far below max price is this aircraft?
	// Lower price = higher score
	if rule.Match.PriceRange.Max > 0 && effectivePrice > 0 {
		priceRatio := 1.0 - (effectivePrice / rule.Match.PriceRange.Max)
		score += priceRatio * 50.0
	}

	// Condition score: better condition = higher score
	if rule.Match.ConditionMin > 0 && aircraft.Condition > 0 {
		condScore := (aircraft.Condition - rule.Match.ConditionMin) / (100.0 - rule.Match.ConditionMin)
		score += math.Max(0, condScore) * 20.0
	}

	// Age score: younger = better
	if rule.Match.MaxAge > 0 && aircraft.Age >= 0 {
		ageScore := 1.0 - (float64(aircraft.Age) / float64(rule.Match.MaxAge))
		score += math.Max(0, ageScore) * 20.0
	}

	// Immediate purchase bonus
	if aircraft.OfferType == "immediate" {
		score += 10.0
	}

	// Priority bonus from rule
	score += float64(rule.Priority)

	return score
}

// BestMatch returns the single best match for an aircraft.
func (e *Engine) BestMatch(aircraft *parser.AircraftOffer) *MatchResult {
	results := e.Evaluate(aircraft)
	if len(results) == 0 {
		return nil
	}
	return results[0]
}

// contains checks if a string slice contains a value (case-insensitive).
func contains(slice []string, value string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, value) {
			return true
		}
	}
	return false
}

// SetRules updates the rules in the engine.
func (e *Engine) SetRules(rules []config.RuleConfig) {
	e.rules = rules
}

// GetRules returns the current rules.
func (e *Engine) GetRules() []config.RuleConfig {
	return e.rules
}