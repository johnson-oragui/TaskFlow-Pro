## **1. Core Concepts**

**Tasks** are the heart of TaskFlow — they represent individual units of work that belong to a **team** and are optionally assigned to a **user**.
Every task can have properties such as **status**, **priority**, **due dates**, **tags**, **attachments**, and **activity history**.

---

## **2. Task Fields & Behaviors**

| Field        | Type         | Description                                                                 |
|--------------|--------------|-----------------------------------------------------------------------------|
| Title        | String       | Required. A short description of the task                                  |
| Description  | Text         | Optional. More detailed explanation                                         |
| Status       | Enum         | One of: `Todo`, `In Progress`, `Done`                                      |
| Priority     | Enum         | One of: `Low`, `Medium`, `High`                                            |
| Due Date     | Date         | Optional. Used to highlight upcoming/overdue tasks                         |
| Assigned To  | UUID         | FK to `users.id` — must be a member of the same team                       |
| Created By   | UUID         | FK to `users.id` — the user who created the task                           |
| Labels       | JSON/String  | List of user-defined tags (e.g., "frontend", "bug")                        |
| Attachments  | List/File    | Optional. File upload support (linked to file storage or DB path)          |

---

## **3. Database Schema**

```sql
tasks
---------
id             UUID (PK)
team_id        UUID (FK to teams)
title          VARCHAR
description    TEXT
status         ENUM('Todo', 'In Progress', 'Done') DEFAULT 'Todo'
priority       ENUM('Low', 'Medium', 'High') DEFAULT 'Medium'
due_date       DATE
assigned_to    UUID (FK to users)
created_by     UUID (FK to users)
labels         JSONB
created_at     TIMESTAMP
updated_at     TIMESTAMP
```

---

## **4. Task Relationships**

```
teams ───────────────┐
                     ▼
                  tasks ───── assigned_to ───▶ users
                     │
         ┌───────────┴────────────┐
         ▼                        ▼
  attachments                 activity_logs
```

- A `team` has many `tasks`
- A `task` can have an optional assignee (who must be a team member)
- A `task` has many `attachments` and `activity_logs`

---

## **5. Task Permissions Logic**

| Role     | Allowed Actions                                                            |
|----------|-----------------------------------------------------------------------------|
| Admin    | Full CRUD on all tasks in the team                                         |
| Manager  | Full CRUD, but can't delete tasks created by Admin unless it's their own   |
| Member   | Can only create tasks for themselves or mark assigned tasks as completed   |

**Examples:**
- Member can't assign a task to another member unless they're an Admin/Manager.
- Member can only view tasks within their team.

---

## **6. Use Case Scenarios**

### Task Creation (POST `/teams/:teamId/tasks`)
- Payload: `title`, `description`, `priority`, `assigned_to`, `due_date`, `labels`
- Backend checks:
  - `teamId` exists and user is a team member
  - `assigned_to` is a member of the same team
  - `title` is non-empty
- On success, task is inserted with a default status of `Todo`.

---

### Task Assignment
- A manager or admin assigns a task to a team member.
- If an assignee is not in the team → reject.
- Update `assigned_to` field in DB.

---

### Status Update (PATCH `/tasks/:id`)
- Only the assignee or a manager/admin can change the status of a task.
- Allowed status transitions:  
  `Todo → In Progress → Done`, or rollback if needed.

---

### Labeling
- Labels are stored as a JSON array or a comma-separated string.
- These help with filtering: "show all `frontend` tasks due this week".

---

### Deleting a Task
- Only an admin or the creator of the task can delete it.
- Task is soft-deleted (optional) or hard-deleted depending on your business rules.

---

## **7. Task Attachment Support (File Uploads)**

- Endpoint: `POST /tasks/:id/attachments`
- Uploads are linked to the task and stored in a separate `attachments` table or file storage system.
- Only team members can view/download attachments.

```sql
attachments
------------
id        UUID
task_id   UUID (FK to tasks)
file_url  TEXT
uploaded_by UUID (FK to users)
created_at TIMESTAMP
```

---

## **8. Task Activity Log**

Every action taken on a task (status update, assignment, edit, comment) is logged.

```sql
activity_logs
----------------
id         UUID
task_id    UUID
team_id    UUID
user_id    UUID
action     TEXT     -- e.g., "status_updated", "assigned_to_changed"
meta       JSONB    -- structured data for UI (old/new values)
created_at TIMESTAMP
```

Example:
```json
{
  "action": "status_updated",
  "meta": {
    "from": "Todo",
    "to": "In Progress"
  }
}
```

---

## **9. Sample Task APIs**

| Method | Endpoint                            | Description                        |
|--------|-------------------------------------|------------------------------------|
| POST   | `/teams/:teamId/tasks`              | Create a new task                  |
| GET    | `/teams/:teamId/tasks`              | List all tasks for a team          |
| GET    | `/tasks/:id`                        | Get task details                   |
| PATCH  | `/tasks/:id`                        | Update task (status, assignee)     |
| DELETE | `/tasks/:id`                        | Delete task                        |
| POST   | `/tasks/:id/attachments`            | Upload file to a task              |
| GET    | `/tasks/:id/activities`             | View change logs for a task        |

---

## **10. Dashboard Integration**

Tasks feed directly into team dashboards with metrics like:

- Total tasks
- Tasks completed this week
- Tasks overdue
- Tasks by priority
- Tasks per member

---

## **11. Pagination, Filters, and Search**

Enable flexible queries:
```
GET /teams/:id/tasks?status=Todo&priority=High&due_before=2025-04-15&page=2&limit=10
```

Add search indexing on `title`, `labels`, and `description`.
