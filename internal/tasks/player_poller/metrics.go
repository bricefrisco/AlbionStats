package player_poller

import (
	"albionstats/internal/sqlite"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
)

type VMMetric struct {
	Metric     map[string]string `json:"metric"`
	Values     []float64         `json:"values"`
	Timestamps []int64           `json:"timestamps"`
}

func PushToVictoriaMetrics(stats []sqlite.PlayerStats, logger *slog.Logger) error {
	if len(stats) == 0 {
		return nil
	}

	batch := make([]VMMetric, 0)
	for _, stat := range stats {
		metrics := playerStatsToMetrics(stat)
		batch = append(batch, metrics...)
	}

	body, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("marshal batch: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8428/api/v1/import", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	logger.Info("pushed to victoria metrics", "num_metrics", len(batch), "status", resp.StatusCode)

	return nil
}

func playerStatsToMetrics(stats sqlite.PlayerStats) []VMMetric {
	labels := map[string]string{
		"region":    stats.Region,
		"player_id": stats.PlayerID,
	}

	ts := stats.TS.UnixMilli()

	metrics := []VMMetric{
		m("albion_player_kill_fame", labels, stats.KillFame, ts),
		m("albion_player_death_fame", labels, stats.DeathFame, ts),

		m("albion_player_pve_total", labels, stats.PveTotal, ts),
		m("albion_player_pve_royal", labels, stats.PveRoyal, ts),
		m("albion_player_pve_outlands", labels, stats.PveOutlands, ts),
		m("albion_player_pve_avalon", labels, stats.PveAvalon, ts),
		m("albion_player_pve_hellgate", labels, stats.PveHellgate, ts),
		m("albion_player_pve_corrupted", labels, stats.PveCorrupted, ts),
		m("albion_player_pve_mists", labels, stats.PveMists, ts),

		m("albion_player_gather_fiber_total", labels, stats.GatherFiberTotal, ts),
		m("albion_player_gather_fiber_royal", labels, stats.GatherFiberRoyal, ts),
		m("albion_player_gather_fiber_outlands", labels, stats.GatherFiberOutlands, ts),
		m("albion_player_gather_fiber_avalon", labels, stats.GatherFiberAvalon, ts),

		m("albion_player_gather_hide_total", labels, stats.GatherHideTotal, ts),
		m("albion_player_gather_hide_royal", labels, stats.GatherHideRoyal, ts),
		m("albion_player_gather_hide_outlands", labels, stats.GatherHideOutlands, ts),
		m("albion_player_gather_hide_avalon", labels, stats.GatherHideAvalon, ts),

		m("albion_player_gather_ore_total", labels, stats.GatherOreTotal, ts),
		m("albion_player_gather_ore_royal", labels, stats.GatherOreRoyal, ts),
		m("albion_player_gather_ore_outlands", labels, stats.GatherOreOutlands, ts),
		m("albion_player_gather_ore_avalon", labels, stats.GatherOreAvalon, ts),

		m("albion_player_gather_rock_total", labels, stats.GatherRockTotal, ts),
		m("albion_player_gather_rock_royal", labels, stats.GatherRockRoyal, ts),
		m("albion_player_gather_rock_outlands", labels, stats.GatherRockOutlands, ts),
		m("albion_player_gather_rock_avalon", labels, stats.GatherRockAvalon, ts),

		m("albion_player_gather_wood_total", labels, stats.GatherWoodTotal, ts),
		m("albion_player_gather_wood_royal", labels, stats.GatherWoodRoyal, ts),
		m("albion_player_gather_wood_outlands", labels, stats.GatherWoodOutlands, ts),
		m("albion_player_gather_wood_avalon", labels, stats.GatherWoodAvalon, ts),

		m("albion_player_gather_all_total", labels, stats.GatherAllTotal, ts),
		m("albion_player_gather_all_royal", labels, stats.GatherAllRoyal, ts),
		m("albion_player_gather_all_outlands", labels, stats.GatherAllOutlands, ts),
		m("albion_player_gather_all_avalon", labels, stats.GatherAllAvalon, ts),

		m("albion_player_crafting_total", labels, stats.CraftingTotal, ts),
		m("albion_player_crafting_royal", labels, stats.CraftingRoyal, ts),
		m("albion_player_crafting_outlands", labels, stats.CraftingOutlands, ts),
		m("albion_player_crafting_avalon", labels, stats.CraftingAvalon, ts),

		m("albion_player_fishing_fame", labels, stats.FishingFame, ts),
		m("albion_player_farming_fame", labels, stats.FarmingFame, ts),
		m("albion_player_crystal_league_fame", labels, stats.CrystalLeagueFame, ts),
	}

	if m, ok := m_float("albion_player_fame_ratio", labels, stats.FameRatio, ts); ok {
		metrics = append(metrics, m)
	}

	return metrics
}

func m(name string, labels map[string]string, val int64, ts int64) VMMetric {
	l := make(map[string]string, len(labels)+1)
	for k, v := range labels {
		l[k] = v
	}
	l["__name__"] = name

	return VMMetric{
		Metric:     l,
		Values:     []float64{float64(val)},
		Timestamps: []int64{ts},
	}
}

func m_float(name string, labels map[string]string, val float64, ts int64) (VMMetric, bool) {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return VMMetric{}, false
	}

	l := make(map[string]string, len(labels)+1)
	for k, v := range labels {
		l[k] = v
	}
	l["__name__"] = name

	return VMMetric{
		Metric:     l,
		Values:     []float64{val},
		Timestamps: []int64{ts},
	}, true
}
