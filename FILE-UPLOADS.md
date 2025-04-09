
## **1. Core Concept of File Uploads**

File uploads are a common feature in many applications, allowing users to attach documents, images, or other files to tasks, projects, or team discussions.
In the **TaskFlow** system, files will be uploaded and associated with tasks, and each file will have a record of metadata (e.g., uploader, file type, size) for management.

### **Purpose of File Uploads:**
- **Collaboration:** Allow team members to share files related to tasks.
- **Storage:** Attachments such as images, documents, spreadsheets, or PDFs can be stored and retrieved.
- **Tracking:** Keeping track of who uploaded which files and when can help with auditing and accountability.
- **Access Control:** Only authorized users can view, download, or delete files based on roles.

---

## **2. Database Modeling for File Uploads**

### File Upload Table
```sql
files
---------
id              UUID (PK)
task_id         UUID (FK to tasks, optional)
team_id         UUID (FK to teams)
user_id         UUID (FK to users) -- User who uploaded the file
file_name       VARCHAR
file_type       VARCHAR  -- mime type, e.g., 'image/png', 'application/pdf'
file_size       BIGINT   -- Size in bytes
file_url        VARCHAR  -- URL to access the file
created_at      TIMESTAMP
updated_at      TIMESTAMP
```

### **Columns Explanation:**
- **id**: A unique identifier for the file.
- **task_id**: (Optional) The task this file is related to.
- **team_id**: The team the file belongs to.
- **user_id**: The user who uploaded the file.
- **file_name**: The name of the file (may be stored with its original name or a sanitized name).
- **file_type**: The MIME type of the file (e.g., `application/pdf`, `image/png`).
- **file_size**: The size of the file in bytes (can help with storage planning or imposing file size limits).
- **file_url**: The URL where the file can be accessed or downloaded.
- **created_at**: Timestamp when the file was uploaded.
- **updated_at**: Timestamp of the last modification, if any (e.g., if the file is replaced).

---

## **3. Relationships**

```
users ─────────────────────────────┐
                                    ▼
                             files ──▶ tasks
                                    │
                                    ▼
                                teams
```

### Relationships:
- A **user** uploads **files** (one-to-many relationship).
- A **file** is related to one **task** (many-to-one).
- A **file** belongs to one **team** (many-to-one).
- **File metadata** (e.g., size, type) helps organize and identify each file, while the **file_url** points to the actual location of the file.

---

## **4. File Handling in the Backend**

The process of handling file uploads involves several steps, from accepting the file, validating it, saving it to storage (local or cloud), and storing metadata in the database.

### **File Upload Process:**
1. **Receive the File**:
   - The client sends the file via a `multipart/form-data` HTTP POST request.
2. **Validate the File**:
   - Ensure the file meets size limits and the allowed types (e.g., only images, PDFs, etc.).
3. **Save the File**:
   - If using **local storage**, the file is saved to the server's file system.
   - If using **cloud storage (e.g., AWS S3, Google Cloud Storage, cloudinary)**, the file is uploaded to the cloud, and the URL is retrieved.
4. **Store Metadata**:
   - The database entry is created for the file, including its URL, type, size, and relationship to a task and user.
5. **Return a Response**:
   - The file's metadata (e.g., URL) is returned to the client for reference.

### **Cloud Storage Example (AWS S3)**:
- If using **AWS S3** for file storage:
  - The backend would use the AWS SDK to upload the file to an S3 bucket.
  - After the upload, the file's URL in the S3 bucket is stored in the `file_url` column in the database.

### **Local File Storage**:
- For **local storage**, files can be stored in a designated folder on the server, with a reference to that path stored in the database.
- File system management tools can be used for handling file naming, access permissions, and folder structures.

---

## **5. Use Case Scenarios for File Uploads**

### **File Upload to Task (task-related)**

- **Event**: A user uploads a file related to a task (e.g., a design document or an analysis report).
- **Data captured**:
  - **User**: The user uploading the file.
  - **Task**: The task the file is associated with.
  - **File**: The file's metadata (name, type, size, URL).
- **Example Response**:
  ```json
  {
    "file_id": "file-123",
    "file_name": "design_document.pdf",
    "file_url": "https://s3.amazonaws.com/bucket/task-123/design_document.pdf",
    "file_size": 1024,
    "file_type": "application/pdf",
    "task_id": "task-123",
    "uploaded_by": "user-456",
    "created_at": "2025-04-09T14:00:00Z"
  }
  ```

---

### **File Replacement/Update**

- **Event**: A user replaces an existing file on a task (e.g., uploading a newer version of a design).
- **Data captured**:
  - **User**: The user replacing the file.
  - **Task**: The task to which the file is attached.
  - **Old File**: The previous file (can be archived or deleted).
  - **New File**: The new file's metadata.
- **Example Response**:
  ```json
  {
    "file_id": "file-124",
    "file_name": "updated_design_document.pdf",
    "file_url": "https://s3.amazonaws.com/bucket/task-123/updated_design_document.pdf",
    "file_size": 2048,
    "file_type": "application/pdf",
    "task_id": "task-123",
    "uploaded_by": "user-456",
    "created_at": "2025-04-10T14:00:00Z"
  }
  ```

---

## **6. API Endpoints for File Uploads**

| Method | Endpoint                           | Description                                  |
|--------|------------------------------------|----------------------------------------------|
| POST   | `/tasks/:taskId/upload`            | Upload a file for a specific task            |
| GET    | `/files/:fileId`                   | Retrieve file metadata (URL, size, etc.)     |
| DELETE | `/files/:fileId`                   | Delete a file (from server or cloud)         |
| GET    | `/tasks/:taskId/files`             | Retrieve all files uploaded for a task       |

---

## **7. File Upload API Flow Example**

### Upload File to Task
```json
POST /tasks/task-123/upload
Content-Type: multipart/form-data

{
    "file": <binary file data>
}
```

**Backend Response**:
```json
{
    "file_id": "file-123",
    "file_name": "project_plan.pdf",
    "file_url": "https://s3.amazonaws.com/bucket/task-123/project_plan.pdf",
    "file_size": 2048,
    "file_type": "application/pdf",
    "task_id": "task-123",
    "uploaded_by": "user-456",
    "created_at": "2025-04-09T14:00:00Z"
}
```

### Delete a File
```json
DELETE /files/file-123
```

**Backend Response**:
```json
{
    "message": "File deleted successfully."
}
```

---

## **8. File Permissions and Access Control**

**Permissions** for file access can be set based on user roles:
- **Admins** can delete, replace, and view all files across teams.
- **Managers** can only view and delete files within their team.
- **Members** can only view files attached to tasks they are working on.

Example middleware for permission check (Go):
```go
func CanUserDeleteFile(userID, fileID string) bool {
    file := GetFileByID(fileID)
    return file.UserID == userID || IsUserAdmin(userID)
}
```

---

## **9. Real-Time File Notifications**

- **WebSockets** or **Server-Sent Events (SSE)** can notify users when a new file has been uploaded or when an existing file is updated.
  
Example:
```go
go func() {
    for {
        select {
        case fileUpload := <-fileUploadChannel:
            // Notify users in the task's team about the new file
            socket.SendToUsersInTeam(fileUpload.TeamID, fileUpload)
        }
    }
}()
```

---

## **10. Example File Upload API Calls**

### File Upload to Task
```json
POST /tasks/123/upload
{
    "file": <binary file data>
}
```

### File Deletion
```json
DELETE /files/123
```
