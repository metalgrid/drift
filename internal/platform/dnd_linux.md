# Linux Drag-and-Drop Notes

This documents the Dolphin drag-and-drop compatibility fix for the Linux UI drop zone.

## What changed

- The drop target now accepts multiple GTK payload types instead of only `glib.TypeString`.
- `GdkFileList` is accepted as the primary payload type for file-manager drags.
- URI-string parsing remains as a fallback for `text/uri-list` style payloads.
- Drop actions now allow `Copy|Move` negotiation to avoid denied cursor when the source prefers move.

## Implementation

- File: `internal/platform/ui_linux.go`
- Drop target setup uses:
  - `gtk.NewDropTarget(glib.TypeInvalid, gdk.ActionCopy|gdk.ActionMove)`
  - `dropTarget.SetGTypes([]glib.Type{gdk.GTypeFileList, glib.TypeString})`
- Dropped values are normalized by `droppedPaths(...)`.

## Regression coverage

- File: `internal/platform/parse_uris_linux_test.go`
- Test cases cover single/multiple URIs, CRLF input, malformed lines, non-file schemes, encoded paths, and empty input.
