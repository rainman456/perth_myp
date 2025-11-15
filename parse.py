import os
from pathlib import Path
from datetime import datetime
import mimetypes
import argparse
import sys

OUTPUT_BASE = "code_chunk"
MAX_CHUNK_CHARS = 400000  # Approximate safe limit for ~100k tokens (adjust if needed)
MAX_FILE_PART_CHARS = 300000  # Leave room for headers in chunks

# Extended ignore directories and files
IGNORE_DIRS = {
    ".anchor",
    ".git",
    "node_modules",
    "target",
    "__pycache__",
    ".idea",
    ".vscode",
    "venv",
    "dist",
    "build",
    "test-ledger",
    "public",
    #"api",
    "bank",
    "config",
    "middleware",
    "tests",
    "utils"
}
IGNORE_FILES = {
    ".DS_Store",
    "Thumbs.db",
    ".gitignore",
    ".prettierignore",
    "cargo.bak",
    "Cargo.lock",
    "yarn.lock",
    "README.md",
    "Dockerfile",
    "parse.py",
    "package-lock.json",
    "go.sum",
    "go.mod",
    "swagger.yaml",
    "internal.md",
    "setup.sh",
    "start-validator.sh",
    "CONTRIBUTING.md",
    "DEPLOYMENT.md",
    "LICENSE",
    "pnpm-lock.yaml",
    "SETUP.md",
    "update.txt",
    "banks.json"
}


def generate_tree(root: Path, prefix="", depth=0, max_depth=5):
    """Generate a tree-like structure for the given folder, skipping ignored dirs."""
    if depth > max_depth:
        return [prefix + "└── ... (max depth reached)"]

    entries = sorted(
        [
            e
            for e in root.iterdir()
            if e.name not in IGNORE_DIRS and e.name not in IGNORE_FILES
        ],
        key=lambda e: (e.is_file(), e.name.lower()),
    )
    result = []
    for i, entry in enumerate(entries):
        connector = "└── " if i == len(entries) - 1 else "├── "
        result.append(f"{prefix}{connector}{entry.name}{'/' if entry.is_dir() else ''}")
        if entry.is_dir():
            extension = "    " if i == len(entries) - 1 else "│   "
            result.extend(
                generate_tree(entry, prefix + extension, depth + 1, max_depth)
            )
    return result


def detect_content_type(file_path: Path):
    """Enhanced content type detection using both extension and mimetypes."""
    extension_mapping = {
        ".rs": "rust",
        ".go": "go",
        ".py": "python",
        ".js": "javascript",
        ".ts": "typescript",
        ".java": "java",
        ".c": "c",
        ".cpp": "cpp",
        ".cs": "csharp",
        ".rb": "ruby",
        ".php": "php",
        ".html": "html",
        ".css": "css",
        ".scss": "scss",
        ".toml": "toml",
        ".json": "json",
        ".md": "markdown",
        ".txt": "plaintext",
        ".sh": "shell",
        ".yml": "yaml",
        ".yaml": "yaml",
        ".sql": "sql",
        ".xml": "xml",
        ".dockerfile": "dockerfile",
        ".ipynb": "json",
    }

    language = extension_mapping.get(file_path.suffix.lower())
    if language:
        return language

    mime_type, _ = mimetypes.guess_type(file_path)
    return mime_type.split('/')[-1] if mime_type else "plaintext"


def get_file_stats(file_path: Path):
    """Get file statistics including size and modification time."""
    try:
        stats = file_path.stat()
        content = file_path.read_text(encoding="utf-8", errors="ignore")
        return {
            "size": f"{stats.st_size / 1024:.2f} KB",
            "modified": datetime.fromtimestamp(stats.st_mtime).strftime(
                "%Y-%m-%d %H:%M:%S"
            ),
            "lines": len(content.splitlines()),
        }
    except Exception:
        return None


def split_content(content: str, max_chars: int):
    """Split content into parts based on character limit, preserving lines."""
    parts = []
    lines = content.splitlines(True)  # Preserve newlines
    current = []
    current_size = 0
    for line in lines:
        line_len = len(line)
        if current_size + line_len > max_chars and current:
            parts.append(''.join(current))
            current = []
            current_size = 0
        current.append(line)
        current_size += line_len
    if current:
        parts.append(''.join(current))
    return parts


def parse_codebase(root_folder: str, output_base: str = OUTPUT_BASE):
    """Parse codebase and generate markdown in chunks."""
    try:
        root = Path(root_folder).resolve()
        if not root.exists():
            raise FileNotFoundError(f"Folder not found: {root_folder}")
        if not root.is_dir():
            raise NotADirectoryError(f"Path is not a directory: {root_folder}")

        # Generate tree structure
        tree_lines = [f"📁 {root.name}"] + generate_tree(root)
        tree_structure = f"```tree\n{'\n'.join(tree_lines)}\n```"

        # Collect file paths and stats
        file_paths = []
        file_count = 0
        total_size = 0.0
        file_types = set()

        for file_path in root.rglob("*"):
            if any(
                part in IGNORE_DIRS or part in IGNORE_FILES for part in file_path.parts
            ):
                continue

            if file_path.is_file():
                file_stats = get_file_stats(file_path)
                if file_stats:
                    file_count += 1
                    total_size += float(file_stats["size"].split()[0])
                    file_types.add(file_path.suffix)
                    file_paths.append(file_path)

        # Build summary
        summary = "\n".join([
            "\n---\n## 📊 Summary\n",
            f"- Total files: {file_count}\n",
            f"- Total size: {total_size:.2f} KB\n",
            f"- File types: {', '.join(sorted(ft or 'unknown' for ft in file_types))}\n",
        ])

        # Header
        header = "\n".join([
            f"# Codebase Analysis: {root.name}\n",
            f"Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n",
            "---\n\n## 📂 Project Structure\n",
            tree_structure,
            "\n---\n\n## 📄 File Contents\n",
        ])

        # Now build file md parts
        file_md_parts = []
        for file_path in file_paths:
            try:
                content = file_path.read_text(encoding="utf-8", errors="ignore")
                language = detect_content_type(file_path)
                file_stats = get_file_stats(file_path)  # Already computed, but recompute for simplicity

                relative_path = file_path.relative_to(root)

                file_header = f"### {relative_path}\n"
                if file_stats:
                    file_header += f"- Size: {file_stats['size']}\n"
                    file_header += f"- Lines: {file_stats['lines']}\n"
                    file_header += f"- Last Modified: {file_stats['modified']}\n"

                content_parts = split_content(content, MAX_FILE_PART_CHARS)
                for i, part_content in enumerate(content_parts, 1):
                    if len(content_parts) > 1:
                        part_header = f"{file_header} (Part {i}/{len(content_parts)})\n" if i == 1 else f"### {relative_path} (continued, part {i}/{len(content_parts)})\n"
                    else:
                        part_header = file_header
                    file_md_parts.append(
                        f"{part_header}"
                        f"```{language}\n"
                        f"{part_content.rstrip()}\n"
                        f"```\n\n---\n"
                    )
            except Exception as e:
                print(f"⚠️ Skipping file {file_path}: {e}", file=sys.stderr)

        # Now chunk the output
        current_chunk = []
        current_size = 0
        chunk_num = 1

        # Add header to first chunk
        header_len = len(header)
        if header_len > MAX_CHUNK_CHARS:
            raise ValueError("Header and tree too large for a single chunk.")
        current_chunk.append(header)
        current_size += header_len

        # Add file md parts
        for part in file_md_parts:
            part_len = len(part)
            if current_size + part_len > MAX_CHUNK_CHARS:
                # Write current chunk
                output_path = Path(f"{output_base}_{chunk_num}.md")
                output_path.write_text("".join(current_chunk), encoding="utf-8")
                print(f"✅ Chunk {chunk_num} exported to {output_path.resolve()}")
                chunk_num += 1
                current_chunk = []
                current_size = 0

            current_chunk.append(part)
            current_size += part_len

        # Add summary to last chunk
        summary_len = len(summary)
        if current_size + summary_len > MAX_CHUNK_CHARS:
            # Write current, start new for summary
            output_path = Path(f"{output_base}_{chunk_num}.md")
            output_path.write_text("".join(current_chunk), encoding="utf-8")
            print(f"✅ Chunk {chunk_num} exported to {output_path.resolve()}")
            chunk_num += 1
            current_chunk = [summary]
        else:
            current_chunk.append(summary)

        # Write final chunk
        if current_chunk:
            output_path = Path(f"{output_base}_{chunk_num}.md")
            output_path.write_text("".join(current_chunk), encoding="utf-8")
            print(f"✅ Chunk {chunk_num} exported to {output_path.resolve()}")

    except Exception as e:
        print(f"❌ Error: {e}", file=sys.stderr)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Export codebase to chunked Markdown files for processing"
    )
    parser.add_argument("folder", help="Path to the codebase folder")
    parser.add_argument(
        "-o", "--output", default=OUTPUT_BASE, help="Base name for output markdown chunks (e.g., code_chunk)"
    )
    args = parser.parse_args()
    parse_codebase(args.folder, args.output)