## 🐍 Opgradering af Backend (Python 2 til 3)

Som en del af moderniseringen af projektet har vi opgraderet kildekoden fra Python 2 til Python 3. Dette sikrer, at vi overholder de nyeste standarder for sikkerhed, performance og syntaks.

### 🛠️ Gennemførelse
Opgraderingen blev udført automatisk ved hjælp af værktøjet `2to3`, som håndterer oversættelsen af ældre syntaks (f.eks. `print`-statements og `import`-logik) til den moderne Python 3-standard.

**Kommando anvendt:**
```bash
2to3 -w app.py
```


## Project Update: Migration & Framework Integration
**Date:** 05/02/2026

---

### Step 1: Structural Refactoring
The project architecture has been reorganized to accommodate a dual-stack environment. This refactor establishes a clear separation between the existing codebase and the new, optimized logic written in **Go**, allowing for a phased migration without service interruption.

* **Legacy Logic:** Isolated existing modules to maintain stability while the transition is underway.
* **Refined Logic:** Established a dedicated space for high-performance Go implementations.

### Step 2: Web Framework Integration
To support the refined services, the **Gorilla Mux** toolkit has been integrated as the primary web routing layer.

* **Dependency Added:** `github.com/gorilla/mux`
* **Purpose:** Leveraged Gorilla Mux for its robust routing capabilities, including advanced pattern matching and middleware support, which will facilitate the handling of complex API endpoints in the refined code.

---

### Summary of Changes
| Feature | Action | Status |
| :--- | :--- | :--- |
| **Project Layout** | Segregated Legacy and Refined codebases | Completed |
| **Go Modules** | Initialized workspace and dependencies | Completed |
| **Web Framework** | Integrated `gorilla/mux` | Installed |


Malthe - oprettet epics og uploaded WhoKnows_Go_Gorilla_ProjectPlan.pdf til projektet. 
Ligger under Devops_Gorillas/Files/WhoKnows_Go_Gorilla_ProjectPlan.pdf

## Project Update: Go Backend Restructure & Endpoints
**Date:** 09/02/2026

---

### Step 1: Project Structure Update
We reorganized the Go project to a more maintainable structure in `main.go`, separating concerns for pages, handlers, models, and services. The new layout improves readability and prepares the codebase for additional features.

### 16/02/2026
## ☁️ Azure Deployment

The project is deployed on a Microsoft Azure virtual machine with the following configuration:

| Property      | Value              |
| ------------- | ------------------ |
| **Host Name** | `gorilla-maskinen` |
| **Public IP** | `51.120.83.21`     |
| **Platform**  | Microsoft Azure    |

### 👤 Available Users

Currently, only the following user accounts are provisioned on the server:

* `azureuser`
* `gustav`

### 🔐 Access

Connect via SSH:

```bash
ssh azureuser@51.120.83.21
```

> Note: Update this section if additional users or infrastructure changes are made.
