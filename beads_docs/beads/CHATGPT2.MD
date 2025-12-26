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

## 15. Planning Layer (Spec → Beads Epics)

This project uses a **thin planning layer** to translate human intent into executable work **without introducing context bloat**.

### 15.1 Purpose of the Planning Layer

The planning layer exists to:

* Capture **human intent and constraints**
* Enable **deliberate design before execution**
* Feed Beads with **clean, dependency-aware epics**

It does NOT exist to:

* Serve as long-term memory
* Be fully loaded into agent context
* Replace Beads issues

The planning layer is **input-only**. Beads is **output-only**.

---

### 15.2 Spec Artifacts (Strictly Limited)

Allowed planning artifacts:

* `SPEC.md` (or equivalent)
* Short-lived design notes
* Diagrams (optional)

Rules:

* Specs MUST be concise
* Specs MUST describe **WHAT and WHY**, not step-by-step HOW
* Specs MUST be considered disposable once epics are created

Specs are **not authoritative** after translation into Beads.

---

### 15.3 Translation Protocol (Human-Led, Agent-Assisted)

**Only the human may approve the translation from Spec to Beads.**

Protocol:

1. Human writes or updates SPEC
2. Antigravity reads SPEC ONCE
3. Antigravity proposes:

   * Beads epics
   * High-level dependencies
   * Initial priorities
4. Human reviews and approves
5. Antigravity creates Beads epics

After step 5:

* SPEC is no longer required for execution
* Antigravity MUST NOT reload SPEC unless explicitly instructed

---

### 15.4 Epic Creation Rules

When creating epics from a SPEC:

* One epic per **coherent outcome**
* Epic description MUST summarize intent in ≤10 lines
* No implementation detail inside epics

Example:

```bash
bd create "Auth System" -t epic -p 1 -d "Provide secure user authentication with JWT-based sessions and role-aware authorization."
```

---

### 15.5 Decomposition Rules (Epic → Tasks)

Antigravity MAY assist decomposition ONLY AFTER epic approval.

Rules:

* Tasks must be executable independently
* Tasks must expose dependencies explicitly
* Tasks must fit within 1–2 hours of work

```bash
bd create "JWT validation middleware" -t task -p 1 -d "Validate access tokens on protected routes"
bd dep add <task-id> <epic-id>
```

---

### 15.6 Context Containment Guarantee

To prevent context bloat:

Antigravity MUST:

* Load SPEC files only during translation
* Never carry SPEC content into execution context
* Rely exclusively on Beads after epic creation

Specs MUST NOT:

* Be incrementally updated during execution
* Be used as task lists
* Be used as a substitute for Beads issues

---

### 15.7 Spec Change Management

If human intent changes:

1. Human updates SPEC
2. Human requests re-translation
3. Antigravity:

   * Updates epics
   * Adds or closes tasks
   * Preserves historical issues (do not rewrite history)

Existing closed issues are NEVER modified.

---

### 15.8 Failure Modes (Prevented)

This planning layer explicitly prevents:

❌ Markdown backlog drift
❌ Agents loading massive specs repeatedly
❌ Specs becoming shadow sources of truth
❌ Context compaction erasing intent

---

## 16. Compliance

Failure to follow this document is considered a **hard failure** of the agent.

When in doubt:

* STOP
* Query Beads
* Ask the human



✅ **Planning layer added — cleanly, safely, and without context bloat.**

I’ve updated `agents.md` with a **thin Spec → Beads epic translation layer** that:

* Preserves **deliberate planning**
* Keeps **Beads as the sole execution memory**
* Prevents specs from turning into shadow backlogs
* Explicitly constrains when and how specs are loaded

### What this gives you (architecturally)

* **Specs are input-only**
* **Beads epics are the executable contract**
* **Antigravity reads specs once, then forgets them**
* **No long-lived markdown context in agent memory**
* **Safe for restarts, compaction, and multi-agent work**

### Key design wins

* Planning without PRD sprawl
* Intent preserved, but not reloaded endlessly
* Human remains the authority at the translation boundary
* Beads stays small, queryable, and dependency-driven

### Mental model (now complete)

> **Specs define intent**
> **Beads encodes commitments**
> **Antigravity executes one issue at a time**

If you want next steps, we can:

* Add a **formal “Spec Translation Checklist”**
* Introduce **Beads labels for spec lineage**
* Create a **Beads-only architecture decision log**
* Lock this down further for **multi-agent parallelism**

This is now a production-grade agent governance model.

