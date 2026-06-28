# Daylik — Product Overview

Daylik is a gamified habits tracker. Users build daily habits, track quantitative progress,
earn EXP and achievements. The core loop is:
build a habit → log daily progress → earn EXP → level up → unlock achievements.

---

## Core Concepts

### Habits

- Each user creates their own habits
- Habits are **quantitative** — the user sets a daily target value and logs a number each day
(e.g. "Run 5km", log 3.2km today)
- Habits have a **difficulty tier**: Easy / Medium / Hard
- Difficulty maps to a fixed EXP reward per completion (easy is 100 - 200, medium is 200-400, hard is 500 and more)
- Habits belong to a **category** (e.g. Health, Learning, Mindfulness, Productivity)
- Habits can be **deleted** — soft-deleted in the DB, history is preserved
- When creating a habit, the user sets a daily EXP reward. The total EXP cap across all habits is **1000/day**
(e.g. 3 habits: 300 + 600 + 100 = 1000)

### Daily Progress

- Each day the user logs their progress value per habit
- A habit is considered **completed** for the day when logged value >= daily target
- There is a clear **"day complete"** state when all habits are done — triggers a visual reward or bonus EXP
- A **weekly summary** shows total completions vs. total possible (e.g. "18/21 habits this week")

### Streaks

- A streak counts consecutive days where all habits were completed
- Users get a **grace period**: one "freeze" per week that prevents a streak from breaking on a missed day
- Dashboard shows **current streak** and **longest streak ever**

### EXP & Levels

- Users earn EXP by completing habits daily
- EXP accumulates into a **level** via a leveling curve
- Leveling curve is **linear** — each level requires a fixed amount of EXP (e.g. 1000 EXP per level)
  - This keeps progression feeling consistent at all levels and avoids frustration at level 50-100+
  - Exponential curves are intentionally avoided

### Achievements

- Achievements are **global** — defined by the app, the same for all users
- Earning an achievement grants a fixed EXP bonus
- Early achievements (reachable in week 1) are critical for new user retention
- Examples:
  - "First Step" — complete a habit for the first time
  - "On a Roll" — 5-day streak
  - "Unstoppable" — 30-day streak
  - "Century" — complete 100 habits total
  - "Full Week" — complete all habits every day for a full week
  - "Health Nut" — complete 10 habits in the Health category

---

## Features

### Dashboard

- GitHub-style contribution grid — one cell per day, colored by completion rate
- Prominent display of: current streak, longest streak ever, total EXP, current level
- Per-habit stats available on drill-down

### Weekly Recap

- Delivered as an in-app notification (push notification in the future)
- Shows: habits completed / total possible, EXP earned this week, streak status, achievements unlocked
