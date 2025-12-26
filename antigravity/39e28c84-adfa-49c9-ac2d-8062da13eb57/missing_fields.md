# Missing Ticket Form Fields

## Fields in Database BUT NOT in Form

### 1. **labels** (NEW - Beads tags) ❌
- **Type:** `TEXT` (JSON array: `["backend", "frontend"]`)
- **Purpose:** Beads labels for categorization
- **Currently:** Not in form, defaults to `[]`
- **Fix:** Add multi-select autocomplete

### 2. **tags** (EXISTING - Legacy) ❌  
- **Type:** `TEXT` (JSON array)
- **Purpose:** Original tagging system (backward compat)
- **Currently:** Not in form, defaults to `[]`
- **Note:** Different from `labels`

### 3. **parent_id** (NEW - Epics) ❌
- **Type:** `INTEGER` (foreign key to tickets)
- **Purpose:** Link subtasks to parent epic
- **Currently:** Not in form, defaults to `NULL`
- **Fix:** Add parent ticket selector

### 4. **dependencies** (NEW - Blockers) ❌
- **Type:** `TEXT` (JSON array of ticket IDs: `[1, 2, 3]`)
- **Purpose:** Track blocking tickets
- **Currently:** Not in form, defaults to `[]`
- **Fix:** Add multi-select of tickets

### 5. **discovery_context** (NEW) ❌
- **Type:** `TEXT`
- **Purpose:** Link to parent issue for context
- **Currently:** Not in form, defaults to `NULL`
- **Optional field**

### 6. **closed_reason** (NEW) ❌
- **Type:** `TEXT`
- **Purpose:** Notes when closing ticket
- **Currently:** Not in form, defaults to `NULL`
- **Fix:** Show when status = CLOSED

---

## Current Form State (line 165-190 in Tickets.tsx)

```tsx
const [title, setTitle] = useState("");
const [description, setDescription] = useState("");
const [status, setStatus] = useState("OPEN");
const [priority, setPriority] = useState("MEDIUM");
const [type, setType] = useState("TASK");
const [assigneeId, setAssigneeId] = useState<number | null>(null);
// ❌ Missing: labels, tags, parentId, dependencies
```

## Current API Payload (line 173)

```tsx
const payload = {
    title,
    description,
    status,
    priority,
    type,
    assigneeId: assigneeId || undefined
    // ❌ Missing: labels, tags (dependencies not needed for create)
};
```

---

## Quick Fix: Add Missing Fields

### 1. Add State Variables

```tsx
const [labels, setLabels] = useState<string[]>([]);
const [tags, setTags] = useState<string[]>([]); // Legacy support
const [parentId, setParentId] = useState<number | null>(null);
const [dependencies, setDependencies] = useState<number[]>([]);
```

### 2. Update Payload

```tsx
const payload = {
    title,
    description,
    status,
    priority,
    type,
    labels,     // NEW
    tags,       // NEW (optional for backward compat)
    assigneeId: assigneeId || undefined
};
```

### 3. Add Form Fields

**Labels (Priority 1):**
```tsx
<div>
  <label>Labels</label>
  <Autocomplete
    multiple
    value={labels}
    onChange={(_, val) => setLabels(val)}
    options={["backend", "frontend", "security", "ui", "database"]}
    freeSolo  // Allow custom labels
  />
</div>
```

**Tags (Optional - Legacy):**
```tsx
<div>
  <label>Tags (Legacy)</label>
  <Input 
    value={tags.join(", ")} 
    onChange={(e) => setTags(e.target.value.split(",").map(t => t.trim()))}
    placeholder="tag1, tag2, tag3"
  />
</div>
```

**Parent Ticket (For Epics):**
```tsx
<div>
  <label>Parent Epic</label>
  <Select value={parentId} onChange={(_, val) => setParentId(val)}>
    <Option value={null}>None (Top-level)</Option>
    {tickets.filter(t => t.type === "EPIC").map(t => (
      <Option value={t.id}>{t.title}</Option>
    ))}
  </Select>
</div>
```

---

## Summary Table

| Field | In DB | In Form | Being Saved | Priority |
|-------|-------|---------|-------------|----------|
| title | ✅ | ✅ | ✅ | - |
| description | ✅ | ✅ | ✅ | - |
| status | ✅ | ✅ | ✅ | - |
| priority | ✅ | ✅ | ✅ | - |
| type | ✅ | ✅ | ✅ | - |
| assignee_id | ✅ | ✅ | ✅ | - |
| **labels** | ✅ | ❌ | ❌ | **HIGH** |
| **tags** | ✅ | ❌ | ❌ | LOW |
| **parent_id** | ✅ | ❌ | ❌ | MED |
| **dependencies** | ✅ | ❌ | ❌ | MED |
| beads_id | ✅ | N/A | Auto | - |
| issue_type | ✅ | N/A | Auto (from type) | - |

**Key Missing:** `labels`, `tags`, `parent_id`, `dependencies`
