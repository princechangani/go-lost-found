# 🐛 Bug Report & 🚀 Feature Suggestions — go-lost-found

> Generated: March 7, 2026  
> Covers: Backend (Go/Fiber/MongoDB) + Frontend (React/TypeScript/Tailwind)

---

## 🐛 BUGS FOUND

---

### BACKEND

---

#### BUG-001 — `models.go`: `NewUser` uses wrong ObjectID format
**File:** `internal/models/models.go` — `NewUser()`  
**Severity:** 🔴 Critical

```go
// ❌ WRONG — primitive.NewObjectID().String() returns "ObjectID(\"...\")" not a hex ID
ID: primitive.NewObjectID().String(),

// ✅ FIX
ID: primitive.NewObjectID().Hex(),
```
**Impact:** The stored user ID will be `ObjectID("abc123...")` instead of a clean `abc123...` hex string, causing lookup failures in FindUserByEmail and JWT flows.

---

#### BUG-002 — `itemRepository.go` + `categoryRepository.go`: MongoDB queries use plain string `_id` instead of `primitive.ObjectID`
**Files:** `internal/repository/itemRepository.go`, `categoryRepository.go`, `contactRepository.go`, `userRepository.go`  
**Severity:** 🔴 Critical

```go
// ❌ WRONG — MongoDB stores _id as ObjectID, not plain string
bson.M{"_id": id}

// ✅ FIX — convert string to ObjectID first
oid, err := primitive.ObjectIDFromHex(id)
if err != nil { return error }
bson.M{"_id": oid}
```
**Impact:** `GetItemsByID`, `GetCategoryById`, `GetContactById`, `DeleteItems`, `UpdateItems`, `DeleteCategory`, `DeleteContact` will ALL return "no documents found" because the types mismatch.

---

#### BUG-003 — `itemRepository.go`: `CreateItems` sets ID incorrectly
**File:** `internal/repository/itemRepository.go` — `CreateItems()`  
**Severity:** 🔴 Critical

```go
// ❌ WRONG
if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
    item.ID = oid.String()  // returns "ObjectID(...)" not a hex string
}

// ✅ FIX
if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
    item.ID = oid.Hex()
}
```

---

#### BUG-004 — `categoryRepository.go`: Same `oid.String()` bug in `SaveCategory` and `SaveAllCategory`
**File:** `internal/repository/categoryRepository.go`  
**Severity:** 🔴 Critical

Same fix as BUG-003: `.String()` → `.Hex()`

---

#### BUG-005 — `item.go`: Image uploaded as `"image"` key in FormFile but frontend sends `"file"` key
**File:** `internal/handler/item.go` vs `frontend/src/pages/ReportItem.tsx`  
**Severity:** 🔴 Critical

```go
// Backend reads:
file, err := c.FormFile("image")

// Frontend sends:
formData.append("file", file);   // ← key is "file" not "image"
```
**Impact:** Images are NEVER saved. The backend silently skips image upload since `err != nil`.

---

#### BUG-006 — `item.go`: `DeleteItems` uses GET method in route definition
**File:** `main.go`  
**Severity:** 🟠 High

```go
// ❌ WRONG — DELETE action registered on GET
protected.Get("/:id", handler.DeleteItems)

// ✅ FIX
protected.Delete("/:id", handler.DeleteItems)
```
**Impact:** Browsing `/items/{id}` deletes the item! `GetItemsByID` and `DeleteItems` both respond to the same `GET /:id` route — only the first registered one (`GetItemsByID`) actually runs.

---

#### BUG-007 — `auth.go`: `Login` reuses `RegisterRequest` struct (includes `FcmToken` as required)
**File:** `internal/handler/auth.go`  
**Severity:** 🟠 High

```go
// ❌ Login handler uses RegisterRequest which has validate:"required" on FcmToken
func Login(c *fiber.Ctx) error {
    var req RegisterRequest  // FcmToken is "required" but login doesn't need it
```
**Impact:** Any login attempt without an `fcmToken` field would fail validation (if a validator middleware is added later). The struct re-use is semantically wrong. Create a dedicated `LoginRequest` struct.

---

#### BUG-008 — `contact.go`: `UpdateContact` calls `CreateContact` instead of an update function
**File:** `internal/handler/contact.go` — `UpdateContact()`  
**Severity:** 🟠 High

```go
// ❌ WRONG — this INSERTS a new contact instead of updating
err := repository.CreateContact(&contact)

// ✅ FIX — should call repository.UpdateContact(&contact) which doesn't exist yet
```
Also: there are **two consecutive `return c.JSON(contact)` statements** (unreachable code).

---

#### BUG-009 — `category.go`: `UpdateCategory` ignores the request body and fetches old data
**File:** `internal/handler/category.go` — `UpdateCategory()`  
**Severity:** 🟠 High

```go
// ❌ BUG — parses body into `category`, then immediately overwrites it with DB data
if err := c.BodyParser(&category); err != nil { ... }
category, err := repository.GetCategoryById(id)   // ← overwrites the parsed body!
category, err = repository.SaveCategory(category)  // ← saves old data, not new
```
**Impact:** `UpdateCategory` never actually updates anything — it just re-saves the existing record.

---

#### BUG-010 — `notification.go` (repository): `GetNotificationByID` returns wrong type
**File:** `internal/repository/notificationRepository.go`  
**Severity:** 🟡 Medium

```go
// ❌ WRONG — returns `error` but should return `(models.Notification, error)`
func GetNotificationByID(id string) error {
    var notification models.Notification
    err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&notification)
    return err  // notification is decoded but never returned!
}
```

---

#### BUG-011 — `utils/jwt.go`: JWT secret loaded at package init time (before `.env` is loaded)
**File:** `internal/utils/jwt.go`  
**Severity:** 🟠 High

```go
// ❌ WRONG — evaluated at package init, before godotenv.Load() in main()
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// ✅ FIX — use a function or lazy-load
func getJWTSecret() []byte {
    return []byte(os.Getenv("JWT_SECRET"))
}
```
**Impact:** `jwtSecret` is always empty `[]byte{}`, meaning all JWTs are signed with an empty secret — a severe security vulnerability.

---

#### BUG-012 — `main.go`: SMTP credentials hardcoded in source code
**File:** `main.go`  
**Severity:** 🔴 Critical (Security)

```go
// ❌ NEVER do this — email + app password in source code
services.NewEmailService("smtp.gmail.com", 465, "princechangani.dev@gmail.com", "mosg tvys putn ltts")
```
**Fix:** Read from environment variables (`os.Getenv("SMTP_USER")`, `os.Getenv("SMTP_PASS")`).

---

#### BUG-013 — `contact.go`: `CreateContact` handler doesn't return a proper success response
**File:** `internal/handler/contact.go`  
**Severity:** 🟡 Medium

```go
// ❌ Returns raw contact object with no status/statusCode/message envelope
return c.JSON(contact)

// ✅ FIX
return c.Status(201).JSON(fiber.Map{
    "status": true,
    "statusCode": "201",
    "message": "Contact created successfully",
    "data": contact,
})
```
Frontend checks `response.status === 200` but the backend returns 200 by default — fragile coupling.

---

#### BUG-014 — `main.go`: `contact` and `notification` routes are never registered
**File:** `main.go`  
**Severity:** 🔴 Critical

The handlers for `contact` (`CreateContact`, `GetContacts`, etc.) and `SendNotificationHandler` are **never wired to any route**. They exist but are unreachable.

---

#### BUG-015 — `item.go`: `CreateItems` sends notification failure as fatal (breaks item creation)
**File:** `internal/handler/item.go`  
**Severity:** 🟠 High

```go
// ❌ If notification fails, we return 500 — but the item was ALREADY saved!
err = services.SendMultiNotification(...)
if err != nil {
    return c.Status(500).JSON(...)  // item is saved but client gets error
}
```
**Fix:** Log notification errors, don't fail the entire request.

---

### FRONTEND

---

#### BUG-016 — `axiosInstance.ts`: Authorization header sends raw token without `Bearer ` prefix
**File:** `frontend/src/api/axiosInstance.ts`  
**Severity:** 🟠 High

```ts
// ❌ WRONG — sends the token as-is (no "Bearer " prefix)
config.headers.Authorization = `${token}`;

// ✅ FIX
config.headers.Authorization = `Bearer ${token}`;
```
Backend middleware does: `strings.TrimPrefix(authHeader, "Bearer ")` — if the interceptor doesn't add it, the raw token gets sent to the API but `ReportItem.tsx` manually adds `Authorization: localStorage.getItem('token')` (also missing `Bearer `), causing inconsistency.

---

#### BUG-017 — `App.tsx`: Route guard uses `location.pathname` from `window.location` (not React Router)
**File:** `frontend/src/App.tsx`  
**Severity:** 🟠 High

```tsx
// ❌ WRONG — `location` here is window.location, not React Router's location
if (location.pathname === '/' && !isAuthenticated) {
    return <AuthPage ... />;
}
```
This means the auth guard ONLY works on `/`. If the user visits `/lost-items` directly while unauthenticated, they see the page without auth. Use `<Navigate>` in route definitions or `useLocation()`.

---

#### BUG-018 — `App.tsx`: `isAuthenticated` state never updates on login (stale state)
**File:** `frontend/src/App.tsx`  
**Severity:** 🟠 High

The `isAuthenticated` state is only set once on mount. When `AuthPage` calls `onLoginSuccess`, it sets `isAuthenticated(true)` and navigates to `/home` — but the main App layout still won't render routes correctly because the auth check logic is broken (see BUG-017).

---

#### BUG-019 — `ReportItem.tsx`: `formData.append("file", file)` but backend reads `"image"`
**File:** `frontend/src/pages/ReportItem.tsx`  
**Severity:** 🔴 Critical  
*(Same as BUG-005 — confirmed from both sides)*

---

#### BUG-020 — `ReportItem.tsx`: `email` field is never appended to FormData
**File:** `frontend/src/pages/ReportItem.tsx`  
**Severity:** 🟠 High

```ts
// ❌ The backend reads c.FormValue("email") but frontend never sends it
formData.append("contact", contact);  // sent as "contact"
// Missing: formData.append("email", contact);
```

---

#### BUG-021 — `LostItems.tsx` / `FoundItems.tsx`: `Items` type used without importing
**File:** `frontend/src/pages/LostItems.tsx`, `FoundItems.tsx`, `ContactDialog.tsx`  
**Severity:** 🟡 Medium

`Items` and `ItemsResponse` interfaces are declared globally in `types/Item.ts` without `export`. This relies on ambient declaration which is fragile and non-standard. They should be exported and imported explicitly.

---

#### BUG-022 — `ContactDialog.tsx`: `itemId` not reset correctly in initial state
**File:** `frontend/src/components/ContactDialog.tsx`  
**Severity:** 🟡 Medium

```tsx
// ❌ itemId is set once on mount — if `item` prop changes, contactData.itemId won't update
const [contactData, setContactData] = React.useState<Contact>({
    itemId: item?.id ?? "",  // stale on re-render
    ...
});

// ✅ FIX — use useEffect to sync
useEffect(() => {
    setContactData(prev => ({ ...prev, itemId: item?.id ?? "" }));
}, [item]);
```

---

#### BUG-023 — `endpoints.ts`: `CREATE_ITEM` endpoint is `/items/create` but route is `POST /items`
**File:** `frontend/src/api/endpoints.ts` vs `main.go`  
**Severity:** 🔴 Critical

```ts
// ❌ Frontend calls /items/create
CREATE_ITEM: "/items/create"

// ✅ Backend registers:
protected.Post("", handler.CreateItems)  // → /api/v1/items
```
Fix: `CREATE_ITEM: "/items"`

---

#### BUG-024 — `Home.tsx`: Missing `React` import (uses JSX without import in non-automatic JSX transform)
**File:** `frontend/src/pages/Home.tsx`  
**Severity:** 🟡 Medium (depends on TSConfig/Vite settings)

```tsx
// ❌ Missing import
const Home: React.FC = () => { ... }  // uses React.FC without importing React
```

---

#### BUG-025 — `service/email.go`: Email body set as `text/plain` but HTML content is sent
**File:** `internal/service/email.go`  
**Severity:** 🟡 Medium

```go
// ❌ HTML body sent as plain text — tags will appear raw in email
msg.SetBody("text/plain", body)

// ✅ FIX
msg.SetBody("text/html", body)
```

---

## 🚀 NEW FEATURES TO ADD

---

### FEATURE-001 — 🔐 User Registration via Frontend
**Currently:** No signup UI exists — `AuthPage` only has login. A "Sign up" link exists but leads nowhere.  
**Add:** A register form/page with email, password, role selection, and FCM token support.  
**Files to create/edit:** `frontend/src/pages/RegisterPage.tsx`, `App.tsx`, `AuthPage.tsx`

---

### FEATURE-002 — 🖼️ Image Preview Before Upload
**Currently:** Only shows the filename after selection.  
**Add:** A thumbnail preview of the selected image in `ReportItem.tsx`.  
**Files:** `frontend/src/pages/ReportItem.tsx`

---

### FEATURE-003 — 📄 Item Detail Page
**Currently:** No dedicated detail page — items are shown in a card grid only.  
**Add:** A `/items/:id` route with full item details, larger image, and contact button.  
**Files to create:** `frontend/src/pages/ItemDetail.tsx`

---

### FEATURE-004 — 🔔 In-App Notification Bell
**Currently:** Notifications are only push (FCM). No in-app notification list.  
**Add:** A notification bell icon in Navbar showing recent notifications fetched from `/api/v1/notifications/all`.  
**Backend:** Add route `GET /api/v1/notifications/all` (handler and route already exist but not wired).  
**Files:** `GlassmorphismNavbar.tsx`, new `NotificationPanel.tsx`

---

### FEATURE-005 — 🔒 Protected Route Component
**Currently:** Auth guard is broken (BUG-017). Unauthenticated users can access all routes directly.  
**Add:** A `ProtectedRoute` wrapper component that redirects to `/` if no token is present.  
**Files to create:** `frontend/src/components/ProtectedRoute.tsx`

---

### FEATURE-006 — 🗑️ Delete / Edit Own Item
**Currently:** No UI for a user to manage their own reported items.  
**Add:** "My Items" page showing items reported by the logged-in user with edit/delete buttons.  
**Backend:** Filter `GET /items/all` by email from JWT, or add `GET /items/my`.  
**Files to create:** `frontend/src/pages/MyItems.tsx`

---

### FEATURE-007 — 📊 Dashboard / Statistics on Home Page
**Currently:** Home page is static with no live data.  
**Add:** Live counters — Total Lost, Total Found, Total Reunited — fetched from the API.  
**Files:** `frontend/src/pages/Home.tsx`

---

### FEATURE-008 — ♻️ Pagination / Infinite Scroll for Item Lists
**Currently:** All items are fetched and displayed at once — poor performance at scale.  
**Add:** Backend pagination (`?page=1&limit=10`) and frontend pagination controls.  
**Backend files:** `itemRepository.go`, `item.go`  
**Frontend files:** `LostItems.tsx`, `FoundItems.tsx`

---

### FEATURE-009 — 🌐 Rate Limiting
**Currently:** No rate limiting on any endpoint — vulnerable to brute force and spam.  
**Add:** `github.com/gofiber/fiber/v2/middleware/limiter` on auth and item creation routes.  
**Files:** `main.go`

---

### FEATURE-010 — ✅ Input Validation with `go-playground/validator`
**Currently:** Backend parses JSON/form data but doesn't validate fields (no min length, email format, etc.).  
**Add:** Struct validation tags are already partially defined but a validator is never called.  
**Files:** All handler files

---

### FEATURE-011 — 🗓️ Expiry / Auto-Archive of Old Items
**Currently:** Items stay in the database forever regardless of status.  
**Add:** A background goroutine (or cron job) that marks items older than 30 days as "archived".  
**Files to create:** `internal/service/archiver.go`

---

### FEATURE-012 — 🔍 Full-Text Search on Backend
**Currently:** Search is done entirely on the frontend by filtering already-fetched items.  
**Add:** `GET /items/all?q=wallet` MongoDB text-index search on backend.  
**Files:** `itemRepository.go`, `item.go`

---

### FEATURE-013 — 📍 Google Maps Location Picker
**Currently:** Location is a plain text input.  
**Add:** An embedded map (Google Maps or Leaflet.js) allowing users to pin exact location.  
**Files:** `frontend/src/pages/ReportItem.tsx`

---

### FEATURE-014 — 🔑 Refresh Token / Token Expiry Handling
**Currently:** JWT is valid for 24h. When it expires, API calls fail silently.  
**Add:** Axios response interceptor that detects 401 and redirects to login, plus refresh token support.  
**Files:** `frontend/src/api/axiosInstance.ts`, `internal/utils/jwt.go`

---

### FEATURE-015 — 📧 Email Confirmation on Registration
**Currently:** Users can register without verifying their email.  
**Add:** Send a verification email on register with a token link.  
**Files:** `internal/handler/auth.go`, `internal/service/email.go`

---

## 📋 PRIORITY SUMMARY

| Priority | Bug/Feature | File(s) | Status |
|----------|-------------|---------|--------|
| 🔴 Fix Now | BUG-001: ObjectID `.Hex()` | `models.go` | ✅ Fixed |
| 🔴 Fix Now | BUG-002: ObjectID in queries | All repositories | ✅ Fixed |
| 🔴 Fix Now | BUG-003: CreateItems InsertedID `.Hex()` | `itemRepository.go` | ✅ Fixed |
| 🔴 Fix Now | BUG-004: SaveCategory/SaveAllCategory `.Hex()` | `categoryRepository.go` | ✅ Fixed |
| 🔴 Fix Now | BUG-005/019: FormFile key mismatch | `item.go`, `ReportItem.tsx` | ✅ Fixed |
| 🔴 Fix Now | BUG-006: DELETE route uses GET | `main.go` | ✅ Fixed |
| 🔴 Fix Now | BUG-011: JWT secret loaded before .env | `jwt.go` | ✅ Fixed |
| 🔴 Fix Now | BUG-012: Hardcoded SMTP credentials | `main.go` | ✅ Fixed |
| 🔴 Fix Now | BUG-014: Contact/notification routes not registered | `main.go` | ✅ Fixed |
| 🔴 Fix Now | BUG-023: Wrong CREATE_ITEM endpoint URL | `endpoints.ts` | ✅ Fixed |
| 🟠 Fix Soon | BUG-007: LoginRequest struct | `auth.go` | ✅ Fixed |
| 🟠 Fix Soon | BUG-008: UpdateContact calls CreateContact | `contact.go` | ✅ Fixed |
| 🟠 Fix Soon | BUG-009: UpdateCategory ignores body | `category.go` | ✅ Fixed |
| 🟠 Fix Soon | BUG-015: Notification failure breaks item creation | `item.go` | ✅ Fixed |
| 🟠 Fix Soon | BUG-016: Missing Bearer prefix | `axiosInstance.ts` | ✅ Fixed |
| 🟠 Fix Soon | BUG-017/018: Broken auth guard + stale state | `App.tsx` | ✅ Fixed |
| 🟠 Fix Soon | BUG-020: email not sent in FormData | `ReportItem.tsx` | ✅ Fixed |
| 🟡 Improve | BUG-010: GetNotificationByID wrong return type | `notificationRepository.go` | ✅ Fixed |
| 🟡 Improve | BUG-013: CreateContact response envelope | `contact.go` | ✅ Fixed |
| 🟡 Improve | BUG-021: Items type not exported/imported | `Item.ts`, all pages | ✅ Fixed |
| 🟡 Improve | BUG-022: ContactDialog stale itemId | `ContactDialog.tsx` | ✅ Fixed |
| 🟡 Improve | BUG-024: Missing React import in Home.tsx | `Home.tsx` | ✅ Fixed |
| 🟡 Improve | BUG-025: Email body text/plain → text/html | `email.go` | ✅ Fixed |
| 🚀 Add | FEATURE-005: ProtectedRoute | New component | ⏳ Pending |
| 🚀 Add | FEATURE-001: Register Page | New page | ⏳ Pending |
| 🚀 Add | FEATURE-003: Item Detail Page | New page | ⏳ Pending |
| 🚀 Add | FEATURE-004: Notification Bell | Navbar | ⏳ Pending |
| 🚀 Add | FEATURE-007: Dashboard stats | `Home.tsx` | ⏳ Pending |

