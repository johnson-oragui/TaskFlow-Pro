

## **1. Core Concept of Real-time Chat (WebSockets)**

Real-time chat functionality allows users to communicate instantly within teams, making it ideal for task collaboration.
The goal is to implement **WebSocket-based communication** that can facilitate seamless, bidirectional messaging between users in real time.
Messages can be tied to tasks, teams, or general team discussions.

### **Purpose of Real-time Chat:**
- **Instant Communication**: Enable users to message each other or communicate within teams regarding tasks, feedback, or project updates.
- **Task Collaboration**: Associate messages with specific tasks or projects, providing clear context.
- **Seamless Interaction**: No need for users to refresh their browser or make manual requests; the app stays updated in real-time.
- **Group Chats**: Allow team members to join or leave channels or rooms based on team roles, enabling chat specific to teams, projects, or tasks.

---

## **2. WebSocket Connections**

WebSocket is a communication protocol that allows bidirectional communication between the client and server.
In this case, WebSocket will allow for continuous communication as users send or receive messages without the need to refresh or make new HTTP requests.

### **WebSocket Flow:**
1. **Client Connection**: The client initiates a WebSocket connection to the backend. This is a long-lived connection that allows both the client and server to send messages anytime.
2. **Authentication**: The client sends an authentication token (e.g., JWT) with the WebSocket connection request to authenticate the user.
3. **Create/Join Chat Room**: Once authenticated, the user can join a room (either based on a specific task, project, or team).
4. **Send/Receive Messages**: Users can send and receive messages in real-time within the room.
5. **Close Connection**: When the client disconnects or the session ends, the WebSocket connection is closed.

### WebSocket URL:
```plaintext
ws://127.0.0.1:7006/ws/chat/{teamId}/{taskId}
```

- **teamId**: The identifier for the team where the chat is happening.
- **taskId** (optional): The identifier for the task associated with the conversation.

---

## **3. Database Modeling for Chat Messages**

The messages exchanged over WebSocket can be saved into a database for persistence (e.g., showing historical messages or tracking activity).
A simple database model can store messages, users, and associated tasks or teams.

### **Message Table**
```sql
messages
---------
id              UUID (PK)
team_id         UUID (FK to teams)
task_id         UUID (FK to tasks, nullable)
user_id         UUID (FK to users)
message         TEXT
timestamp       TIMESTAMP
is_edited       BOOLEAN
edited_at       TIMESTAMP (nullable)
```

### **Columns Explanation:**
- **id**: Unique identifier for each message.
- **team_id**: The team the message belongs to.
- **task_id**: (Optional) The task associated with the message.
- **user_id**: The user who sent the message.
- **message**: The content of the chat message.
- **timestamp**: When the message was sent.
- **is_edited**: A flag to indicate if the message has been edited.
- **edited_at**: The timestamp when the message was edited.

### **Additional Relationships:**
- Messages are sent by a **user** and belong to a **team**.
- Messages may be associated with a **task** (optional).

---

## **4. WebSocket Handling in the Backend**

### **WebSocket Connection Flow:**
1. **Connect**: When a user connects via WebSocket, the backend validates the JWT token to authenticate the user.
2. **Join Room**: Based on the team and optional task, the user joins a specific WebSocket room to receive messages for that team or task.
3. **Send Message**: When a user sends a message, the server broadcasts that message to the connected users in the corresponding room.
4. **Store Message**: The backend saves the message in the database to keep a record of all communications.
5. **Receive Message**: All connected clients (users in the room) receive the message in real-time.
6. **Close Connection**: When the client disconnects, the WebSocket connection is closed.
---

## **5. Use Case Scenarios for Real-time Chat**

### **Message Sent to Task or Team**
- **Event**: A user sends a message regarding a task (e.g., asking a question about the task or providing an update).
- **Data captured**:
  - **Message Content**: The content of the message.
  - **User**: The user who sent the message.
  - **Task/Team**: The task or team the message is related to.
  - **Timestamp**: When the message was sent.
- **Backend stores**:
  - A new record is created in the `messages` table.
  - The message is broadcasted to all users in the corresponding WebSocket room (team/task).

**Example Response**:
```json
{
  "user_id": "user-456",
  "team_id": "team-123",
  "task_id": "task-789",
  "content": "The task needs more details on the design. Can someone update it?",
  "timestamp": "2025-04-09T14:00:00Z"
}
```

---

### **Edit a Message**
- **Event**: A user edits a previously sent message (e.g., correcting a typo or adding additional information).
- **Data captured**:
  - **Message ID**: The identifier of the message being edited.
  - **New Content**: The updated message content.
  - **Edited Timestamp**: When the message was edited.
- **Backend stores**:
  - The `is_edited` flag is set to `true`.
  - The `edited_at` timestamp is updated.

---

## **6. API Endpoints for Real-time Chat**

| Method | Endpoint                           | Description                                  |
|--------|------------------------------------|----------------------------------------------|
| GET    | `/ws/chat/{teamId}/{taskId}`       | Establish WebSocket connection for chat      |
| GET    | `/tasks/{taskId}/messages`         | Retrieve all historical messages for a task  |
| POST   | `/tasks/{taskId}/messages`         | Send a new message to a task (via HTTP)      |

---

## **7. Real-time Messaging Flow Example**

**Backend Response**: A WebSocket connection is established and the client is subscribed to the `team-123/task-789` room.

### Sending a Message:
```json
{
  "user_id": "user-456",
  "team_id": "team-123",
  "task_id": "task-789",
  "content": "Let's prioritize this task next week.",
  "timestamp": "2025-04-09T14:00:00Z"
}
```

**Backend Response**: The message is broadcasted to all connected clients in the same room (team or task).

---

## **8. Permissions and Access Control**

Access control for chat messages can be managed based on **user roles**:
- **Admins**: Can send, edit, and delete messages in any team/task.
- **Managers**: Can send and edit messages within their team or associated tasks.
- **Members**: Can only send messages to tasks they are assigned to.

---

## **9. Example Real-time Chat Scenarios**

### **User sends a message to a task**
1. **User** connects to the WebSocket and joins a room for the team and task.
2. **User** sends a message regarding the task.
3. **Backend** broadcasts the message to all connected clients in that room.
4. The message is stored in the database for future reference.
