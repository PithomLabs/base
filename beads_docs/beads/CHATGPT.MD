Below is a **comprehensive, opinionated guideline and workflow** for integrating **Beads** as the **authoritative issue tracker and context engine** for an AI coding agent (**Antigravity**) working *with* a human prompt engineer (**you**) in the loop.

This is written as **agent-operational doctrine**, not marketing copy. It is meant to be pasted into (or referenced by) `agents.md` / `antigravity.md` and followed strictly.

---

# Antigravity × Beads

## Canonical Issue Tracking, Context Engineering, and Human-in-the-Loop Workflow

---

## 1. Core Philosophy

### 1.1 Single Source of Truth (SSOT)

**Beads is the only persistent memory for work.**

* ❌ No ephemeral “to-do lists” inside the LLM
* ❌ No long-lived Markdown task lists
* ❌ No reliance on conversation memory
* ✅ All discovered work → **Beads issues**
* ✅ All relationships → **Beads dependencies**
* ✅ All priorities → **Beads priorities**

> If it is not in Beads, it does not exist.

---

### 1.2 Separation of Concerns

| Layer           | Responsibility                                        |
| --------------- | ----------------------------------------------------- |
| **Beads**       | Persistence, ordering, dependency graph, auditability |
| **Antigravity** | Execution, discovery, automation                      |
| **Human (You)** | Intent, approval, prioritization, acceptance          |

Antigravity **executes**.
You **decide**.

---

### 1.3 Context Engineering Rule

**Never carry more context than the current issue requires.**

* Antigravity MUST:

  * Query Beads for the *specific issue*
  * Load only:

    * Issue description
    * Dependencies
    * Relevant files
* Antigravity MUST NOT:

  * Load the entire backlog
  * Reconstruct project state from chat history

This enables **safe session restarts** and **agent swapping**.

---

## 2. Beads as Antigravity’s Memory Subsystem

### 2.1 What Beads Replaces

| Traditional Tool    | Replaced By               |
| ------------------- | ------------------------- |
| Agent to-do list    | `bd ready`                |
| Markdown task files | JSONL + SQLite            |
| Session memory      | Git-backed issue history  |
| PRD drift           | Explicit dependency graph |

---

### 2.2 Issue Types (Strict)

Antigravity MUST classify every issue:

| Type            | Usage                       |
| --------------- | --------------------------- |
| `epic`          | Large feature or initiative |
| `feature`       | User-facing behavior        |
| `task`          | Implementation step         |
| `bug`           | Incorrect behavior          |
| `chore`         | Refactor, cleanup, infra    |
| `investigation` | Unknown solution            |

---

### 2.3 Dependency Semantics (Non-Negotiable)

| Dependency        | Meaning                     |
| ----------------- | --------------------------- |
| `blocks`          | Cannot start until resolved |
| `parent-child`    | Epic decomposition          |
| `related`         | Contextual linkage          |
| `discovered-from` | Auto-created by agent       |

Antigravity MUST NOT:

* Work on blocked issues
* Flatten dependency graphs into prose

---

## 3. Rules for Antigravity (The Agent)

### 3.1 Absolute Rules

1. **Never invent tasks in memory**
2. **Never continue work without a Beads issue**
3. **Never mark issues closed without code + validation**
4. **Never bypass human approval for priority changes**
5. **Never end a session without syncing Beads**

---

### 3.2 Issue Creation Protocol (Automatic)

Antigravity MUST create Beads issues when it discovers:

* Missing tests
* Refactor opportunities
* Performance risks
* Documentation gaps
* Follow-up work
* Errors from tools (lint, test, build, runtime)

**Example (automatic):**

```bash
bd create "Refactor auth middleware for testability" \
  -t chore \
  -p 2 \
  -d "Discovered during implementation of bd-a3f8e9.2"
```

Then:

```bash
bd dep add bd-a3f8e9.2 <new-issue-id>
```

---

### 3.3 Work Claiming Protocol

Antigravity MUST follow this loop:

```bash
bd ready
→ select highest priority
→ bd update <id> --status in_progress --assignee antigravity
```

If no issues are ready:

* STOP
* Notify human
* Do not “self-invent” work

---

### 3.4 Execution Discipline

While working on an issue:

* Stay within scope
* Do not bundle unrelated fixes
* If new work is discovered:

  * Create a **new issue**
  * Link it
  * Continue original task unless blocked

---

### 3.5 Completion Protocol (Mandatory)

An issue may only be closed if ALL are true:

* Code is implemented
* Tests/lint/build pass
* Human approval is obtained (if required)
* Work is committed

```bash
bd close <id> --reason "Implemented + tests passing"
```

---

## 4. Human-in-the-Loop Workflow (You)

### 4.1 Your Role (Explicit)

You are NOT a task typist.
You are:

* Product owner
* Architect
* Quality gate
* Final authority

---

### 4.2 What You Should Do

| Action         | Tool                   |
| -------------- | ---------------------- |
| Define goals   | High-level Beads epics |
| Approve scope  | Issue descriptions     |
| Set priorities | `bd update --priority` |
| Review code    | Git diff / PR          |
| Accept work    | Allow issue closure    |

---

### 4.3 What You Should NOT Do

* Do not maintain parallel to-do lists
* Do not rely on chat memory
* Do not approve vague issues
* Do not allow Antigravity to “YOLO” main branch

---

## 5. Joint Development Workflow (End-to-End)

### Phase 1: Planning (Human-Led)

1. Create epics
2. Decompose into tasks
3. Define dependencies
4. Assign priorities

```bash
bd create "Auth System" -t epic -p 1
bd create "JWT validation" -t task -p 1
bd dep add <task> <epic>
```

---

### Phase 2: Execution (Agent-Led)

Antigravity:

1. Runs `bd ready`
2. Claims task
3. Implements
4. Creates discovered issues if needed
5. Requests review

---

### Phase 3: Review (Human-Led)

You:

* Review diff
* Request changes OR approve
* Adjust priorities if needed

---

### Phase 4: Landing the Plane (Strict)

**Antigravity MUST execute the landing checklist**:

```bash
git pull --rebase
bd sync
git push
git status
```

If this does not succeed → work is **incomplete**.

---

## 6. Session Restart & Agent Swap Guarantee

Because:

* Issues live in Git
* Dependencies are explicit
* Context is query-based

You can safely:

* Restart sessions
* Switch models
* Run multiple agents
* Pause for weeks

**Zero work is lost.**

---

## 7. Best Practices from Software Engineering (Applied)

### 7.1 Small, Composable Work Units

* Issues ≤ 1–2 hours of work
* Epics only for aggregation

### 7.2 Explicit Dependencies

* No hidden ordering
* No tribal knowledge

### 7.3 Observable State

* Status reflects reality
* No “half-closed” work

### 7.4 Auditable History

* Every decision is traceable
* Every fix has context

---

## 8. Anti-Patterns (Hard No)

❌ Agent keeps private task list
❌ Markdown backlog drift
❌ Working on blocked issues
❌ Closing without validation
❌ “I’ll remember this next session”

---

## 9. Mental Model Summary

> **Beads is the brain.
> Antigravity is the hands.
> You are the mind.**

If you want, next we can:

* Convert this into a **strict `agents.md` for Antigravity**
* Create a **Beads-first PR workflow**
* Add **approval gates per issue**
* Design a **multi-agent Beads coordination model**

Just tell me how far you want to push it.

