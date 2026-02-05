## 🐍 Opgradering af Backend (Python 2 til 3)

Som en del af moderniseringen af projektet har vi opgraderet kildekoden fra Python 2 til Python 3. Dette sikrer, at vi overholder de nyeste standarder for sikkerhed, performance og syntaks.

### 🛠️ Gennemførelse
Opgraderingen blev udført automatisk ved hjælp af værktøjet `2to3`, som håndterer oversættelsen af ældre syntaks (f.eks. `print`-statements og `import`-logik) til den moderne Python 3-standard.

**Kommando anvendt:**
```bash
2to3 -w app.py