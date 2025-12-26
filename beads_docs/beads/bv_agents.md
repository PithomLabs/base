### Using bv as an AI sidecar

  bv is a fast terminal UI for Beads projects (.beads/beads.jsonl). It renders lists/details and precomputes dependency metrics (PageRank, critical path, cycles, etc.) so you instantly see blockers and execution order. For agents, it’s a graph sidecar: instead of parsing JSONL or risking hallucinated traversal, call the robot flags to get deterministic, dependency-aware outputs.

  - bv --robot-help — shows all AI-facing commands.
  - bv --robot-insights — JSON graph metrics (PageRank, betweenness, HITS, critical path, cycles) with top-N summaries for quick triage.
  - bv --robot-plan — JSON execution plan: parallel tracks, items per track, and unblocks lists showing what each item frees up.
  - bv --robot-priority — JSON priority recommendations with reasoning and confidence.
  - bv --robot-recipes — list recipes (default, actionable, blocked, etc.); apply via bv --recipe <name> to pre-filter/sort before other flags.
  - bv --robot-diff --diff-since <commit|date> — JSON diff of issue changes, new/closed items, and cycles introduced/resolved.

  Use these commands instead of hand-rolling graph logic; bv already computes the hard parts so agents can act safely and quickly.
```
