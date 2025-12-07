# Concepts

This application consists of three views:

- **Main View** :: Search and preview options
- **Help View** :: Display this help page
- **Value View** :: Show the current value of an option
- **Scope Select View** :: Select scope to use

A **purple border** indicates the active (focused) view. Keybinds will only work
in the context of the currently active view.

To quit this application, press `Ctrl+C` or `Esc` from the main view.

---

## Main View

The main view appears when the application starts. It contains two _windows_:

- **Search Window** (left) :: User types to see available options
- **Preview Window** (right) :: Displays info about the selected option

Press `<Tab>` to switch focus between the two windows.

Press `<Enter>` to view the current value of a selected option, if available;
this will open the **value view**.

Press `Ctrl+O` to open the scope select view.

### Search Window

There are two modes of search: **fuzzy search** and **regex search**. Fuzzy
search is the default mode, and uses ranked approximate string matching. Regex
mode allows using RE2-style regular expressions for more exact matching.

Switch between these modes using `Ctrl+F`. Fuzzy mode is indicated by a `> `
prompt, while regex mode is indicated by a `(^$) ` prompt in the search bar.

Use the `Up` + `Down` arrows to navigate the results. As you move through the
list, the **Preview Window** updates automatically.

### Preview Window

Shows detailed information about the selected option.

Use the arrow keys or `h`, `j`, `k`, and `l` to scroll around.

## Value View

Displays the current value of the selected option (if it can be evaluated).

Use the arrow keys or `h`, `j`, `k`, and `l` to scroll around.

Press `<Esc>` or `q` to close this window.

## Scope Select View

Shows all available scopes defined in the configuration, if there is more than
one. May not always be applicable.

Use the arrow keys or `j` and `k` to scroll the list.

Use `/` to search the list of scopes.

To switch to the selected scope, press `Enter`; if successful, this redirects
back to the main view automatically.

Press `<Esc>` or `q` to close this window.

## Help View

Use the arrow keys or `h`, `j`, `k`, and `l` to scroll around.

Press `<Esc>` or `q` to close this window.
