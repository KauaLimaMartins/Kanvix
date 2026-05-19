# Kanvix — Major Update Plan

## 1. Rebrand to "Kanvix"
- Update app title, metadata in `__root.tsx`
- Update sidebar branding in `Sidebar.tsx`
- Update any "Lovable App" references

## 2. WYSIWYG Editor Fix
- Replace current Tiptap setup with a solid, production-ready config
- Add extensions: StarterKit (bold/italic/lists/code/headings) + Underline
- Install: `@tiptap/extension-underline`
- Add toolbar with buttons for: bold, italic, underline, H1/H2, bullet list, ordered list, code block
- Fix content persistence: use `onUpdate` with debounce, sync via `editor.getHTML()` to store
- Use `useEditor` with proper `content` initialization and `editor.commands.setContent` when taskId changes

## 3. Auth + Workspace Selection Flow
- New routes:
  - `/login` — mocked login screen, "Kanvix" branding, sets auth in store
  - `/workspaces` — workspace picker grid with icons/colors
- Update `src/routes/index.tsx` to redirect based on auth state
- Add `auth` slice to zustand store: `{ isAuthenticated, login(), logout() }`
- Add icon + color to Workspace type; update seed data (Company, Freelances, Personal Projects)
- Animated cards with hover effects (framer-motion)

## 4. Labels Management
- Add `Label` type: `{ id, workspaceId, name, color }`
- Add labels CRUD to zustand store
- New route: `/w/$workspaceId/labels` — list/create/edit/delete labels
- Add "Labels" item to sidebar
- Update Task type: `labels: string[]` are label IDs (already are strings; clarify semantics)
- Remove mock hardcoded labels from `mockData.ts` — seed a few real labels per workspace
- TaskDrawer: replace any free-form label input with a multiselect of workspace labels

## 5. Animations (framer-motion)
- Install `framer-motion`
- Animate: sidebar items hover, project cards, workspace cards, kanban card drag (already via dnd-kit), drawer slide
- Page transitions via `AnimatePresence` in `__root.tsx` Outlet wrapper
- Subtle: opacity + 4px Y translate, 200ms

## 6. UI polish
- Refine spacing, typography in sidebar, topbar, kanban
- Better empty states

## Files to create
- `src/routes/login.tsx`
- `src/routes/workspaces.tsx`
- `src/routes/w.$workspaceId.labels.tsx`
- `src/components/task/RichEditor.tsx` (rewrite of TaskEditor)
- `src/components/labels/LabelPicker.tsx`

## Files to edit
- `src/types.ts` — Label, Workspace icon/color
- `src/store/useAppStore.ts` — auth slice, labels CRUD
- `src/lib/mockData.ts` — workspaces with icons/colors, seed labels, remove hardcoded task labels
- `src/routes/__root.tsx` — title, AnimatePresence
- `src/routes/index.tsx` — redirect logic
- `src/components/layout/Sidebar.tsx` — branding, Labels link, animations
- `src/components/layout/Topbar.tsx` — branding
- `src/components/task/TaskDrawer.tsx` — use new editor + label picker
- `src/components/kanban/Card.tsx` — hover anim, label chips from store
- `src/components/kanban/Column.tsx`, `Board.tsx` — minor anim polish
- `src/routes/w.$workspaceId.index.tsx` — animated project cards

## Out of scope
- Real auth, real DB persistence (still localStorage — "persist in database" interpreted as persisted store)
