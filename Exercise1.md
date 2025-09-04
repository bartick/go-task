# Answers To Additional Questions in Exercise 1

## 1. Which indexes are used in `/tasks (GET)` and why?

The given `CREATE TABLE` statement defines `idx_parent` and `idx_category`. For a generic `GET /tasks` endpoint that fetches all top-level tasks, neither of these would be the primary index used for sorting.

However, the `queryGetTaskHierarchy` in `model/tasks.go` gives a strong clue about how tasks are sorted: `ORDER BY th.priority DESC, th.created_at ASC`. A typical `GET /tasks` endpoint would likely use this same sorting logic.


## 2. Explain the importance of `idx_parent` in the `/tasks/:id/subtasks` API route.

The `idx_parent` on the `parent_task_id` column is absolutely critical for the performance of the `/tasks/:id/subtasks` route.

This route is powered by the `GetTaskWithSubtasks` function, which uses a recursive SQL query. The core of this query's recursive step is this join:
`... FROM tasks t INNER JOIN task_hierarchy th ON t.parent_task_id = th.id`

Why `idx_parent` is essential:

 - The Recursive Lookup: This query starts with a single task ID and then repeatedly needs to find all tasks whose parent_task_id matches an ID it has already found. This happens for every level of the hierarchy.

 - Scenario WITHOUT an Index: For every parent task in the tree, the database would be forced to perform a full table scan on the tasks table to find its children. If you have a project with 20 tasks that have children, the database might scan the entire table 20 times. This is incredibly inefficient and scales terribly.

 - Scenario WITH an Index: The `idx_parent` acts like a pre-built map. When the database needs to find children for a specific `parent_task_id`, it uses the index to look up that ID and go directly to the exact rows for its children almost instantly. This lookup is extremely fast (logarithmic time) compared to a full table scan (linear time).

 - In short, `idx_parent` prevents the database from doing an expensive full-table search for every node in the task tree, making the feature viable for any non-trivial number of tasks.


## 3. How can task search by `status` and `due_date` in the same query be improved? And why?

A common use case is to find tasks with a specific status that are due by a certain date. For example: "Find all 'in_progress' tasks that are overdue". The query would look like this:

```sql
SELECT * FROM tasks
WHERE status = 'in_progress' AND due_date < NOW();
```

Improvement: Create a composite index on (status, due_date).

```sql
CREATE INDEX idx_status_due_date ON tasks (status, due_date);
```

Why this specific index is the best solution:

 - Filtering Efficiency: A composite index creates a multi-level sorted structure. The data is first sorted by `status`. Then, within each group of identical statuses (e.g., within all the 'todo' tasks), the data is sorted again by `due_date`.

 - How the Database Uses It: When the `WHERE` clause is executed, the database can use this index very efficiently:

   1. It uses the first part of the index (`status`) to instantly jump to the block of all tasks where `status = 'in_progress'`.

   2. It then scans this much smaller, pre-sorted block, using the second part of the index (`due_date`) to find all tasks that meet the date condition. It stops scanning as soon as the `due_date` is past the one specified.

 - Why Not Single Indexes? If we only had separate indexes on `status` and `due_date`, the database would have to pick one, find all matching rows (which could be a lot), and then manually filter those results based on the second condition. The composite index allows the database to use both criteria simultaneously at the index level, which is significantly faster.