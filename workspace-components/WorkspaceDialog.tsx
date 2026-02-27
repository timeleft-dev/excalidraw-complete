import React, { useEffect, useState } from "react";

interface WorkspaceFileMeta {
  id: string;
  name: string;
  updated_at: string;
  updated_by: string;
}

interface Props {
  currentFileId: string | null;
  currentFileName: string | null;
  onLoad: (id: string, name: string) => void;
  onClose: () => void;
}

const fmt = (iso: string) => {
  try {
    return new Date(iso).toLocaleString(undefined, {
      dateStyle: "medium",
      timeStyle: "short",
    });
  } catch {
    return iso;
  }
};

export const WorkspaceDialog: React.FC<Props> = ({
  currentFileId,
  currentFileName,
  onLoad,
  onClose,
}) => {
  const [files, setFiles] = useState<WorkspaceFileMeta[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deleting, setDeleting] = useState<string | null>(null);

  useEffect(() => {
    fetch("/api/workspace/files")
      .then((r) => {
        if (!r.ok) throw new Error("Failed to load workspace");
        return r.json();
      })
      .then((data) => {
        setFiles(data);
        setLoading(false);
      })
      .catch((e) => {
        setError(e.message);
        setLoading(false);
      });
  }, []);

  const handleDelete = async (id: string, name: string) => {
    if (!window.confirm(`Delete "${name}"?`)) return;
    setDeleting(id);
    try {
      await fetch(`/api/workspace/files/${id}`, { method: "DELETE" });
      setFiles((prev) => prev.filter((f) => f.id !== id));
    } catch {
      alert("Failed to delete file.");
    } finally {
      setDeleting(null);
    }
  };

  // Detect dark mode from the html element class Excalidraw sets
  const isDark = document.documentElement.classList.contains("dark");
  const bg = isDark ? "#1e1e2e" : "#ffffff";
  const border = isDark ? "#3c3c4e" : "#e0e0e0";
  const text = isDark ? "#e0e0e0" : "#1a1a1a";
  const subtext = isDark ? "#888" : "#666";
  const rowHover = isDark ? "#2a2a3e" : "#f5f5f5";
  const btnPrimary = "#6965db";
  const btnDanger = isDark ? "#5c3a3a" : "#fee2e2";
  const btnDangerText = isDark ? "#ff8a8a" : "#dc2626";

  return (
    <div
      style={{
        position: "fixed",
        inset: 0,
        background: "rgba(0,0,0,0.5)",
        zIndex: 9999,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <div
        style={{
          background: bg,
          border: `1px solid ${border}`,
          borderRadius: 8,
          width: "min(640px, 92vw)",
          maxHeight: "80vh",
          display: "flex",
          flexDirection: "column",
          boxShadow: "0 20px 60px rgba(0,0,0,0.4)",
          color: text,
          fontFamily: "inherit",
        }}
      >
        {/* Header */}
        <div
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            padding: "16px 20px",
            borderBottom: `1px solid ${border}`,
          }}
        >
          <h2 style={{ margin: 0, fontSize: 18, fontWeight: 600 }}>
            Workspace
          </h2>
          <button
            onClick={onClose}
            style={{
              background: "none",
              border: "none",
              cursor: "pointer",
              color: subtext,
              fontSize: 20,
              lineHeight: 1,
              padding: "2px 6px",
            }}
            aria-label="Close"
          >
            ×
          </button>
        </div>

        {/* Current file indicator */}
        {currentFileName && (
          <div
            style={{
              padding: "10px 20px",
              borderBottom: `1px solid ${border}`,
              fontSize: 13,
              color: subtext,
            }}
          >
            Currently open:{" "}
            <strong style={{ color: text }}>{currentFileName}</strong>
          </div>
        )}

        {/* File list */}
        <div style={{ overflowY: "auto", flex: 1 }}>
          {loading && (
            <p style={{ textAlign: "center", padding: 32, color: subtext }}>
              Loading…
            </p>
          )}
          {error && (
            <p style={{ textAlign: "center", padding: 32, color: "#ef4444" }}>
              {error}
            </p>
          )}
          {!loading && !error && files.length === 0 && (
            <p style={{ textAlign: "center", padding: 32, color: subtext }}>
              No files saved yet. Use{" "}
              <strong>Save to Workspace</strong> to save your current drawing.
            </p>
          )}
          {!loading &&
            !error &&
            files.map((file) => (
              <div
                key={file.id}
                style={{
                  display: "flex",
                  alignItems: "center",
                  padding: "12px 20px",
                  borderBottom: `1px solid ${border}`,
                  gap: 12,
                  background:
                    file.id === currentFileId ? rowHover : "transparent",
                }}
                onMouseEnter={(e) =>
                  ((e.currentTarget as HTMLDivElement).style.background =
                    rowHover)
                }
                onMouseLeave={(e) =>
                  ((e.currentTarget as HTMLDivElement).style.background =
                    file.id === currentFileId ? rowHover : "transparent")
                }
              >
                {/* File info */}
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div
                    style={{
                      fontWeight: 500,
                      fontSize: 14,
                      whiteSpace: "nowrap",
                      overflow: "hidden",
                      textOverflow: "ellipsis",
                    }}
                  >
                    {file.name}
                    {file.id === currentFileId && (
                      <span
                        style={{
                          marginLeft: 8,
                          fontSize: 11,
                          color: btnPrimary,
                          fontWeight: 400,
                        }}
                      >
                        (open)
                      </span>
                    )}
                  </div>
                  <div style={{ fontSize: 12, color: subtext, marginTop: 2 }}>
                    {fmt(file.updated_at)}
                    {file.updated_by && ` · ${file.updated_by}`}
                  </div>
                </div>

                {/* Actions */}
                <button
                  onClick={() => {
                    onLoad(file.id, file.name);
                    onClose();
                  }}
                  style={{
                    padding: "5px 12px",
                    borderRadius: 6,
                    border: "none",
                    background: btnPrimary,
                    color: "#fff",
                    fontSize: 13,
                    cursor: "pointer",
                    fontWeight: 500,
                    flexShrink: 0,
                  }}
                >
                  Open
                </button>
                <button
                  onClick={() => handleDelete(file.id, file.name)}
                  disabled={deleting === file.id}
                  style={{
                    padding: "5px 10px",
                    borderRadius: 6,
                    border: "none",
                    background: btnDanger,
                    color: btnDangerText,
                    fontSize: 13,
                    cursor: deleting === file.id ? "not-allowed" : "pointer",
                    fontWeight: 500,
                    flexShrink: 0,
                  }}
                >
                  {deleting === file.id ? "…" : "Delete"}
                </button>
              </div>
            ))}
        </div>

        {/* Footer */}
        <div
          style={{
            padding: "12px 20px",
            borderTop: `1px solid ${border}`,
            fontSize: 12,
            color: subtext,
          }}
        >
          Use <strong>Save to Workspace</strong> in the menu to save the current
          drawing.
        </div>
      </div>
    </div>
  );
};
