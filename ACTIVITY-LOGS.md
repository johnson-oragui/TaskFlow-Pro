
## **1. Core Concept of Activity Logs**

An **Activity Log** tracks significant actions or events within the application, such as:

- Task updates (e.g., status change, assignee change)
- Team changes (e.g., member joined, role assigned)
- File uploads
- Comments or notes added
- Any user interactions that require auditing or accountability.

**Purpose of Activity Logs:**
- **Audit Trail:** Provides visibility into who performed what action and when.
- **Transparency:** Makes it easier to trace changes to tasks, teams, or projects.
- **Notifications:** Activity logs can trigger real-time notifications for users when tasks are updated or comments are added.

---

## **2. Database Modeling**

### Activity Log Table
```sql
activity_logs
----------------
id            UUID (PK)
user_id       UUID (FK to users)
team_id       UUID (FK to teams)
task_id       UUID (FK to tasks, optional)
action        ENUM('task_created', 'task_updated', 'task_assigned', 'comment_added', 'status_changed', 'role_assigned', etc.)
meta          JSONB   -- Additional metadata like old/new values, changes made, etc.
created_at    TIMESTAMP
```

### **Columns Explanation**:
- **user_id**: The ID of the user who performed the action.
- **team_id**: The team where the action was performed.
- **task_id**: (Optional) The ID of the task involved, if applicable.
- **action**: A string that indicates what kind of action occurred. This can be an enum value like `task_created`, `status_changed`, `assigned_to_updated`, etc.
- **meta**: A JSONB field that stores additional information relevant to the action. For example:
  - When a task's status is updated, it may store:
    ```json
    {
      "from": "Todo",
      "to": "In Progress"
    }
    ```
  - When a user is assigned to a task, it could store:
    ```json
    {
      "assigned_to": "user_id",
      "previous_assignee": "previous_user_id"
    }
    ```
- **created_at**: Timestamp of when the action occurred.

---

## **3. Activity Log Relationships**

```
users ─────────────────────────────┐
                                    ▼
                             activity_logs ──▶ tasks
                                    │
                                    ▼
                               teams
```

### Relationships:
- A **user** can create multiple **activity logs** (one-to-many).
- An **activity log** is associated with one **team** (many-to-one).
- An **activity log** may involve a **task** (many-to-one).
- **Meta data** in the logs stores extra context about changes, providing a detailed history for each task/team action.

---

## **4. Use Case Scenarios for Activity Logs**

### **Task Created (task_created)**

- **Event**: A new task is created by a user.
- **Data captured**: 
  - **User**: The user who created the task.
  - **Task**: Task details like title, description, assigned member, priority.
  - **Meta**: Any details related to task creation (e.g., initial labels or tags).
- **Action logged**:
  ```json
  {
    "action": "task_created",
    "meta": {
      "title": "Create Project Plan",
      "assigned_to": "user_id",
      "priority": "High"
    }
  }
  ```

---

### **Task Status Change (status_changed)**

- **Event**: A user updates the task's status (e.g., from `Todo` to `In Progress`).
- **Data captured**:
  - **User**: The user who made the update.
  - **Task**: The task being updated.
  - **Meta**: The previous and current status values.
- **Action logged**:
  ```json
  {
    "action": "status_changed",
    "meta": {
      "from": "Todo",
      "to": "In Progress"
    }
  }
  ```

---

### **Task Assignment (task_assigned)**

- **Event**: A task is assigned to a user (or reassigned).
- **Data captured**:
  - **User**: The user who performed the assignment.
  - **Task**: The task being assigned.
  - **Meta**: Information about the user being assigned or reassigned.
- **Action logged**:
  ```json
  {
    "action": "task_assigned",
    "meta": {
      "assigned_to": "user_id",
      "previous_assignee": "previous_user_id"
    }
  }
  ```

---

### **Comment Added (comment_added)**

- **Event**: A comment is added to a task.
- **Data captured**:
  - **User**: The user who added the comment.
  - **Task**: The task the comment belongs to.
  - **Meta**: The content of the comment added.
- **Action logged**:
  ```json
  {
    "action": "comment_added",
    "meta": {
      "comment": "This task is urgent, please prioritize it!"
    }
  }
  ```

---

### **Role Assigned (role_assigned)**

- **Event**: A user’s role within a team changes.
- **Data captured**:
  - **User**: The user whose role is being changed.
  - **Team**: The team where the role change occurred.
  - **Meta**: The old and new roles.
- **Action logged**:
  ```json
  {
    "action": "role_assigned",
    "meta": {
      "user_id": "user_id",
      "from": "member",
      "to": "manager"
    }
  }
  ```

---

## **5. Activity Log API Endpoints**

| Method | Endpoint                        | Description                                      |
|--------|---------------------------------|--------------------------------------------------|
| POST   | `/teams/:teamId/activity_logs`  | Create a new activity log                        |
| GET    | `/teams/:teamId/activity_logs`  | Retrieve all activity logs for a team           |
| GET    | `/tasks/:taskId/activity_logs`  | Retrieve all activity logs for a specific task  |
| GET    | `/users/:userId/activity_logs`  | Retrieve all activity logs for a user           |

---

## **6. Real-Time Activity Logging**

- **WebSockets** or **Server-Sent Events (SSE)** can be used to broadcast activity log updates in real-time to clients. 
  - When an activity log is created, notify relevant team members, managers, or anyone assigned to the task.
  
  Example: Upon a status update for a task, all team members can be notified in real-time.

- WebSockets API:
  ```go
  go func() {
      for {
          select {
          case activity := <-activityLogChannel:
              // Broadcast the activity log to team members
              socket.SendToUser(activity.UserID, activity)
          }
      }
  }()
  ```

---

## **7. Activity Log Analytics & Dashboard**

The activity logs can also be aggregated into a dashboard for tracking user actions, task progress, and project audits:
- **Dashboard Stats**:
  - Total number of tasks created in a period
  - Number of comments added
  - Most active users
  - Task status change frequency
  
Example queries:
- "Show me the number of tasks created today."
- "Show the most active user in the last week."

---

## **8. Permissions and Data Privacy**

- **Admins** have full access to all logs across teams and tasks.
- **Managers** can view logs related to their teams.
- **Members** can only view logs for their assigned tasks and actions taken on those tasks.

Example (middleware for permissions check):

```go
func RequireAdminAccess(teamID, userID string) bool {
    member := GetTeamMember(teamID, userID)
    return member.Role.Name == "admin"
}
```

---

## **9. Example Activity Log API Calls**

### Create Activity Log (Task Created)
```json
POST /teams/123/activity_logs
{
    "user_id": "user-123",
    "team_id": "team-123",
    "task_id": "task-123",
    "action": "task_created",
    "meta": {
        "title": "New Task",
        "assigned_to": "user-456",
        "priority": "High"
    }
}
```
