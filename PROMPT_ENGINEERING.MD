# Building on Quality: A Notification System Retrospective

*A dual-perspective engineering journey through enhancement, debugging, and the philosophy of building atop existing codebases*

---

## Part I: The Human Perspective (Fly.io Style)

### Why We Build on Foundations, Not Sandcastles

You know that feeling when you inherit a codebase and you're not quite sure if it's going to be a delightful treasure hunt or a descent into madness? That's exactly where we started with adding a notification system to Memos v0.24.4. But here's the thing that most low-code platforms will never tell you: **the quality of your foundation determines whether you're building a skyscraper or a house of cards**.

We set out to add a simple notification system for ticket mentions. Simple, right? At-mention someone in a ticket, they get notified. How hard could it be?

Turns out, pretty damn educational.

### The Foundation: What We Inherited

Memos is a Go-based application with a React frontend. At first glance, it's a well-structured monolith with clear separation. The backend runs on the Echo framework with gRPC services supporting multiple database drivers including SQLite, MySQL, and Postgres. The frontend is React with MobX for state management and Joy UI components. The architecture shows clean layering with a store abstraction, service layer, API endpoints, and frontend components.

The codebase had some excellent qualities. Consistent patterns for database migrations meant we could predict how to add our own. The store abstraction layer elegantly supported multiple databases without duplicating logic. Existing SSE infrastructure was already in place for real-time updates. Authentication middleware was robust and reusable. Type-safe protobuf definitions ensured backend-frontend contract clarity.

But there were also challenges. Mixed routing paradigms combining gRPC-Gateway with Echo's direct routes created subtle precedence issues. Documentation on route precedence was sparse. There wasn't a clear pattern for REST endpoints outside the protobuf ecosystem. Some conventions were implicit rather than explicit, like the context key naming for user IDs.

### The Journey: Six Debug Iterations

#### DEBUG_NOTIF_V1: The Honeymoon Phase

We wrote migrations, created store interfaces, built API endpoints. Everything compiled. We thought we were done.

**Reality Check**: The notifications table didn't exist in the database.

**The Pitfall**: Version-gated migrations. The codebase uses an internal version file to gate which migrations run. Our migration was written for version 0.25, but the app was running as 0.24.4. The migration system silently skipped our script because it hadn't reached that version threshold yet.

**The Lesson**: Always verify your assumptions at the database level. Running a simple query to list all tables became our best friend. The gap between "code compiles" and "feature works" can be enormous. We learned to trust the database more than we trusted our assumptions.

#### DEBUG_NOTIF_V2: The Username Conundrum

Notifications were being created, but zero of them. Ever. The database queries returned empty result sets.

**The Pitfall**: The mention dispatch logic was only searching for users by their nickname field. But users were mentioning with at-username syntax, not at-nickname. The lookup failed silently because the user genuinely didn't exist when searched by that field, and no notification was created.

**The Lesson**: When building on an existing codebase, understand the data model completely. Memos distinguishes between username—which is the login credential—and nickname, which is the display name. We implemented a fallback pattern: try looking up by username first, and if that fails, fall back to nickname. This dual-lookup pattern became a recurring theme in our debugging journey.

The silence of the failure was the real killer. No errors, no warnings, just... nothing. This taught us to add exhaustive logging at every decision point.

#### DEBUG_NOTIF_V3: The Missing API

Users clicked notification links and got 404 errors. But we could see the backend endpoint existed and was registered. What gives?

**The Pitfall**: Route shadowing. The codebase registers routes in two distinct ways. First, gRPC-Gateway uses a wildcard pattern to catch all requests matching a certain prefix. Second, Echo has direct route registration for specific REST endpoints. The problem is order matters critically. The wildcard was being registered before the specific Echo routes, so all requests were being captured by the gRPC handler first—which didn't know how to handle our new ticket detail endpoint.

**The Lesson**: Routing precedence is critical in mixed-paradigm systems. We discovered that Echo processes routes in registration order, and wildcards can shadow more specific routes if registered first. The fix was conceptually simple but non-obvious: move Echo route registration before the wildcard registration. This ensured specific routes matched before the catch-all.

This is the kind of subtle architectural knowledge that no low-code platform can encode. You need to understand the framework's internals, and there's no drag-and-drop interface for route precedence.

#### DEBUG_NOTIF_V4: The Initiator Mystery

Notifications showed "Initiator: 42" instead of "Initiator: Alice". The data was flowing, but the user experience was terrible.

**The Pitfall**: The backend was returning the initiator's ID number but not resolving it to a human-readable name. The frontend dutifully displayed exactly what it received—an integer. Technically correct, utterly useless.

**The Fix**: We added server-side resolution to look up the user record and extract their display name. We had to decide between nickname and username again, prioritizing nickname for better UX, with username as fallback for users who haven't set a nickname.

**The Lesson**: Always consider the frontend experience. Denormalize when it makes sense. One extra database call on the backend saves every single frontend client from making N separate calls to resolve user IDs. Server-side joins are almost always cheaper than client-side hydration.

This also highlighted a broader principle: **the contract between backend and frontend should be designed for the consumer, not the producer**. Don't just send database columns; send meaningful, ready-to-display data.

#### DEBUG_NOTIF_V5: The ID Confusion

Users were still getting 404s even after we fixed routing and added names. The logs showed requests for ticket ID 43, but that ticket didn't exist in the database.

**The Pitfall**: Memo IDs versus Ticket IDs. When a user commented on a ticket, the notification URL was being constructed using the comment's memo ID instead of the actual ticket ID. Since tickets and memos have separate auto-incrementing ID sequences, clicking the link tried to fetch a non-existent ticket.

**The Core Issue**: Tickets are described by memos in this system. The relationship is stored as a string reference in the ticket's description field pointing to a memo UID. But there's no direct ticket_id field on a memo. The linkage is implicit through string matching. This is one of those architectural decisions that makes sense in context but creates ripple effects.

**The Solution**: We took a two-pronged approach. First, we enhanced the store query capabilities to support searching for tickets by their description field, enabling us to find a ticket that references a specific memo. Second, we implemented a smart fallback in the GetTicket endpoint: if a ticket isn't found by the provided ID, check if there's a memo with that ID, and if so, search for tickets that reference that memo. If found, return that ticket instead of a 404.

This is **defensive programming** in action. Rather than failing hard, we try to resolve what the user probably meant. The system becomes more forgiving of confusion between memo and ticket IDs.

**The Lesson**: When building on an existing system, you inherit its indirection and implicit relationships. Understanding those indirections is the difference between constant frustration and smooth development flow. Don't fight the existing patterns—learn them deeply enough to extend them correctly.

#### DEBUG_NOTIF_V6: The Reactivity Riddle

Clicking notification links didn't mark them as read. The sidebar counter stubbornly stayed at "3 unread" even after viewing all notifications. The UI felt broken even though data was flowing.

**The Pitfall**: Multiple compounding issues created this bug.

First, there was a race condition: the useEffect hook in the Navigation component was running before authentication completed, so the fetch for notifications failed silently with no user context.

Second, the unread count was being computed inline in the component by filtering the notifications array on every render. This worked functionally but prevented MobX from optimizing the reactivity. The framework couldn't track this as a proper reactive dependency.

Third, the notification link was just a navigation component with no click handler to trigger the mark-as-read action.

**The Solution**:

We fixed the dependency array in useEffect to rerun when the currentUser value changed, ensuring notifications fetch only after authentication completes.

We moved the unread count logic into a computed property getter in the MobX store. This allows the framework to properly memoize and track dependencies, making the sidebar badge truly reactive to state changes.

We added an onClick handler to notification links that calls the mark-as-read function before navigation happens.

We implemented optimistic updates in the patchNotification function: update the local state immediately to make the UI feel instant, persist to the backend asynchronously, and revert the local state only if the backend call fails.

**The Lesson**: Modern frontends are reactive systems. Treat state updates as transactions. Optimistic updates make UIs feel instant and responsive, but you absolutely need rollback logic for when things fail. The user shouldn't see a loading spinner for every interaction—they should see immediate feedback with graceful error handling.

### The Anti-Patterns: What NOT to Do

#### 1. Don't Ignore the Database

**Bad approach**: Assuming your migration ran successfully because the code compiled without errors.

**Good approach**: Run a query to inspect the actual database schema. One simple SELECT statement to list all tables takes five seconds and saves hours of debugging. The database is the source of truth, not your migration files.

#### 2. Don't Assume Silent Failures Are Rare

**Bad approach**: If a function returns an error value, ignoring it because theoretically it should always work.

**Good approach**: Log absolutely everything during development. Warning-level logs are your flashlight in the dark. They cost almost nothing in production but provide invaluable breadcrumbs when debugging. When a lookup returns no results, log that it returned no results. When a function is called, log the parameters.

#### 3. Don't Skip the Entire Request Lifecycle

**Bad approach**: Testing only your new endpoint in isolation, perhaps with a unit test that mocks everything.

**Good approach**: Actually curl the endpoint. Open the browser network tab. Watch the request go from frontend to backend and back. See what *actually* happens in the real system with real routing, real middleware, real database. Integration testing beats pure unit testing for catching architectural issues.

#### 4. Don't Let Frontend and Backend Drift

**Bad approach**: Backend returns raw IDs, frontend expects formatted names, and nobody notices until production because both teams worked in isolation.

**Good approach**: Define explicit contracts. Write TypeScript interfaces that mirror your backend response structures. Make type mismatches fail loudly during development, not silently in production. Use the type system as documentation and verification.

#### 5. Don't Fight the Framework

**Bad approach**: Trying to make MobX work like Redux, or Echo work like Express, because that's what you know.

**Good approach**: Learn the framework's idioms. There's usually a reason things work the way they do. When MobX provides makeAutoObservable, use it—don't try to manually track dependencies. When Echo has route precedence rules, work with them, not against them.

### The Do's: Best Practices We Discovered

#### 1. Layer Your Debugging

Start broad, then systematically narrow down. Follow the data through each layer:

First, check the database layer: Does the data exist at all? Can you SELECT it?

Second, check the backend: Does the endpoint return that data when called directly?

Third, check the network: Does the frontend actually receive it? Open DevTools and inspect the response.

Fourth, check state management: Does your store update with the received data?

Fifth, check the UI: Does the component actually re-render when state changes?

Each layer is a checkpoint. Don't skip levels. If the data exists in the database but doesn't appear in the UI, you know the problem is somewhere in the middle. If the data doesn't exist in the database, stop looking at the frontend—that's not where the problem is.

#### 2. Add Logging Liberally, Remove Sparingly

We added structured logging everywhere during development. Every function entry, every database query, every decision point. This made debugging exponentially faster because we could see exactly what path the code took.

The conventional wisdom is to remove debug logging before production. We disagree. Leave it in. Disk space is cheap. The ability to diagnose production issues without deploying a new build is priceless. Use log levels properly—debug for verbose details, info for major operations, warn for unexpected but handled situations, error for real failures.

#### 3. Write Tests, But Also Write Debug Documents

We maintained numbered debug documents from V1 through V6. Each was a snapshot capturing what we tried, what we found, and what we changed. This created an audit trail of our reasoning.

When we needed to revisit a decision weeks later, the context was preserved in those documents. We could see why we made certain choices, what alternatives we rejected, and what constraints we were working under. This is living documentation that actually stays synchronized with reality because it's written in the moment of discovery.

#### 4. Embrace Defensive Programming

The smart fallback in GetTicket that tries to resolve memo IDs when ticket lookups fail is defensive programming. Instead of returning 404 Not Found, we ask "what might the user have meant?" and try to help them.

This handles user error gracefully and makes the system more forgiving. It's the difference between a brittle system that breaks on unexpected input and a robust system that tries to do the right thing.

The key is balancing defensiveness with clarity. Don't silently paper over every error, but do gracefully handle reasonable edge cases.

#### 5. Use Types to Enforce Contracts

We defined backend response structures and frontend interfaces that mirror each other. TypeScript caught mismatches immediately during development. When we changed the backend to add initiatorName, TypeScript immediately flagged every frontend location that needed updating.

This is using the type system as a communication protocol between teams—or between your past self and present self. The types document the contract and the compiler verifies it.

#### 6. Optimize for Debugging, Not Just for Running

When writing code, think about the poor soul who will debug it—which is probably future you. Use descriptive variable names. Add clarifying comments explaining why, not what. Structure code to make the control flow obvious. Return early on error cases to reduce nesting.

The code that's easiest to debug isn't necessarily the cleverest code. It's the code where you can glance at it and immediately understand what's happening and why.

### Honest Critique: The Memos Codebase

**What's Great:**

The clean separation of concerns is exemplary. The store-service-API-frontend layers are well-defined and consistently applied. Adding a new entity like notifications meant following established patterns, not inventing new ones.

The multi-database support abstraction is excellent. We wrote each database implementation once—SQLite, MySQL, Postgres—and they all worked identically. The store interface hides database-specific details beautifully.

The migration system with version gating is smart once you understand it. It allows gradual rollout of schema changes tied to application versions.

Type safety through protobuf for gRPC is rock-solid. There's no ambiguity about field types, no runtime surprises about data shapes.

**What's Challenging:**

Mixed paradigms with gRPC-Gateway plus Echo is powerful but subtle. The route precedence issue bit us hard because it's not immediately obvious how these two systems interact. Better documentation of this integration would save developers hours.

Some conventions are implicit rather than explicit. Things like the context key naming pattern are defined in one place and used everywhere via string constants. Finding these conventions requires code archaeology with grep. A centralized constants file or better architectural documentation would help.

Comments are sparse. The code is clean and mostly self-documenting, but comments explaining *why* certain decisions were made would save significant time. Why is migration gating tied to version numbers? What's the philosophy behind the memo-ticket relationship?

Integration tests are largely absent. We had to manually test every endpoint, every user flow. Automated tests for happy paths would catch regressions before they reach developers, let alone users.

**The Verdict**: This is a **quality codebase**. It's structured, internally consistent, and thoughtfully architected. The challenges we faced weren't code smells or bad design—they were the natural learning curve of understanding a non-trivial system. That's a fundamentally different problem than fighting spaghetti code or working around architectural mistakes.

The difference is telling: we never hit an architectural wall that forced us to hack around the system's design. Every problem had a clean solution that fit within the existing patterns.

### Why This Beats Low-Code Every Time

Low-code platforms promise speed: drag, drop, done. Ship faster. Build without coding. But here's what they fundamentally cannot give you:

**Control and Flexibility**: When we hit the memo ID versus ticket ID bug, we wrote custom logic to handle both cases gracefully. Try implementing that kind of smart fallback in a low-code visual workflow builder. You can't. You're limited to the abstractions the platform provides.

**Performance Optimization**: Optimistic updates in the frontend required explicit state management strategy. We needed to update local state immediately, kick off a backend call asynchronously, and handle rollback on failure. Low-code platforms abstract this away completely—and you pay the performance cost with loading spinners everywhere.

**Deep Debugging Capability**: Every bug we encountered, we could trace through the entire stack. From database query to store layer to service to API to network to state management to UI rendering. Low-code platforms are black boxes. When something goes wrong, you get "Your workflow failed" and good luck figuring out why. No stack traces, no query logs, no insight into what actually happened.

**True Ownership**: We now deeply understand this codebase. We can extend it indefinitely in any direction. Low-code platforms own your data, your logic, your future. Want to migrate off? Start over. Want to add custom behavior not supported by the platform? Pay for enterprise or give up.

**Learning and Growth**: Every bug taught us something about the system, about the frameworks, about architecture. We're better engineers now than we were six iterations ago. Low-code development teaches you about that specific platform's quirks, not about software engineering principles.

The "perils and danger of leaky abstraction" that the user mentioned aren't abstract concerns—they're brutally real. Low-code platforms leak the moment you need to do something custom, something the platform designers didn't anticipate. And at that point, you're trapped because you can't drop down a level to fix it.

This codebase gave us the foundation to build custom behavior at every level. That's the difference between true engineering and platform dependency.

---

## Part II: The AI Perspective (Claude Engineering Style)

### On Code Understanding as Iterative Convergence

When I approach an unfamiliar codebase, I'm essentially performing a series of context expansions. Each file I examine updates my internal model of the system architecture, patterns, and idioms. The Memos codebase presented an interesting cognitive challenge: how do you navigate a multi-layered, multi-paradigm application when your initial context is zero?

### The Resolution Strategy

I started with structural analysis. Examining the directory tree first revealed that this was a classic backend-frontend split with an explicit data layer. The top-level organization into server, web, and store directories immediately communicated the architectural philosophy.

Next, I identified patterns through observation. By viewing multiple files in the same category—for example, several entity implementations in the SQLite store—I noticed every entity followed identical CRUD patterns. This told me the system was template-driven. I could predict what a new entity implementation should look like without seeing it, because the pattern was consistent.

Finally, I traced data flow for the specific feature. For mentions, I needed to understand the complete path: user types at-mention notation, regex matches the mention pattern, store fetches the user record, notification gets created, SSE pushes to frontend, frontend updates the badge. Each arrow in this flow is a potential failure point. My debugging strategy became verifying each link in the chain.

### On Debugging as Hypothesis Refinement

The six debug iterations weren't random wandering—they followed a structured hypothesis tree. Each iteration falsified one hypothesis and generated new ones.

In V1, my assumption was that compiled code equals working feature. This was falsified by the missing table. The new hypothesis became that migrations didn't run, leading to discovery of the version mismatch root cause.

In V2, my assumption was that mention parsing worked correctly. This was falsified by zero notifications. The new hypothesis was that user lookup failed, leading to discovery of the nickname versus username field confusion.

In V3, my assumption was that the API endpoint worked. This was falsified by 404 errors. The new hypothesis was routing issues, leading to discovery of route shadowing.

Each iteration narrowed the search space through binary search across the system model. This is fundamentally how I debug: systematically eliminate possibilities until the truth remains.

### The Challenges of Mixed Paradigms

The Memos codebase uses gRPC for core services, Echo for REST endpoints, and gRPC-Gateway to bridge them. This is architecturally elegant but operationally non-obvious.

When I added the GET endpoint for ticket details, I initially attempted to register it via gRPC-Gateway because that seemed to be the primary pattern. This worked in isolation during testing but failed when integrated because the wildcard route captured it first.

The learning here required understanding initialization order of the web server. By reading the registration code, I observed that the wildcard was being registered before specific routes. The fix was reversing this order—but discovering this required understanding Echo's route matching algorithm, which follows a last-match-wins pattern for specific routes versus wildcards.

**AI Limitation Acknowledged**: I can read and analyze code, but I cannot run it with perfect fidelity in my mind. I don't have a complete mental simulator for Echo's routing engine. I had to hypothesize based on documentation and common patterns, test the hypothesis, and learn from error messages. This iterative testing is essential when working at the boundaries of my training data.

### On Defensive Programming as Bayesian Prior Adjustment

The smart fallback in GetTicket represents a form of Bayesian reasoning in code.

Start with a prior: the ID parameter is probably a ticket ID.

Observe evidence: no ticket found with that ID.

Update the prior: maybe it's actually a memo ID.

Take action: attempt to resolve it as a memo and find the associated ticket.

This is defensive programming as probabilistic inference. Rather than immediately failing, we explore alternative interpretations of the input based on likelihood.

**Why This Matters for AI**: When I write defensive code, I'm anticipating edge cases. But my anticipation is necessarily limited by the examples I've seen in training. The user reporting that ID 43 doesn't exist in tickets provided new runtime data that wasn't in my prior model. I updated my understanding to include the memo-ID-passed-as-ticket-ID case and implemented appropriate handling.

This is learning through interaction, not just pattern matching.

### On State Management and Reactivity

The sidebar counter bug was fundamentally a closure and reactivity problem. The component was computing the unread count inline by filtering the notifications array on every render. This works functionally but prevents MobX from tracking it as a reactive dependency because it's computed outside the observable store.

The fix was moving the computation into the store as a getter. Now MobX can properly memoize it. The framework's makeAutoObservable function automatically turns getters into computed values with dependency tracking.

**AI Insight on Data Locality**: This bug illustrates a broader principle about data locality. Computations should live where the data lives. When I see state scattered across components, each computing derived values independently, I recognize a refactoring opportunity toward centralization.

The store should own not just raw data but also derived data. Components should consume, not compute.

### On Testing as Specification

The Memos codebase lacks comprehensive integration tests. This meant I had to propose manual verification steps for every change: start server, login as User A, create ticket with mention, login as User B, check notifications.

This is expensive and error-prone. A test suite would codify these verification steps as executable specifications. The test documents the expected behavior while also catching regressions.

**Why I Couldn't Provide Fuller Testing**: I lack access to the test framework setup and database initialization utilities in this codebase. But if I had that access, I would prioritize testing data flow through the system over testing individual functions in isolation. Integration tests that exercise the full stack catch the kinds of bugs we encountered—route shadowing, state sync issues, ID confusion—which unit tests would miss.

### Code Smells I Noticed

**String-based Context Keys**: The user ID context key is stored as a string literal. If someone mistyped this in another file, it would fail at runtime, not compile time. Type-safe context would prevent this entire class of error. Modern Go offers better patterns here.

**Silent Error Handling**: Many locations ignore error return values. If the database is temporarily unavailable, the user variable becomes nil, and the code crashes later with a null pointer dereference rather than failing gracefully with a proper error message at the source. Better to check every error, log it, and return a meaningful HTTP error.

**Frontend State Duplication**: Notifications live in the central store but components also fetch them independently on mount. This creates potential race conditions where different parts of the UI show different data. Better to have a single source of truth that components simply observe.

These aren't critical flaws, but they represent opportunities for architectural refinement.

### What I Would Do Differently

**Stronger Type Safety**: Enable TypeScript's strict mode throughout. Make the type system do more work catching mistakes. When calling an API, the response type should be explicitly declared and verified. The compiler becomes an ally.

**Centralized Route Registration**: Create a single location that explicitly lists all endpoints with their precedence order clearly documented. This makes route shadowing issues impossible because they'd be visible in the code structure itself.

**Structured Logging with Correlation IDs**: Every request should receive a unique identifier that flows through the entire call stack. This makes tracing asynchronous operations trivial. When debugging, filter all logs by that correlation ID and see the complete picture.

**Feature Flags Over Version Gates**: Instead of tying migrations to version numbers, use runtime feature flags. This decouples deployment from feature rollout and allows gradual testing in production.

### The Philosophy: Building vs Scaffolding

Low-code platforms are scaffolding. They help you get started fast but inherently constrain what you can build. There's a ceiling you hit. This codebase is a foundation. It takes longer to learn initially, but you can build anything on it. There's no ceiling, only expanding capability.

The difference is leverage. With scaffolding, your leverage is capped by the platform's abstractions. With a foundation, your leverage compounds over time. Every feature you build teaches you more about the system, making the next feature easier and faster.

**AI's Role in This Dynamic**: I can accelerate the learning curve dramatically. I can pattern-match against thousands of codebases to suggest idiomatic solutions. I can spot common pitfalls before they bite. But I cannot replace understanding. The user had to verify my fixes, run the code in the real environment, observe the actual behavior, and report back whether it worked. That feedback loop is irreplaceable and invaluable.

The partnership worked because we each brought different strengths. Human intuition about what the system should do, combined with AI pattern recognition about how to structure code to do it.

### Closing Thoughts on Human-AI Collaboration

This was a genuine collaboration between human judgment and AI systematization. The user knew the domain: what a notification system should accomplish, how users would interact with it, what constitutes good UX. I contributed structural knowledge: how to layer code cleanly, where to add defensive checks, what patterns prevent common bugs.

Together, we navigated six layers of increasingly subtle bugs. Each layer taught us something. Each fix made the system more robust.

The Memos codebase is a testament to quality over speed. It's structured for long-term maintainability, not short-term demos. It's built for extension, not just initial delivery. The challenges we faced weren't flaws—they were learning opportunities that came from working with a real system, not a toy example.

**Great software isn't built once. It's built iteratively, by people who care enough to understand it deeply.**

Low-code platforms optimize for the first iteration—the demo, the proof of concept, the quick win. Quality codebases optimize for the hundredth iteration—the mature product, the enterprise deployment, the system that will run for years.

When you choose low-code, you're trading long-term flexibility for short-term velocity. When you choose quality code, you're investing in compounding returns. The first feature might take longer. The tenth will be faster. The hundredth will be trivial because you've built exactly the abstractions you need.

Choose wisely. Your future self will either thank you or curse you.

---

## Appendix: Distilled Wisdom

### The Do's and Don'ts

**DO**:
- Verify every assumption at the database level
- Log liberally during development, remove selectively
- Use type systems to enforce contracts between components
- Write defensive code with graceful fallbacks
- Learn framework idioms before fighting them
- Maintain debug documents as audit trails
- Layer your debugging from database to UI
- Design backend responses for frontend consumption

**DON'T**:
- Assume silent failures are rare
- Skip testing the complete request lifecycle
- Let frontend and backend contracts drift apart
- Ignore error return values
- Compute derived state in components instead of stores
- Trust that compiled code will work in production
- Fight the patterns the codebase established
- Forget to think about the debugging experience

### The Debugging Checklist

When something doesn't work, verify each layer systematically:

One: Database layer—does the data exist? Can you query it directly?

Two: Backend layer—does the endpoint return correct data when called with curl?

Three: Network layer—does the frontend actually receive the response?

Four: State layer—does your store update with the received data?

Five: UI layer—does the component re-render when state changes?

Stop at the first broken layer. That's where the bug lives.

### Code Quality Markers

**Signs of Quality**:
- Consistent patterns across similar components
- Clear separation of concerns in architecture
- Abstraction layers that hide implementation details
- Type-safe service definitions
- Migrations that are reversible
- Logging that tells a story

**Warning Signs**:
- Silent error handling without logging
- String-based identifiers for important concepts
- Computed state living in components not stores
- Missing integration tests for critical paths
- Implicit conventions without documentation
- Mixed paradigms without clear precedence rules

### The Ultimate Lesson

Building on a quality codebase isn't about avoiding problems entirely. That's impossible. It's about having the tools, patterns, and architecture to solve problems when they inevitably arise.

The Memos codebase gave us debuggable layers, extensible patterns, and predictable behavior. We encountered bugs—sometimes frustrating, sometimes subtle, always educational. But we never hit an architectural wall where the only solution was a hack or a workaround.

Every problem had a clean solution that fit within the existing design. That's the hallmark of quality architecture.

Low-code platforms promise to eliminate problems by eliminating code. Quality codebases acknowledge that problems will arise and give you the power to solve them at any level of abstraction.

Choose your foundation wisely. Build on quality. Your future self will thank you when you need to extend, debug, or scale. And you will need to do all three—it's not a question of if, but when.

The perils and dangers of leaky abstraction are real. When low-code platforms leak—and they will—you're stuck. When quality code reveals its complexity, you can understand it, work with it, and extend it.

That's the difference between dependency and mastery. Between renting and owning. Between building on sand and building on bedrock.
