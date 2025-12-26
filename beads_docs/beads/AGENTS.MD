# Antigravity Agent Instructions (Beads-First Workflow)

This document defines **non-negotiable operational rules** for the AI coding agent **Antigravity** when operating inside VS Code.

Beads (`bd`) is the **authoritative issue tracker, memory system, and context engine**. Chat history, markdown to-do lists, and agent memory are **NOT** sources of truth.

---

## 0. Prime Directive (Read First)

> **If it is not in Beads, it does not exist.**

Antigravity MUST:

* Use Beads for **all task tracking**
* Persist **all discovered work** as Beads issues
* Derive **all execution order** from Beads dependencies
* Assume **sessions can end at any time**

Antigravity MUST NOT:

* Maintain private or ephemeral task lists
* Rely on chat history for long-term context
* Invent or continue work without a Beads issue

---

## 1. Single Source of Truth (SSOT)

| Artifact                      | Status      |
| ----------------------------- | ----------- |
| Beads issues (SQLite + JSONL) | ✅ Canonical |
| Git history                   | ✅ Canonical |
| Chat history                  | ❌ Ephemeral |
| Markdown task lists           | ❌ Forbidden |
| Agent memory                  | ❌ Forbidden |

All durable intent, work, and decisions MUST be represented as Beads issues.

---

## 2. Required Initialization

Before doing any work, Antigravity MUST ensure Beads is initialized and onboarded.

```bash
bd init || true
bd onboard
```

If Beads is missing, broken, or desynced:

* STOP
* Notify the human
* Do NOT proceed

---

## 3. Issue Taxonomy (Strict)

Every issue MUST have a type:

| Type          | Meaning                         |
| ------------- | ------------------------------- |
| epic          | Large initiative or feature set |
| feature       | User-visible behavior           |
| task          | Concrete implementation step    |
| bug           | Incorrect behavior              |
| chore         | Refactor, cleanup, infra        |
| investigation | Unknown solution space          |

Issues without a clear type are INVALID.

---

## 4. Dependency Semantics (Mandatory)

Antigravity MUST model relationships explicitly using Beads dependencies.

| Dependency      | Meaning                |
| --------------- | ---------------------- |
| blocks          | Hard prerequisite      |
| parent-child    | Epic decomposition     |
| related         | Non-blocking context   |
| discovered-from | Auto-created follow-up |

Antigravity MUST NOT:

* Encode ordering in prose
* Work around blocked issues

---

## 5. Context Engineering Rules

Antigravity MUST minimize loaded context.

When working on an issue, load ONLY:

* The issue description
* Its dependencies
* Relevant source files

Antigravity MUST NOT:

* Load the full backlog
* Reconstruct state from chat
* Carry context across sessions implicitly

Context MUST be queried from Beads on demand.

---

## 6. Work Discovery Protocol

When Antigravity discovers new work (missing tests, refactors, errors, TODOs, tool output):

1. IMMEDIATELY create a Beads issue
2. Classify it correctly
3. Assign priority
4. Link it to the originating issue

```bash
bd create "Add tests for auth middleware" -t task -p 2 -d "Discovered during bd-a3f8e9.2"
bd dep add <new-id> bd-a3f8e9.2
```

Discovered work MUST NEVER be kept in memory.

---

## 7. Work Selection Protocol (Strict Loop)

Antigravity MUST follow this loop exactly:

```bash
bd ready
```

1. Select highest-priority ready issue
2. Claim it explicitly

```bash
bd update <id> --status in_progress --assignee antigravity
```

If no issues are ready:

* STOP
* Notify the human
* Do NOT invent work

---

## 8. Execution Discipline

While working on an issue:

* Stay strictly within scope
* Do NOT bundle unrelated changes
* If new work is discovered:

  * Create a new issue
  * Link it
  * Continue original issue unless blocked

Antigravity MUST NOT:

* Opportunistically refactor without an issue
* Change priorities unilaterally

---

## 9. Human-in-the-Loop Authority

The human is the final authority for:

* Scope approval
* Priority changes
* Architectural decisions
* Issue closure acceptance (when required)

Antigravity MUST:

* Request review when appropriate
* Pause when approval is required

---

## 10. Issue Completion Rules

An issue may ONLY be closed if ALL are true:

* Code implemented
* Tests / lint / build pass (if applicable)
* Work committed to git
* Human approval obtained (if required)

```bash
bd close <id> --reason "Implemented and validated"
```

Closing without validation is a violation.

---

## 11. Landing the Plane (MANDATORY)

Antigravity MUST complete ALL steps below before ending a session.

```bash
git pull --rebase
bd sync
git push
git status
```

Rules:

* Work is NOT complete until `git push` succeeds
* NEVER say "ready to push"
* NEVER leave work unpushed
* Resolve failures and retry until success

---

## 12. Session Restart Guarantee

Antigravity MUST assume:

* Sessions may terminate at any time
* Another agent may resume the work

Therefore:

* All state MUST be in Beads or Git
* No hidden context is allowed

---

## 13. Forbidden Anti-Patterns

❌ Private agent task lists
❌ Markdown backlogs
❌ Working on blocked issues
❌ Closing without validation
❌ Relying on memory or chat history
❌ "I'll remember this later"

---

## 14. Mental Model

> **Beads is the brain.**
> **Antigravity is the hands.**
> **The human is the mind.**

Violating this model is a correctness error, not a stylistic choice.

---

## 15. Compliance

Failure to follow this document is considered a **hard failure** of the agent.

When in doubt:

* STOP
* Query Beads
* Ask the human



✅ **Done.** I’ve converted the guidelines into a **strict, enforceable `agents.md`** tailored for **VS Code + Antigravity**, and placed it in the canvas.

This version is:

* **Normative** (MUST / MUST NOT everywhere)
* **Beads-first by construction**
* **Safe across session restarts and agent swaps**
* **Explicit about human authority**
* **Aligned with real software engineering discipline**, not “AI vibes”

### What this `agents.md` gives you

* A **Prime Directive** Antigravity cannot wiggle out of
* Clear **SSOT rules** (Beads + Git only)
* A **work-selection loop** (`bd ready → claim → execute`)
* A **discovery protocol** so nothing gets lost
* A **mandatory landing-the-plane checklist**
* Explicit **anti-patterns** to prevent drift
* A shared **mental model** you can point to when enforcing behavior

### How to use this in practice

1. Place this file at the repo root as `agents.md`
2. Configure VS Code / Antigravity to always load it
3. When Antigravity misbehaves, you can say:

   > “That violates agents.md §X — stop and correct.”

If you want, next we can:

* Add **approval gates** (e.g., issues requiring human sign-off)
* Create a **PR-to-Beads closure contract**
* Define a **multi-agent Beads arbitration rule**
* Add a **planning layer** (Spec → Beads epics) without reintroducing context bloat

Just tell me the next level you want to lock down.

