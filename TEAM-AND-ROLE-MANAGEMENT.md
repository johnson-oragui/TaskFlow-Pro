

## **1. Core Concepts**

### Teams
- A **Team** represents a working group or organization unit (e.g., "Marketing Team").
- Teams can have many users.
- Teams serve as containers for:
  - Tasks
  - Users
  - Chats
  - Activity logs
  - File uploads

### Users
- A **User** can belong to multiple teams.
- A user’s role **may differ across teams** (e.g., an `admin` in Team A, but a `member` in Team B).

### Roles
- Each user has **one role per team**.
- Roles define **what the user can do** within that team.

---

## **2. Role Types & Permissions**

| Role     | Description                                 | Typical Permissions                              |
|----------|---------------------------------------------|--------------------------------------------------|
| Admin    | Team owner or creator                       | Full control: manage users, tasks, team settings |
| Manager  | Senior team member, project lead            | CRUD on tasks, assign members, view dashboard    |
| Member   | Regular contributor                         | Can view and complete tasks, add comments        |

---

## **3. Database Modeling**

### Entities

#### `users`
```sql
id          UUID (PK)
email       VARCHAR UNIQUE
password    VARCHAR
name        VARCHAR
created_at  TIMESTAMP
```

#### `teams`
```sql
id          UUID (PK)
name        VARCHAR
created_by  UUID (FK to users.id)
created_at  TIMESTAMP
```

#### `roles`
```sql
id    UUID (PK)
name  ENUM('admin', 'manager', 'member')
```

#### `team_members`
(*junction table: many-to-many between `users` and `teams`*)
```sql
id         UUID (PK)
user_id    UUID (FK to users.id)
team_id    UUID (FK to teams.id)
role_id    UUID (FK to roles.id)
joined_at  TIMESTAMP
is_active  BOOLEAN
```

---

## **4. Interconnections (ERD Overview)**

```
users ─────┐
           │       team_members ─────── roles
teams ─────┘               │
                           ▼
                        tasks
```

### Relationships
- One `team` has many `team_members`
- One `user` can be in many `team_members`
- Each `team_member` has one `role`
- Tasks are scoped to a team and optionally assigned to a member (user)

---

## **5. Permissions Logic (at runtime)**

### Access Control Middleware (example)
```go
func RequireRole(teamID, userID string, allowedRoles ...string) bool {
   member := GetTeamMember(teamID, userID)
   for _, role := range allowedRoles {
       if member.Role.Name == role {
           return true
       }
   }
   return false
}
```

logic to protect routes:

```go
if !RequireRole("team-123", currentUser.ID, "admin", "manager") {
   return c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permission"})
}
```

---

## **6. Use Case Scenarios**

### Creating a Team
- User sends `POST /teams` with team name.
- Backend creates a new team with the user as `admin` in `team_members`.

### Inviting a User
- `admin` or `manager` generates an invite token tied to a team.
- The invitee accepts and is added to `team_members` with a default role (e.g., `member`).

### Role Updates
- `admin` can promote/demote users within the same team.
- Update `team_members.role_id`.

### Deleting a Team
- Only `admin` can delete a team.
- Deleting a team also deletes:
  - `team_members`
  - `tasks`
  - `chat history`
  - (optional) `files` and `logs`

---

## **7. Sample API Endpoints**

| Method | Endpoint                      | Description                    |
|--------|-------------------------------|--------------------------------|
| POST   | `/api/teams`                  | Create new team                |
| GET    | `/api/teams/:id`              | Get team info                  |
| POST   | `/api/teams/:id/invite`       | Invite user to team            |
| PATCH  | `/api/teams/:id/role/:uid`    | Change user's role in team     |
| DELETE | `/api/teams/:id/members/:uid` | Remove user from team          |

---

## **8. Scalability Thoughts**

- Use a **Redis pub-sub** system to broadcast role changes in real-time (for chat/dashboard needs it).
- Add **team-specific settings** like file limits, task limits.
- Create **orgs/subteams** later by nesting teams if needed.

---
