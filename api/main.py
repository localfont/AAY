from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional
import sqlite3
import os

app = FastAPI(title="AAR - Alpine AUR Repository API", version="0.1.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

DB_PATH = os.getenv("AAR_DB", "aar.db")


def get_db():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    return conn


def init_db():
    db = get_db()
    db.execute("""
        CREATE TABLE IF NOT EXISTS packages (
            name        TEXT PRIMARY KEY,
            version     TEXT NOT NULL,
            description TEXT,
            maintainer  TEXT,
            url         TEXT,
            repo_url    TEXT,
            depends     TEXT DEFAULT '',
            makedepends TEXT DEFAULT '',
            votes       INTEGER DEFAULT 0,
            created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    """)
    db.execute("""
        CREATE TABLE IF NOT EXISTS comments (
            id         INTEGER PRIMARY KEY AUTOINCREMENT,
            pkgname    TEXT,
            author     TEXT,
            body       TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    """)
    db.commit()

    # Örnek paket ekle
    try:
        db.execute("""
            INSERT INTO packages (name, version, description, maintainer, url, repo_url, depends)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        """, (
            "hello-alpine",
            "1.0.0",
            "AAR test paketi - merhaba dünya",
            "admin",
            "https://github.com/alpineaar/hello-alpine",
            "https://github.com/alpineaar/hello-alpine.git",
            "bash"
        ))
        db.commit()
    except:
        pass


init_db()


class Package(BaseModel):
    name: str
    version: str
    description: str
    maintainer: str
    url: Optional[str] = ""
    repo_url: str
    depends: Optional[List[str]] = []
    makedepends: Optional[List[str]] = []


class Comment(BaseModel):
    author: str
    body: str


def row_to_dict(row):
    d = dict(row)
    d["depends"] = d.get("depends", "").split(",") if d.get("depends") else []
    d["makedepends"] = d.get("makedepends", "").split(",") if d.get("makedepends") else []
    return d


@app.get("/")
def root():
    return {"name": "AAR API", "version": "0.1.0", "status": "running"}


@app.get("/api/search")
def search(q: str):
    db = get_db()
    rows = db.execute(
        "SELECT * FROM packages WHERE name LIKE ? OR description LIKE ? ORDER BY votes DESC",
        (f"%{q}%", f"%{q}%")
    ).fetchall()
    return [row_to_dict(r) for r in rows]


@app.get("/api/info/{name}")
def info(name: str):
    db = get_db()
    row = db.execute("SELECT * FROM packages WHERE name = ?", (name,)).fetchone()
    if not row:
        raise HTTPException(status_code=404, detail="Paket bulunamadı")
    return row_to_dict(row)


@app.post("/api/submit")
def submit(pkg: Package):
    db = get_db()
    try:
        db.execute("""
            INSERT INTO packages (name, version, description, maintainer, url, repo_url, depends, makedepends)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        """, (
            pkg.name, pkg.version, pkg.description,
            pkg.maintainer, pkg.url, pkg.repo_url,
            ",".join(pkg.depends), ",".join(pkg.makedepends)
        ))
        db.commit()
        return {"status": "ok", "message": f"{pkg.name} eklendi"}
    except sqlite3.IntegrityError:
        raise HTTPException(status_code=409, detail="Bu paket zaten mevcut")


@app.put("/api/update/{name}")
def update(name: str, pkg: Package):
    db = get_db()
    row = db.execute("SELECT name FROM packages WHERE name = ?", (name,)).fetchone()
    if not row:
        raise HTTPException(status_code=404, detail="Paket bulunamadı")
    db.execute("""
        UPDATE packages SET version=?, description=?, maintainer=?,
        url=?, repo_url=?, depends=?, makedepends=? WHERE name=?
    """, (
        pkg.version, pkg.description, pkg.maintainer,
        pkg.url, pkg.repo_url,
        ",".join(pkg.depends), ",".join(pkg.makedepends), name
    ))
    db.commit()
    return {"status": "ok"}


@app.delete("/api/delete/{name}")
def delete(name: str):
    db = get_db()
    db.execute("DELETE FROM packages WHERE name = ?", (name,))
    db.commit()
    return {"status": "ok"}


@app.post("/api/vote/{name}")
def vote(name: str):
    db = get_db()
    db.execute("UPDATE packages SET votes = votes + 1 WHERE name = ?", (name,))
    db.commit()
    return {"status": "voted"}


@app.get("/api/comments/{name}")
def get_comments(name: str):
    db = get_db()
    rows = db.execute(
        "SELECT * FROM comments WHERE pkgname = ? ORDER BY created_at DESC", (name,)
    ).fetchall()
    return [dict(r) for r in rows]


@app.post("/api/comments/{name}")
def add_comment(name: str, comment: Comment):
    db = get_db()
    db.execute(
        "INSERT INTO comments (pkgname, author, body) VALUES (?, ?, ?)",
        (name, comment.author, comment.body)
    )
    db.commit()
    return {"status": "ok"}


@app.get("/api/stats")
def stats():
    db = get_db()
    pkg_count = db.execute("SELECT COUNT(*) FROM packages").fetchone()[0]
    top = db.execute(
        "SELECT name, votes FROM packages ORDER BY votes DESC LIMIT 5"
    ).fetchall()
    return {
        "total_packages": pkg_count,
        "top_voted": [dict(r) for r in top]
    }
